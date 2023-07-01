package gonuts

import "reflect"

/* example usage for GetStructFieldNamesByTagsValues:
package main

import (
	"fmt"
	"reflect"
)

type SourceStruct struct {
	Field1 string `filter:"public"`
	Field2 int    `filter:"private"`
	Field3 bool   `filter:"public"`
}

func main() {
	tagsValues := map[string]any{
		"filter": "public",
	}
	fieldNames := GetStructFieldNamesByTagsValues(SourceStruct{}, tagsValues)
	fmt.Printf("Field names: %+v\n", fieldNames)
}
*/
// this function takes a struct and returns a slice of strings containing the names of the fields that have the given combination of a map[string]any with tags and values
func GetStructFieldNamesByTagsValues(source any, tagsValues map[string]any) []string {
	sourceType := reflect.TypeOf(source)
	var filteredFields []string
	for i := 0; i < sourceType.NumField(); i++ {
		field := sourceType.Field(i)
		if fieldHasTagsValues(field, tagsValues) {
			filteredFields = append(filteredFields, field.Name)
		}
	}
	return filteredFields
}

func fieldHasTagsValues(field reflect.StructField, tagsValues map[string]any) bool {
	for tag := range tagsValues {
		val, ok := field.Tag.Lookup(tag)
		if !ok {
			return false
		}
		if field.Tag.Get(tag) != val {
			return false
		}
	}
	return true
}
