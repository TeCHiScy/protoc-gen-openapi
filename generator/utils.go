// Copyright 2020 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package generator

import (
	"slices"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// appendUnique appends a string, to a string slice, if the string is not already in the slice
func appendUnique(s []string, e string) []string {
	if !slices.Contains(s, e) {
		return append(s, e)
	}
	return s
}

// singular produces the singular form of a collection name.
func singular(plural string) string {
	if strings.HasSuffix(plural, "ves") {
		return strings.TrimSuffix(plural, "ves") + "f"
	}
	if strings.HasSuffix(plural, "ies") {
		if plural == "series" || plural == "species" {
			return plural
		} else {
			return strings.TrimSuffix(plural, "ies") + "y"
		}
	}
	if strings.HasSuffix(plural, "s") {
		return strings.TrimSuffix(plural, "s")
	}
	return plural
}

func getValueKind(message protoreflect.MessageDescriptor) string {
	valueField := getValueField(message)
	return valueField.Kind().String()
}

func getValueField(message protoreflect.MessageDescriptor) protoreflect.FieldDescriptor {
	fields := message.Fields()
	return fields.ByName("value")
}
