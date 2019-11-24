package updater // import "go.lsl.digital/updater"

import (
	"errors"
	"reflect"
)

// Common errors for updater
var (
	ErrInvalidInstance      = errors.New("instance must be of type struct")
	ErrInvalidExistingObj   = errors.New("existing object must be of type pointer to struct")
	ErrInvalidDerivedSchema = errors.New("derived schema is empty")
)

// Updater accepts a pointer to an object (typically loaded from database)
// and values to update the object with.
// It updates the existing object with the values.
type Updater func(existing interface{}, values map[string]interface{}) error

// New is a factory function that given an instance of an object will generate an "Updater" function.
func New(instance interface{}) (Updater, error) {
	schema, err := schemaFromInstance(instance)
	if err != nil {
		return nil, err
	}

	return func(existing interface{}, values map[string]interface{}) error {
		valOfExisting := reflect.ValueOf(existing)

		// existing must a pointer to struct
		if valOfExisting.Kind() != reflect.Ptr {
			return ErrInvalidExistingObj
		}
		valOfExisting = valOfExisting.Elem()
		if valOfExisting.Kind() != reflect.Struct {
			return ErrInvalidExistingObj
		}

		for name, propname := range schema {
			updateField(name, propname, values, valOfExisting)
		}

		existing = valOfExisting.Interface()

		return nil
	}, nil
}

// updateField updates a field using either new or existing values
// if no new values found for field, use existing values from existing instance of object
func updateField(name, propname string, values map[string]interface{}, valOfExisting reflect.Value) {
	// get raw from values and check if valid
	raw, ok := values[name]
	if !ok || raw == nil {
		return
	}

	field := valOfExisting.FieldByName(propname)

	// field must be valid and settable
	if !field.IsValid() || !field.CanSet() {
		return
	}

	valOfRaw := reflect.ValueOf(raw)
	fieldType := field.Type()
	if !valOfRaw.Type().ConvertibleTo(fieldType) {
		return
	}

	field.Set(valOfRaw.Convert(fieldType))
}

// schemaFromInstance accepts an instance of an object
// and returns map[fieldname]propname
func schemaFromInstance(instance interface{}) (map[string]string, error) {
	schema := make(map[string]string)

	valElem := reflect.ValueOf(instance)

	// if pointer get value it points to instead
	valElem = reflect.Indirect(valElem)

	if valElem.Kind() != reflect.Struct {
		return nil, ErrInvalidInstance
	}

	for i := 0; i < valElem.NumField(); i++ {
		typeField := valElem.Type().Field(i)
		propname := typeField.Name
		fieldname := typeField.Tag.Get("json")
		if fieldname == "" || fieldname == "-" {
			fieldname = toSnakeCase(propname)
		}
		schema[fieldname] = propname
	}

	if len(schema) == 0 {
		return nil, ErrInvalidDerivedSchema
	}

	return schema, nil
}
