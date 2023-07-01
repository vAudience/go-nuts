package gonuts

import (
	"reflect"
)

/* example usage:
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

type DestinationStruct struct {
	Field1 string
	Field3 bool
}

func main() {
	source := SourceStruct{
		Field1: "value1",
		Field2: 42,
		Field3: true,
	}

	destination := createFilteredStruct(source, "public")

	fmt.Printf("Destination: %+v\n", destination)
}
*/

// CreateFilteredStruct creates a new struct with only the fields that have the given filter value.
// source must be a struct.
func CreateFilteredStruct(source any, filterValue string) any {
	sourceType := reflect.TypeOf(source)
	destinationType := reflect.StructOf(filterStructFields(sourceType, filterValue))
	destination := reflect.New(destinationType).Elem()
	sourceValue := reflect.ValueOf(source)
	for i := 0; i < sourceType.NumField(); i++ {
		field := sourceType.Field(i)
		value := sourceValue.Field(i)
		if field.Tag.Get("filter") == filterValue {
			destinationField := destination.FieldByName(field.Name)
			if destinationField.IsValid() && destinationField.CanSet() {
				destinationField.Set(value)
			}
		}
	}
	return destination.Addr().Interface()
}

func filterStructFields(sourceType reflect.Type, filterValue string) []reflect.StructField {
	var filteredFields []reflect.StructField
	for i := 0; i < sourceType.NumField(); i++ {
		field := sourceType.Field(i)
		if field.Tag.Get("filter") == filterValue {
			filteredFields = append(filteredFields, field)
		}
	}
	return filteredFields
}
