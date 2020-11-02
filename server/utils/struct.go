package utils

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	indentChars = "─"
	lChar       = "└"
	vertChar    = "│"
	horizChar   = "─"
	branchChar  = "├"
)

// PrintStruct prints the tags
func PrintStruct(o interface{}) {
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)
	if v.Kind() != reflect.Struct {
		fmt.Println("Not a struct")
		return
	}
	printSubTags(v, t, reflect.StructField{}, []int{}, t.Name())
}

func printSubTags(val reflect.Value, t reflect.Type, field reflect.StructField, index []int, parentName string) {
	for x := 0; x < val.NumField(); x++ {
		var findex []int
		if len(index) == 0 {
			// Top level => fieild is empty
			findex = append(field.Index, x)
		} else {
			findex = append(index, x)
		}
		sfield := t.FieldByIndex(findex)
		sval := val.Field(x)
		var parname string = parentName
		if field.Name != "" {
			parname = fmt.Sprintf("%s.%s", parentName, field.Name)
		}
		indent := strings.Repeat(indentChars, strings.Count(parname, "."))
		if indent != "" {
			fmt.Println(branchChar+horizChar+horizChar, "Sub Field", parname+"."+sfield.Name)
		} else {
			fmt.Println(branchChar, "Top Field", parname+"."+sfield.Name)
		}
		switch sval.Kind() {
		case reflect.Struct:
			printSubTags(sval, t, sfield, findex, parname)

		case reflect.Slice, reflect.Array:
			// switch sfield.Type.Elem().Kind(){
			// case reflect.Slice:

			// }
			indent += indentChars
			fmt.Println(indent, "sub key", "[]"+sfield.Type.Elem().Name())
			if typeVal := sfield.Tag.Get("type"); typeVal != "" {
				fmt.Println(indent+indentChars, "type:", typeVal)
			}
			if sval.Len() > 0 {
				fmt.Println(indent, "[")
				for x := 0; x < sval.Len(); x++ {
					fmt.Println(indent+indentChars, sval.Index(x))
				}
				fmt.Println(indent, "]")
			} else {
				fmt.Println(indent, "[]")
			}

		default:
			indent += indentChars
			fmt.Println(indent, "sub key", sval.Kind())
			if typeVal := sfield.Tag.Get("type"); typeVal != "" {
				fmt.Println(indent+indentChars, "type:", typeVal)
			}
		}
	}
}