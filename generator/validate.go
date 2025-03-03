package generator

import (
	"log"
	"slices"
	"strconv"

	v3 "github.com/TeCHiScy/protoc-gen-openapi/openapiv3"
	"github.com/envoyproxy/protoc-gen-validate/validate"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (g *OpenAPIv3Generator) addValidationRules(fieldSchema *v3.SchemaOrReference, field protoreflect.FieldDescriptor) {
	validationRules := proto.GetExtension(field.Options(), validate.E_Rules)
	if validationRules == nil {
		return
	}
	fieldRules, ok := validationRules.(*validate.FieldRules)
	if !ok {
		return
	}
	schema, ok := fieldSchema.Oneof.(*v3.SchemaOrReference_Schema)
	if !ok {
		return
	}

	// TODO: implement map validation
	if field.IsMap() {
		return
	}

	if field.IsList() {
		repeatedRules := fieldRules.GetRepeated()
		if repeatedRules == nil {
			// no rules
			return
		}
		// MinItems specifies that this field must have the specified number of
		// items at a minimum
		// MaxItems specifies that this field must have the specified number of
		// items at a maximum
		// Unique specifies that all elements in this field must be unique. This
		// contraint is only applicable to scalar and enum types (messages are not
		// supported).
		// Items specifies the contraints to be applied to each item in the field.
		// Repeated message fields will still execute validation against each item
		// unless skip is specified here.
		// IgnoreEmpty specifies that the validation rules of this field should be
		// evaluated only if the field is not empty
		if repeatedRules.MinItems != nil {
			schema.Schema.MinItems = int64(*repeatedRules.MinItems)
		}
		if repeatedRules.MaxItems != nil {
			schema.Schema.MaxItems = int64(*repeatedRules.MaxItems)
		}

		// pull out the array items field rules
		fieldRules := repeatedRules.Items
		if fieldRules == nil {
			// no item specific rules
			return
		}
		schema := schema.Schema.Items.SchemaOrReference[0]
		fieldRule(fieldRules, field, schema.Oneof.(*v3.SchemaOrReference_Schema))
		return
	}

	fieldRule(fieldRules, field, schema)
}

func fieldRule(r *validate.FieldRules, field protoreflect.FieldDescriptor, schema *v3.SchemaOrReference_Schema) {
	kind := field.Kind()
	switch kind {
	case protoreflect.MessageKind:
		return // TODO: Implement message validators from protoc-gen-validate
	case protoreflect.StringKind:
		setStringRules(r.GetString_(), schema)
	case protoreflect.Int32Kind:
		setInt32Rules(r.GetInt32(), schema)
	case protoreflect.Int64Kind:
		setInt64Rules(r.GetInt64(), schema)
	case protoreflect.EnumKind:
		setEnumRules(r.GetEnum(), field, schema)
	// TODO: implement protoc-gen-validate rules for the following types
	case protoreflect.Sint32Kind, protoreflect.Uint32Kind,
		protoreflect.Sint64Kind, protoreflect.Uint64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind, protoreflect.Sfixed64Kind,
		protoreflect.Fixed64Kind:

	case protoreflect.BoolKind:

	case protoreflect.FloatKind, protoreflect.DoubleKind:

	case protoreflect.BytesKind:

	default:
		log.Printf("(TODO) Unsupported field type: %+v", fullMessageTypeName(field.Message()))
	}
}

func setStringRules(r *validate.StringRules, schema *v3.SchemaOrReference_Schema) {
	if r == nil {
		return
	}

	// set format
	// format is an open value, so you can use any formats, even not those defined by the OpenAPI Specification
	if r.GetEmail() {
		schema.Schema.Format = "email"
	} else if r.GetHostname() {
		schema.Schema.Format = "hostname"
	} else if r.GetIp() {
		schema.Schema.Format = "ip"
	} else if r.GetIpv4() {
		schema.Schema.Format = "ipv4"
	} else if r.GetIpv6() {
		schema.Schema.Format = "ipv6"
	} else if r.GetUri() {
		schema.Schema.Format = "uri"
	} else if r.GetUriRef() {
		schema.Schema.Format = "uri_ref"
	} else if r.GetUuid() {
		schema.Schema.Format = "uuid"
	}

	// set min/max
	if r.GetMinLen() > 0 {
		schema.Schema.MinLength = int64(r.GetMinLen())
	}
	if r.GetMaxLen() > 0 {
		schema.Schema.MaxLength = int64(r.GetMaxLen())
	}

	// set Pattern
	if r.GetPattern() != "" {
		schema.Schema.Pattern = r.GetPattern()
	}
}

func setInt32Rules(r *validate.Int32Rules, schema *v3.SchemaOrReference_Schema) {
	if r == nil {
		return
	}

	if r.Gte != nil {
		schema.Schema.Minimum = float64(r.GetGte())
	}
	if r.Lte != nil {
		schema.Schema.Maximum = float64(r.GetLte())
	}
	if r.Gt != nil {
		schema.Schema.Minimum = float64(r.GetGt())
		schema.Schema.ExclusiveMinimum = true
	}
	if r.Lt != nil {
		schema.Schema.Maximum = float64(r.GetLt())
		schema.Schema.ExclusiveMaximum = true
	}
	if r.Const != nil {
		schema.Schema.Enum = []*v3.Any{{Yaml: strconv.FormatInt(int64(r.GetConst()), 10)}}
	}
	if r.In != nil {
		for _, v := range r.GetIn() {
			schema.Schema.Enum = append(schema.Schema.Enum, &v3.Any{Yaml: strconv.FormatInt(int64(v), 10)})
		}
	}
	if r.NotIn != nil {
		if schema.Schema.Not == nil {
			schema.Schema.Not = &v3.Schema{}
		}
		for _, v := range r.GetNotIn() {
			schema.Schema.Not.Enum = append(schema.Schema.Not.Enum, &v3.Any{Yaml: strconv.FormatInt(int64(v), 10)})
		}
	}
}

func setInt64Rules(r *validate.Int64Rules, schema *v3.SchemaOrReference_Schema) {
	if r == nil {
		return
	}

	if r.Gte != nil {
		schema.Schema.Minimum = float64(r.GetGte())
	}
	if r.Lte != nil {
		schema.Schema.Maximum = float64(r.GetLte())
	}
	if r.Gt != nil {
		schema.Schema.Minimum = float64(r.GetGt())
		schema.Schema.ExclusiveMinimum = true
	}
	if r.Lt != nil {
		schema.Schema.Maximum = float64(r.GetLt())
		schema.Schema.ExclusiveMaximum = true
	}
	if r.Const != nil {
		schema.Schema.Enum = []*v3.Any{{Yaml: strconv.FormatInt(int64(r.GetConst()), 10)}}
	}
	if r.In != nil {
		for _, v := range r.GetIn() {
			schema.Schema.Enum = append(schema.Schema.Enum, &v3.Any{Yaml: strconv.FormatInt(int64(v), 10)})
		}
	}
	if r.NotIn != nil {
		if schema.Schema.Not == nil {
			schema.Schema.Not = &v3.Schema{}
		}
		for _, v := range r.GetNotIn() {
			schema.Schema.Not.Enum = append(schema.Schema.Not.Enum, &v3.Any{Yaml: strconv.FormatInt(int64(v), 10)})
		}
	}
}

func setEnumRules(r *validate.EnumRules, field protoreflect.FieldDescriptor, schema *v3.SchemaOrReference_Schema) {
	if r == nil {
		return
	}

	// TODO 规则对于 comments 的处理
	// we don't check enumRules.DefinedOnly because we already list the set of valid enums
	if r.Const != nil {
		setEnumInRules(field, schema, []int32{r.GetConst()}, false)
	}
	if r.In != nil {
		setEnumInRules(field, schema, r.GetIn(), false)
	}
	if r.NotIn != nil {
		setEnumInRules(field, schema, r.GetNotIn(), true)
	}
}

func setEnumInRules(field protoreflect.FieldDescriptor, schema *v3.SchemaOrReference_Schema, vals []int32, reversed bool) {
	valids := make([]string, 0, len(vals))
	if schema.Schema.Type == "string" {
		m := enumValues(field)
		for _, v := range vals {
			if name, ok := m[protoreflect.EnumNumber(v)]; ok {
				valids = append(valids, string(name))
			}
		}
	} else if schema.Schema.Type == "integer" {
		for _, v := range vals {
			valids = append(valids, strconv.FormatInt(int64(v), 10))
		}
	}

	filtered := make([]*v3.Any, 0, len(schema.Schema.Enum))
	for _, e := range schema.Schema.Enum {
		if !reversed && slices.Contains(valids, e.GetYaml()) {
			filtered = append(filtered, e)
		} else if reversed && !slices.Contains(valids, e.GetYaml()) {
			filtered = append(filtered, e)
		}
	}
	schema.Schema.Enum = filtered
}

// enumValues returns of list of enum ids for a given field
func enumValues(field protoreflect.FieldDescriptor) map[protoreflect.EnumNumber]protoreflect.Name {
	values := field.Enum().Values()
	m := map[protoreflect.EnumNumber]protoreflect.Name{}
	for i := range values.Len() {
		v := values.Get(i)
		m[v.Number()] = v.Name()
	}
	return m
}
