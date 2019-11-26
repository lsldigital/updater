package updater // import "go.lsl.digital/updater"

import (
	"errors"
	"reflect"
)

// Common errors for updater
var (
	ErrInvalidInstance = errors.New("instance must be of type struct")
	ErrInvalidExisting = errors.New("existing must be of type struct")
	ErrEmptyInstance   = errors.New("instance must have at least one field")
	ErrInvalidValType  = errors.New("invalid value type")
)

// Field represents information about a single field of a struct
type Field struct {
	Name     string
	Propname string
	Type     reflect.Type
}

// Updater accepts a pointer to an object (typically loaded from database)
// and values to update the object with.
// It updates the existing object with the values.
type Updater func(existing interface{}, values map[string]interface{}) (interface{}, error)

// New is a factory function that given an instance of an object will generate an "Updater" function.
func New(instance interface{}) (Updater, error) {
	fields, instanceType, err := metaFromInstance(instance)
	if err != nil {
		return nil, err
	}

	return func(existing interface{}, values map[string]interface{}) (interface{}, error) {
		newValOfInstance := valFromType(instanceType)

		valOfExisting, ok := valOfStruct(existing)
		if !ok {
			return nil, ErrInvalidExisting
		}

		for index, field := range fields {
			val := getValue(field, values, valOfExisting)
			if val.IsZero() {
				continue
			}

			setField(index, newValOfInstance, val)
		}

		return newValOfInstance.Interface(), nil
	}, nil
}

// getValue returns a value from either values or field in valOfStruct
// else we return zero value of field
func getValue(field Field, values map[string]interface{}, valOfStruct reflect.Value) reflect.Value {
	// if in values, use it
	if raw, ok := values[field.Name]; ok && raw != nil {
		valOfRaw := reflect.ValueOf(raw)
		if convertibleTo(valOfRaw, field.Type) {
			return valOfRaw
		}
	}

	// else use value from field in existing struct
	fieldExisting := valOfStruct.FieldByName(field.Propname)
	if !fieldExisting.IsValid() || !convertibleTo(fieldExisting, field.Type) {
		return reflect.Zero(field.Type)
	}

	return fieldExisting
}

// convertibleTo checks if val can convert to fieldType
func convertibleTo(val reflect.Value, fieldType reflect.Type) bool {
	return val.Type().ConvertibleTo(fieldType)
}

// setField updates field of valOfStruct at specified index with newVal
func setField(index int, valOfStruct, newVal reflect.Value) {
	valOfStruct.Field(index).Set(newVal)
}

// metaFromInstance accepts an instance of an object
// and returns []Field
func metaFromInstance(instance interface{}) ([]Field, reflect.Type, error) {
	valOfInstance, ok := valOfStruct(instance)
	if !ok {
		return nil, nil, ErrInvalidInstance
	}

	numFields := valOfInstance.NumField()

	if numFields == 0 {
		return nil, nil, ErrEmptyInstance
	}

	schema := make([]Field, numFields, numFields)

	for i := 0; i < numFields; i++ {
		typeField := valOfInstance.Type().Field(i)

		propname := typeField.Name
		fieldname := typeField.Tag.Get("json")
		if fieldname == "" || fieldname == "-" {
			fieldname = toSnakeCase(propname)
		}

		schema[i] = Field{
			Name:     fieldname,
			Propname: propname,
			Type:     typeField.Type,
		}
	}

	return schema, valOfInstance.Type(), nil
}

// valOfStruct checks that instance is a struct
// and returns its reflect.Value
func valOfStruct(instance interface{}) (reflect.Value, bool) {
	valOfInstance := reflect.ValueOf(instance)

	// if pointer get value it points to instead
	valOfInstance = reflect.Indirect(valOfInstance)

	if valOfInstance.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}

	return valOfInstance, true
}

// valFromType returns the reflect.Value from structType
func valFromType(structType reflect.Type) reflect.Value {
	return reflect.New(structType).Elem()
}
