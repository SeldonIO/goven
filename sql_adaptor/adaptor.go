// Package sql_adaptor provides functions to convert a goven query to a valid and safe SQL query.
package sql_adaptor

import (
	"reflect"
	"regexp"
	"strings"
)

func NewDefaultAdaptorFromStruct(gorm reflect.Value) *SqlAdaptor {
	matchers := map[*regexp.Regexp]ParseValidateFunc{}
	fieldMappings := map[string]string{}
	defaultFields := FieldParseValidatorFromStruct(gorm)
	return NewSqlAdaptor(fieldMappings, defaultFields, matchers)
}

// FieldParseValidatorFromStruct
// Don't panic - reflection is only used once on initialisation.
func FieldParseValidatorFromStruct(gorm reflect.Value) map[string]ParseValidateFunc {
	defaultFields := map[string]ParseValidateFunc{}
	e := gorm.Elem()

	for i := 0; i < e.NumField(); i++ {
		varName := strings.ToLower(e.Type().Field(i).Name)
		varType := e.Type().Field(i).Type
		vType := strings.TrimPrefix(varType.String(), "*")

		switch vType {
		case "float32", "float64":
			defaultFields[varName] = DefaultMatcherWithValidator(NumericValidator)
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
			defaultFields[varName] = DefaultMatcherWithValidator(IntegerValidator)
		default:
			defaultFields[varName] = DefaultMatcherWithValidator(NullValidator)
		}
	}
	return defaultFields
}
