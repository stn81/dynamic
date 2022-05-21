package dynamic

import (
	"encoding/json"
	"reflect"
)

// ParseJSON with dynamic field support
func ParseJSON(data []byte, ptr interface{}) (err error) {
	if err = json.Unmarshal(data, ptr); err != nil {
		return err
	}
	return parseDynamic(reflect.ValueOf(ptr), nil, "")
}

func parseDynamic(v reflect.Value, dynFielder DynamicFielder, dynFieldName string) (err error) {
	switch v.Kind() {
	case reflect.Interface:
		if v.IsNil() {
			return nil
		}
		return parseDynamic(v.Elem(), dynFielder, dynFieldName)
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		if v.Type() != DynamicType {
			return parseDynamic(v.Elem(), dynFielder, dynFieldName)
		}
		dynValue := v.Interface().(*Type)
		if len(dynValue.raw) > 0 {
			if dynFielder != nil {
				ptr := dynFielder.NewDynamicField(dynFieldName)
				if ptr != nil {
					if err = ParseJSON(dynValue.raw, ptr); err != nil {
						return err
					}
					dynValue.Value = ptr
					return nil
				}
			}
		}
		if v.CanSet() {
			v.Set(reflect.Zero(v.Type()))
		}
		return nil
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if err = parseDynamic(v.Index(i), dynFielder, dynFieldName); err != nil {
				return err
			}
		}
	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			if err = parseDynamic(iter.Value(), dynFielder, dynFieldName); err != nil {
				return err
			}
		}
	case reflect.Struct:
		if !v.CanAddr() {
			return nil
		}

		dynFielder, _ = v.Addr().Interface().(DynamicFielder)
		typ := v.Type()
		for i := 0; i < v.NumField(); i++ {
			sf := typ.Field(i)
			field := v.Field(i)

			if !field.CanSet() {
				continue
			}

			if err = parseDynamic(field, dynFielder, sf.Name); err != nil {
				return err
			}
		}
	}
	return nil
}
