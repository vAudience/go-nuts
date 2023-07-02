package gonuts

import (
	"reflect"
	"strings"
)

/* example usages:
package main

import (
	"fmt"
)

type SourceStruct struct {
	Field1 string `filter:"public"`
	Field2 int    `filter:"private,admin"`
	Field3 bool   `filter:"public,user"`
	Field4 bool
}

func main() {
	// filter to keep all public fields
	source := SourceStruct{
		Field1: "public",
		Field2: 2,
		Field3: true,
		Field4: false,
	}
	filtered := CreateFilteredStruct(source, []string{"public"}, nil)
	fmt.Println("(1) filter to keep all public fields:", filtered)
	// filter to remove all admin fields
	source = SourceStruct{
		Field1: "public",
		Field2: 2,
		Field3: true,
		Field4: false,
	}
	filtered = CreateFilteredStruct(source, []string{""}, []string{"admin"})
	fmt.Println("(2) filter to remove all admin fields and keep any others:", filtered)
	// filter to keep all admin fields
	source = SourceStruct{
		Field1: "public",
		Field2: 2,
		Field3: true,
		Field4: false,
	}
	filtered = CreateFilteredStruct(source, []string{"admin"}, nil)
	fmt.Println("(3) filter to keep all admin fields:", filtered)
	// filter to keep all fields with any filter value
	source = SourceStruct{
		Field1: "public",
		Field2: 2,
		Field3: true,
		Field4: false,
	}
	filtered = CreateFilteredStruct(source, []string{""}, nil)
	fmt.Println("(4) filter to keep all fields with any filter value:", filtered)
}

*/

// @@Summary CreateFilteredStruct creates a new struct with only the fields that have any of the given filterValuesToKeep AND do not have any of the filterValuesToRemove.
func CreateFilteredStruct(source any, filterValuesToKeep []string, filterValuesToRemove []string) any {
	sourceType := reflect.TypeOf(source)
	sourceValue := reflect.ValueOf(source)
	destinationType := reflect.StructOf(createFilteredStructFields(sourceType, filterValuesToKeep, filterValuesToRemove))
	destinationValue := reflect.New(destinationType).Elem()
	for i := 0; i < destinationType.NumField(); i++ {
		fieldName := destinationType.Field(i).Name
		destinationValue.FieldByName(fieldName).Set(sourceValue.FieldByName(fieldName))
	}
	return destinationValue.Interface()
}

// @@Summary CreateFilteredStructFields creates a new struct with only the fields that have any of the given filterValuesToKeep AND do not have any of the filterValuesToRemove.
func createFilteredStructFields(sourceType reflect.Type, filterValuesToKeep []string, filterValuesToRemove []string) []reflect.StructField {
	var filteredFields []reflect.StructField
	for i := 0; i < sourceType.NumField(); i++ {
		field := sourceType.Field(i)
		tagMapToKeep := map[string][]string{
			"filter": filterValuesToKeep,
		}
		tagMapToRemove := map[string][]string{
			"filter": filterValuesToRemove,
		}
		if fieldHasTagsValues(field, tagMapToKeep, tagMapToRemove) {
			filteredFields = append(filteredFields, field)
		}
	}
	return filteredFields
}

// this function takes a struct and returns a slice of strings containing the names of the fields that have the given combination of a map[string]any with tags and values
func GetStructFieldNamesByTagsValues(source any, tagsValues map[string][]string) []string {
	sourceType := reflect.TypeOf(source)
	var filteredFields []string
	for i := 0; i < sourceType.NumField(); i++ {
		field := sourceType.Field(i)
		if fieldHasTagsValues(field, tagsValues, nil) {
			filteredFields = append(filteredFields, field.Name)
		}
	}
	return filteredFields
}

// @@Summary fieldHasTagsValues returns true if the field has all the given tags and values, and none of the given tags and values.
func fieldHasTagsValues(field reflect.StructField, tagsValuesToKeep map[string][]string, tagsValuesToRemove map[string][]string) bool {
	for tag, value := range tagsValuesToKeep {
		found := false
		for _, v := range value {
			if fieldHasTagValue(field, tag, v) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	for tag, value := range tagsValuesToRemove {
		found := false
		for _, v := range value {
			if fieldHasTagValue(field, tag, v) {
				found = true
				break
			}
		}
		if found {
			return false
		}
	}
	return true
}

// @@Summary fieldHasTagValue returns true if the field has the given tag and value.
func fieldHasTagValue(field reflect.StructField, tag string, value string) bool {
	tagValue := field.Tag.Get(tag)
	if tagValue == "" {
		if value == "" {
			return true
		} else {
			return false
		}
	}
	// split the tagValue by comma and trim the spaces, then strings.Equalfold-compare each slice-element to the given value. if any match, return true.
	for _, v := range strings.Split(tagValue, ",") {
		if strings.EqualFold(strings.TrimSpace(v), value) {
			return true
		}
	}
	return false
}
