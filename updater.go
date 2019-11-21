package updater // import "go.lsl.digital/updater"

import (
	"errors"
	"reflect"
)

// Updater accepts an existing instance of an object (typically loaded from database)
// and values to update the object with.
// It returns an updated version of the object.
type Updater func(existing interface{}, values map[string]interface{}) (interface{}, error)

// New is a factory function that given an instance of an object will generate an "Updater" function.
func New(instance interface{}) (Updater, error) {
	schema := make(map[string]struct{})

	valElem := reflect.ValueOf(instance)

	if valElem.Kind() != reflect.Struct {
		return nil, errors.New("instance must be of type struct")
	}

	for i := 0; i < valElem.NumField(); i++ {
		typeField := valElem.Type().Field(i)
		name := typeField.Name
		schema[name] = struct{}{}
	}

	return func(existing interface{}, values map[string]interface{}) (interface{}, error) {
		typeOfExisting := reflect.TypeOf(existing)

		if typeOfExisting.Kind() != reflect.Struct {
			return nil, errors.New("existing object must be of type struct")
		}

		newElem := reflect.New(typeOfExisting).Elem()
		if !newElem.CanInterface() {
			return nil, errors.New("new element from existing object cannot be casted to interface")
		}

		for name := range schema {
			field := newElem.FieldByName(name)
			updateField(name, values, existing, &field)
		}

		return newElem.Interface(), nil
	}, nil
}

// updateField updates a field using either new or existing values
// if no new values found for field, use existing values from existing instance of object
func updateField(name string, values map[string]interface{}, existing interface{}, field *reflect.Value) {
	if !(field.IsValid() && field.CanSet()) {
		return
	}

	if raw, ok := values[name]; ok && raw != nil {
		valM := reflect.ValueOf(raw)
		if !valM.IsValid() {
			return
		}
		if t := field.Type(); valM.Type().ConvertibleTo(t) {
			if v := valM.Convert(t); v.IsValid() {
				field.Set(v)
			}
		}
	} else if valOfExisting := reflect.ValueOf(existing); valOfExisting.Kind() == reflect.Struct {
		if fieldDest := valOfExisting.FieldByName(name); fieldDest.IsValid() {
			field.Set(fieldDest)
		}
	}
}
