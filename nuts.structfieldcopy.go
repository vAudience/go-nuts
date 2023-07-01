package gonuts

import "reflect"

/* example usage:
func main() {
	source := SourceStruct{
		Field1: "value1",
		Field2: 42,
		Field3: true,
	}

	destination := DestinationStruct{
		Field1: "initialValue",
		Field2: 0,
		Field3: false,
		Field4: 3.14,
	}

	copyFields(source, &destination)

	fmt.Printf("Destination after copying: %+v\n", destination)
}
*/

// copyFields copies all fields from source to destination if they have the same name and type.
// destination must be a pointer to a struct.
func CopyFields(source, destination interface{}) {
	sourceValue := reflect.ValueOf(source)
	destinationValue := reflect.ValueOf(destination).Elem()
	for i := 0; i < sourceValue.NumField(); i++ {
		sourceField := sourceValue.Field(i)
		destinationField := destinationValue.FieldByName(sourceValue.Type().Field(i).Name)
		if destinationField.IsValid() && destinationField.Type() == sourceField.Type() {
			destinationField.Set(sourceField)
		}
	}
}
