// Internal helper function for the nmcli package

package nmcli

import (
	"fmt"
	"reflect"
)

// Given a valid struct generates command value pairs for nmcli
func generate_commands(c ConnDetails) []string {
	output := make([]string, 0)
	// Get type
	t := reflect.TypeOf(c)

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)

		// Get the field tag value
		tag := field.Tag.Get("cmd")

		// Get the field value
		value := reflect.ValueOf(c).Field(i)

		// if value and tag not empty, write command
		if !value.IsZero() && tag != "" {
			switch x := value.Interface().(type) {
			case string:
				output = append(output, []string{tag, value.String()}...)
			case []string:
				items := value.Interface().([]string)
				content := tag+` "`
				for index, item := range items {
					if index == len(items)-1 {
						content = content + item + `"`
					} else {
						content = content + item + `,`
					}
				}
				output = append(output, content)
			default:
				fmt.Println(x)
			}
		}

	}
	return output
}
