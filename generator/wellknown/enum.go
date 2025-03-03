// Copyright 2020 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, softwis
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package wellknown

import (
	"strconv"
	"strings"

	v3 "github.com/TeCHiScy/protoc-gen-openapi/openapiv3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func NewEnumSchema(enumType *string, field protoreflect.FieldDescriptor) *v3.SchemaOrReference {
	schema := &v3.Schema{Format: "enum"}
	var useString bool
	schema.Type = "integer"
	if enumType != nil && *enumType == "string" {
		useString = true
		schema.Type = "string"
	}

	schema.Enum = make([]*v3.Any, 0, field.Enum().Values().Len())
	vals := field.Enum().Values()
	for i := range vals.Len() {
		v := vals.Get(i)
		// skip default unspecified values
		// https://protobuf.dev/programming-guides/proto3/#enum-default
		if strings.HasSuffix(string(v.Name()), "_UNSPECIFIED") ||
			strings.HasSuffix(string(v.Name()), "_UNKNOWN") {
			continue
		}
		var item string
		if useString {
			item = string(v.Name())
		} else {
			item = strconv.FormatInt(int64(v.Number()), 10)
		}
		schema.Enum = append(schema.Enum, &v3.Any{Yaml: item})
	}
	return &v3.SchemaOrReference{
		Oneof: &v3.SchemaOrReference_Schema{
			Schema: schema}}
}

func AddEnumComments(fieldSchema *v3.SchemaOrReference, field protoreflect.FieldDescriptor) {
	schema := fieldSchema.GetSchema()
	if schema == nil {
		return
	}

	if field.IsMap() {
		if mv := field.MapValue(); mv.Kind() == protoreflect.EnumKind {
			addEnumComment(mv, schema, schema.GetAdditionalProperties().GetSchemaOrReference().GetSchema())
		}
		return
	}
	if field.IsList() {
		if field.Kind() == protoreflect.EnumKind {
			addEnumComment(field, schema, schema.GetItems().SchemaOrReference[0].GetSchema())
		}
		return
	}
	if field.Kind() == protoreflect.EnumKind {
		addEnumComment(field, schema, schema)
	}
}

func addEnumComment(field protoreflect.FieldDescriptor, schema *v3.Schema, enumSchema *v3.Schema) {
	vals := field.Enum().Values()
	m := map[string]string{}
	descs := map[string]string{}
	for i := range vals.Len() {
		v := vals.Get(i)
		number := strconv.FormatInt(int64(v.Number()), 10)
		m[string(v.Name())] = number
		if opt, ok := proto.GetExtension(v.Options(), v3.E_Desc).(string); ok {
			descs[number] = opt
		} else {
			descs[number] = string(v.Name())
		}
	}

	var comments []string
	for _, e := range enumSchema.Enum {
		if enumSchema.Type == "string" {
			id := m[e.GetYaml()]
			comments = append(comments, id+": "+descs[id])
		} else if enumSchema.Type == "integer" {
			comments = append(comments, e.GetYaml()+": "+descs[e.GetYaml()])
		}
	}
	if schema.Description != "" {
		schema.Description += " \\\n"
	}
	schema.Description += strings.Join(comments, " \\\n")
}
