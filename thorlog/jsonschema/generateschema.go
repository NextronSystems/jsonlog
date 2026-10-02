package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/NextronSystems/jsonlog"
	"github.com/NextronSystems/jsonlog/thorlog/v3"
	_ "github.com/NextronSystems/jsonlog/thorlog/v3/audittrail"
	"github.com/invopop/jsonschema"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

var objectType = reflect.TypeOf((*jsonlog.Object)(nil)).Elem()

func makeObjectSchema() (mainEntry string, defs map[string]*jsonschema.Schema) {
	var allLogObjects []*jsonschema.Schema
	var logObjectTypes []any
	var reflector jsonschema.Reflector
	reflector.AllowAdditionalProperties = true
	// Walks subdirectories too, so this also covers audittrail.
	err := reflector.AddGoComments("github.com/NextronSystems/jsonlog/thorlog/v3", "../v3")
	if err != nil {
		panic(err)
	}
	defs = map[string]*jsonschema.Schema{}

	// Sort the object type names to have a stable output.
	var objectTypeNames = slices.Collect(maps.Keys(thorlog.LogObjectTypes))
	slices.Sort(objectTypeNames)

	reflector.Mapper = func(r reflect.Type) *jsonschema.Schema {
		if r.Kind() == reflect.Interface {
			if r.Implements(objectType) {
				// r is an interface that implements jsonlog.Object.
				// Since we know all types that implement jsonlog.Object,
				// we can filter for the types that implement the interface,
				// and generate a oneOf schema for them.
				var implementations = &jsonschema.Schema{}
				for _, typename := range objectTypeNames {
					t := thorlog.LogObjectTypes[typename]
					if reflect.TypeOf(t).Implements(r) {
						structName := reflect.TypeOf(t).Elem().Name()
						implementations.OneOf = append(implementations.OneOf, &jsonschema.Schema{
							Ref: "#/$defs/" + structName,
						})
					}
				}
				if _, ok := defs[r.Name()]; !ok {
					defs[r.Name()] = implementations
				}
				return &jsonschema.Schema{
					Ref: "#/$defs/" + r.Name(),
				}
			} else {
				panic(fmt.Sprintf("Use of unknown interface %s", r.Name()))
			}
		}
		return nil
	}
	for _, typename := range objectTypeNames {
		schema := reflector.Reflect(thorlog.LogObjectTypes[typename])
		refName := strings.TrimPrefix(schema.Ref, "#/$defs/")
		typeSchema := schema.Definitions[refName]
		typenameDef, ok := typeSchema.Properties.Get("type")
		if !ok {
			panic("type property not found in " + refName)
		}
		typenameDef.Const = typename
		allLogObjects = append(allLogObjects, schema)
		logObjectTypes = append(logObjectTypes, typename)
		defs[refName] = typeSchema
	}
	var logObjectSchema = &jsonschema.Schema{
		Properties: orderedmap.New[string, *jsonschema.Schema](),
	}
	logObjectSchema.Properties.Set("type", &jsonschema.Schema{
		Type: "string",
		Enum: logObjectTypes,
	})
	logObjectSchema.OneOf = allLogObjects

	const logObjectSchemaName = "object"
	defs[logObjectSchemaName] = logObjectSchema

	return logObjectSchemaName, defs
}

func main() {
	var logEventSchema jsonschema.Schema
	if len(os.Args) == 2 && os.Args[1] == "log" {
		logEventSchema = jsonschema.Schema{
			Version:     jsonschema.Version,
			ID:          "https://www.nextron-systems.com/schemas/thorlog/v3/thor-event.json",
			Definitions: map[string]*jsonschema.Schema{},
			Title:       "ThorEvent",
			OneOf: []*jsonschema.Schema{
				{
					Ref: "#/$defs/Assessment",
				},
				{
					Ref: "#/$defs/Message",
				},
			},
		}
	} else if len(os.Args) == 2 && os.Args[1] == "audittrail" {
		logEventSchema = jsonschema.Schema{
			Version:     jsonschema.Version,
			ID:          "https://www.nextron-systems.com/schemas/thorlog/v3/thor-audit-entry.json",
			Definitions: map[string]*jsonschema.Schema{},
			Title:       "ThorAuditEntry",
			OneOf: []*jsonschema.Schema{
				{
					Ref: "#/$defs/AuditMessage",
				},
				{
					Ref: "#/$defs/AuditRecord",
				},
			},
		}
	} else {
		fmt.Fprintf(os.Stderr, "Usage: %s (log|audittrail)\n", os.Args[0])
		os.Exit(2)
	}

	entry, defs := makeObjectSchema()
	for key, value := range defs {
		logEventSchema.Definitions[key] = value
	}

	flatten(logEventSchema.Definitions[entry], logEventSchema.Definitions)
	prune(&logEventSchema)

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(logEventSchema)
}

func flatten(schema *jsonschema.Schema, definitions jsonschema.Definitions) {
	if schema == nil {
		return
	}
	schema.Version = ""
	schema.ID = ""
	for key, subschema := range schema.Definitions {
		if _, ok := definitions[key]; ok {
			continue
		}
		definitions[key] = subschema
		flatten(subschema, definitions)
	}
	schema.Definitions = nil

	flatten(schema.If, definitions)
	flatten(schema.Then, definitions)
	flatten(schema.Else, definitions)
	for _, subschema := range schema.AnyOf {
		flatten(subschema, definitions)
	}
	for _, subschema := range schema.AllOf {
		flatten(subschema, definitions)
	}
	for _, subschema := range schema.OneOf {
		flatten(subschema, definitions)
	}
}

// prune removes all definitions that are not reachable from the root schema via $ref.
func prune(root *jsonschema.Schema) {
	reachable := map[string]bool{}
	var visit func(schema *jsonschema.Schema)
	visit = func(schema *jsonschema.Schema) {
		// Most subschema fields are nil.
		if schema == nil {
			return
		}
		// Check !reachable[name] so we won't run into an endless loop
		// and we only visit each definition once.
		if name, ok := strings.CutPrefix(schema.Ref, "#/$defs/"); ok && !reachable[name] {
			def, ok := root.Definitions[name]
			if !ok {
				panic("dangling reference " + schema.Ref)
			}
			// Mark the definition before the recursive call.
			reachable[name] = true
			visit(def)
		}
		// Definitions are not visited, since they are only reachable via $ref.
		// Visit all applicators https://www.learnjsonschema.com/2020-12/applicator/
		// and contentSchema, which could also contain references.
		children := slices.Concat(
			[]*jsonschema.Schema{
				schema.Not, schema.If, schema.Then, schema.Else, schema.Items, schema.Contains,
				schema.AdditionalProperties, schema.PropertyNames, schema.ContentSchema,
			},
			schema.AllOf, schema.AnyOf, schema.OneOf, schema.PrefixItems,
			slices.Collect(maps.Values(schema.PatternProperties)),
			slices.Collect(maps.Values(schema.DependentSchemas)),
		)
		for pair := schema.Properties.Oldest(); pair != nil; pair = pair.Next() {
			children = append(children, pair.Value)
		}
		for _, child := range children {
			visit(child)
		}
	}
	visit(root)
	maps.DeleteFunc(root.Definitions, func(name string, _ *jsonschema.Schema) bool { return !reachable[name] })
}
