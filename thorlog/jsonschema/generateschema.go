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
	"github.com/invopop/jsonschema"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

var objectType = reflect.TypeOf((*jsonlog.Object)(nil)).Elem()

func makeObjectSchema() (mainEntry string, defs map[string]*jsonschema.Schema) {
	var allLogObjects []*jsonschema.Schema
	var logObjectTypes []any
	var reflector jsonschema.Reflector
	reflector.AllowAdditionalProperties = true
	err := reflector.AddGoComments("github.com/NextronSystems/jsonlog/thorlog/v3", "../v3")
	if err != nil {
		panic(err)
	}
	defs = map[string]*jsonschema.Schema{}

	// Sort the object type names to have a stable output
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
	logEventSchema := jsonschema.Schema{
		Version:     jsonschema.Version,
		ID:          "https://www.nextron-systems.com/schemas/thorlog/v3/thor-event.json",
		Definitions: map[string]*jsonschema.Schema{},
		Title:       "ThorEvent",
		OneOf: []*jsonschema.Schema{
			{
				Ref: "#/$defs/Finding",
			},
			{
				Ref: "#/$defs/Message",
			},
		},
	}
	entry, defs := makeObjectSchema()
	for key, value := range defs {
		logEventSchema.Definitions[key] = value
	}

	flatten(logEventSchema.Definitions[entry], logEventSchema.Definitions)

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
