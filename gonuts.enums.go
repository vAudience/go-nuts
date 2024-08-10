package gonuts

import (
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"
)

// EnumDefinition represents the structure of an enum
type EnumDefinition struct {
	Name   string   // The name of the enum type
	Values []string // The values of the enum
}

// GenerateEnum generates the enum code based on the provided definition
//
// This function takes an EnumDefinition and returns a string containing
// the generated Go code for a type-safe enum. The generated code includes:
// - A new type based on int
// - Constant values for each enum value
// - A String() method for string representation
// - An IsValid() method to check if a value is valid
// - A Parse[EnumName]() function to convert strings to enum values
// - MarshalJSON() and UnmarshalJSON() methods for JSON encoding/decoding
//
// Parameters:
//   - def: An EnumDefinition struct containing the enum name and values
//
// Returns:
//   - string: The generated Go code for the enum
//   - error: An error if code generation fails
//
// Example usage:
//
//	def := EnumDefinition{
//	    Name:   "Color",
//	    Values: []string{"Red", "Green", "Blue"},
//	}
//	code, err := GenerateEnum(def)
//	if err != nil {
//	    log.Fatalf("Failed to generate enum: %v", err)
//	}
//	fmt.Println(code)
//
// The generated code can be used as follows:
//
//	var c Color = ColorRed
//	fmt.Println(c)                 // Output: Red
//	fmt.Println(c.IsValid())       // Output: true
//	c2, _ := ParseColor("Green")
//	fmt.Println(c2)                // Output: Green
//	jsonData, _ := json.Marshal(c)
//	fmt.Println(string(jsonData))  // Output: "Red"
func GenerateEnum(def EnumDefinition) (string, error) {
	const enumTemplate = `
// Code generated by gonuts. DO NOT EDIT.

package {{.PackageName}}

import (
	"fmt"
	"encoding/json"
)

// {{.Name}} represents an enumeration of {{.Name}} values
type {{.Name}} int

const (
	{{range $index, $value := .Values}}{{if $index}}_{{end}}{{$.Name}}{{$value}} {{$.Name}} = iota
	{{end}}
)

// {{.Name}}Values contains all valid string representations of {{.Name}}
var {{.Name}}Values = []string{
	{{range .Values}}"{{.}}",
	{{end}}
}

// String returns the string representation of the {{.Name}}
func (e {{.Name}}) String() string {
	if e < 0 || int(e) >= len({{.Name}}Values) {
		return fmt.Sprintf("Invalid{{.Name}}(%d)", int(e))
	}
	return {{.Name}}Values[e]
}

// IsValid checks if the {{.Name}} value is valid
func (e {{.Name}}) IsValid() bool {
	return e >= 0 && int(e) < len({{.Name}}Values)
}

// Parse{{.Name}} converts a string to a {{.Name}} value
func Parse{{.Name}}(s string) ({{.Name}}, error) {
	for i, v := range {{.Name}}Values {
		if v == s {
			return {{.Name}}(i), nil
		}
	}
	return {{.Name}}(-1), fmt.Errorf("invalid {{.Name}}: %s", s)
}

// MarshalJSON implements the json.Marshaler interface for {{.Name}}
func (e {{.Name}}) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for {{.Name}}
func (e *{{.Name}}) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v, err := Parse{{.Name}}(s)
	if err != nil {
		return err
	}
	*e = v
	return nil
}
`

	packageName, err := getPackageName()
	if err != nil {
		return "", fmt.Errorf("failed to get package name: %w", err)
	}

	data := struct {
		PackageName string
		EnumDefinition
	}{
		PackageName:    packageName,
		EnumDefinition: def,
	}

	var buf strings.Builder
	tmpl, err := template.New("enum").Parse(enumTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// getPackageName attempts to determine the package name of the current directory
func getPackageName() (string, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", nil, parser.PackageClauseOnly)
	if err != nil {
		return "", err
	}

	for name := range pkgs {
		return name, nil
	}

	return "", fmt.Errorf("no package found in current directory")
}

// WriteEnumToFile generates the enum code and writes it to a file
//
// This function generates the enum code based on the provided EnumDefinition
// and writes it to a file. It also formats the generated code for readability.
//
// Parameters:
//   - def: An EnumDefinition struct containing the enum name and values
//   - filename: The name of the file to write the generated code to
//
// Returns:
//   - error: An error if code generation or file writing fails
//
// Example usage:
//
//	def := EnumDefinition{
//	    Name:   "Color",
//	    Values: []string{"Red", "Green", "Blue"},
//	}
//	err := WriteEnumToFile(def, "color_enum.go")
//	if err != nil {
//	    log.Fatalf("Failed to generate enum file: %v", err)
//	}
//
// This will create a file named "color_enum.go" in the current directory
// with the generated enum code. The generated file will contain a type-safe
// Color enum with values ColorRed, ColorGreen, and ColorBlue, along with
// helper methods for string conversion, validation, and JSON marshaling/unmarshaling.
func WriteEnumToFile(def EnumDefinition, filename string) error {
	code, err := GenerateEnum(def)
	if err != nil {
		return err
	}

	// Parse the generated code
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse generated code: %w", err)
	}

	// Format the AST
	var buf strings.Builder
	err = format.Node(&buf, fset, file)
	if err != nil {
		return fmt.Errorf("failed to format code: %w", err)
	}

	// Write the formatted code to file
	err = os.WriteFile(filename, []byte(buf.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
