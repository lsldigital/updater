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
	schema, err := schemaFromInstance(instance)
	if err != nil {
		return nil, err
	}

	return func(existing interface{}, values map[string]interface{}) (interface{}, error) {
		newElem, err := newElementFromStruct(existing)
		if err != nil {
			return nil, err
		}

		for name, propname := range schema {
			field := newElem.FieldByName(propname)
			updateField(name, propname, values, existing, &field)
		}

		return newElem.Interface(), nil
	}, nil
}

// updateField updates a field using either new or existing values
// if no new values found for field, use existing values from existing instance of object
func updateField(name, propname string, values map[string]interface{}, existing interface{}, field *reflect.Value) {
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
		if fieldDest := valOfExisting.FieldByName(propname); fieldDest.IsValid() {
			field.Set(fieldDest)
		}
	}
}

// schemaFromInstance accepts an instance of an object
// and returns map[fieldname]propname
func schemaFromInstance(instance interface{}) (map[string]string, error) {
	schema := make(map[string]string)

	valElem := reflect.ValueOf(instance)

	// if pointer get value it points to instead
	if valElem.Kind() == reflect.Ptr {
		valElem = valElem.Elem()
	}

	if valElem.Kind() != reflect.Struct {
		return nil, errors.New("instance must be of type struct")
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
		return nil, errors.New("derived schema is empty")
	}

	return schema, nil
}

// newElementFromStruct creates a new element from existing struct
func newElementFromStruct(existing interface{}) (*reflect.Value, error) {
	typeOfExisting := reflect.TypeOf(existing)

	// if pointer get value it points to instead
	if typeOfExisting.Kind() == reflect.Ptr {
		typeOfExisting = typeOfExisting.Elem()
	}

	if typeOfExisting.Kind() != reflect.Struct {
		return nil, errors.New("existing object must be of type struct")
	}

	newElem := reflect.New(typeOfExisting).Elem()
	if !newElem.CanInterface() {
		return nil, errors.New("new element from existing object cannot be casted to interface")
	}

	return &newElem, nil
}
