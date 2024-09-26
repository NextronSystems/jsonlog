package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/NextronSystems/jsonlog"
	"github.com/NextronSystems/jsonlog/thorlog/v3"
	"github.com/invopop/jsonschema"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func makeObjectSchema() *jsonschema.Schema {
	var allLogObjects []*jsonschema.Schema
	var logObjectTypes []any
	var reflector jsonschema.Reflector
	reflector.Mapper = func(r reflect.Type) *jsonschema.Schema {
		if r.Kind() == reflect.Interface {
			switch r {
			case reflect.TypeOf((*jsonlog.Object)(nil)).Elem():
				return &jsonschema.Schema{
					Ref: "#/$defs/object",
				}
			case reflect.TypeOf((*thorlog.Permissions)(nil)).Elem():
				return &jsonschema.Schema{
					OneOf: []*jsonschema.Schema{
						{
							Ref: "#/$defs/WindowsPermissions",
						},
						{
							Ref: "#/$defs/UnixPermissions",
						},
					},
				}
				return nil
			default:
				panic(fmt.Sprintf("Use of unknown interface %s", r.Name()))
			}
		}
		return nil
	}
	for typename, object := range thorlog.LogObjectTypes {
		schema := reflector.Reflect(object)
		refName := strings.TrimPrefix(schema.Ref, "#/$defs/")
		typeSchema := schema.Definitions[refName]
		typenameDef, ok := typeSchema.Properties.Get("type")
		if !ok {
			panic("type property not found in " + refName)
		}
		typenameDef.Const = typename
		allLogObjects = append(allLogObjects, schema)
		logObjectTypes = append(logObjectTypes, typename)
	}
	var logObjectSchema = &jsonschema.Schema{
		Properties: orderedmap.New[string, *jsonschema.Schema](),
	}
	logObjectSchema.Properties.Set("summary", &jsonschema.Schema{Type: "string"})
	logObjectSchema.Properties.Set("type", &jsonschema.Schema{
		Type: "string",
		Enum: logObjectTypes,
	})
	logObjectSchema.AnyOf = allLogObjects

	return logObjectSchema
}

func main() {
	logEventSchema := jsonschema.Schema{
		Version: jsonschema.Version,
		ID:      "https://www.nextron-systems.com/schemas/thorlog/v3/log-event.json",
		Ref:     "#/$defs/event",
		Definitions: map[string]*jsonschema.Schema{
			"event": {
				Type: "object",
				AnyOf: []*jsonschema.Schema{
					{
						Ref: "#/$defs/Finding",
					},
					{
						Ref: "#/$defs/Message",
					},
				},
			},
		},
	}
	logEventSchema.Definitions["object"] = makeObjectSchema()

	flatten(logEventSchema.Definitions["object"], logEventSchema.Definitions)

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
}
