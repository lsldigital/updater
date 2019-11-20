package updater // import "go.lsl.digital/updater"

import (
	"reflect"
)

// Updater accepts an instance of an object (typically loaded from database)
// and graphql resolve params.
// It returns an updated version of the object.
type Updater func(values map[string]interface{}, dest interface{}) interface{}

// New is a factory function that given a schematic will generate an updater function
// which itself accepts an instance of an object (typically loaded from database)
// and graphql resolve params.
// It returns an updated version of the object.
func New(element interface{}) Updater {
	schema := make(map[string]reflect.StructField)

	//TODO use element to create schema
	valElem := reflect.ValueOf(element)

	if valElem.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < valElem.NumField(); i++ {
		typeField := valElem.Type().Field(i)
		name := typeField.Name
		schema[name] = typeField
	}

	return func(values map[string]interface{}, dest interface{}) interface{} {
		typeDest := reflect.TypeOf(dest)

		if typeDest.Kind() != reflect.Struct {
			return nil
		}

		newEl := reflect.New(typeDest).Elem()
		if !newEl.CanInterface() {
			return nil
		}

		for name, _ := range schema {
			fieldUpdater(name, name, values, dest, &newEl)
		}

		return newEl.Interface()
	}
}

// fieldUpdater is a factory function for field
func fieldUpdater(name, propname string, values map[string]interface{}, dest interface{}, newEl *reflect.Value) {
	newField := newEl.FieldByName(propname)
	if !(newField.IsValid() && newField.CanSet()) {
		return
	}

	if raw, ok := values[name]; ok && raw != nil {
		valM := reflect.ValueOf(raw)
		if !valM.IsValid() {
			return
		}
		if t := newField.Type(); valM.Type().ConvertibleTo(t) {
			if v := valM.Convert(t); v.IsValid() {
				newField.Set(v)
			}
		}
	} else if valDest := reflect.ValueOf(dest); valDest.Kind() == reflect.Struct {
		if fieldDest := valDest.FieldByName(propname); fieldDest.IsValid() {
			newField.Set(fieldDest)
		}
	}
}
