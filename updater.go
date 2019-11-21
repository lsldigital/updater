package updater // import "go.lsl.digital/updater"

import (
	"reflect"
)

// Updater accepts an existing object (typically loaded from database)
// It returns an updated version of the object.
type Updater func(existing interface{}, values map[string]interface{}) interface{}

// New is a factory function that given an instance of an object will generate an updater function.
func New(instance interface{}) Updater {
	schema := make(map[string]struct{})

	valElem := reflect.ValueOf(instance)

	if valElem.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < valElem.NumField(); i++ {
		typeField := valElem.Type().Field(i)
		name := typeField.Name
		schema[name] = struct{}{}
	}

	return func(existing interface{}, values map[string]interface{}) interface{} {
		typeOfExisting := reflect.TypeOf(existing)

		if typeOfExisting.Kind() != reflect.Struct {
			return nil
		}

		newElem := reflect.New(typeOfExisting).Elem()
		if !newElem.CanInterface() {
			return nil
		}

		for name := range schema {
			field := newElem.FieldByName(name)
			updateField(name, values, existing, &field)
		}

		return newElem.Interface()
	}
}

// updateField updates a field
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
