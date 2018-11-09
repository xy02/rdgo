package rdgo

import (
	"reflect"
	"strings"
)

type Field string

func (f Field) PathValue(target Data, tagKey string) Data {
	names := strings.Split(string(f), ".")
	v := reflect.ValueOf(target)
	v = pathValue(v, names, 0, tagKey)
	if v == (reflect.Value{}) {
		return v
	}
	r := v.Interface()
	return r
}

func pathValue(v reflect.Value, names []string, index int, tagKey string) reflect.Value {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		v = v.Elem()
	}
	// v = reflect.Indirect(v)
	name := names[index]
	switch v.Kind() {
	case reflect.Struct:
		if tagKey != "" {
			name = getFieldNameByTag(v, name, tagKey)
		}
		v = v.FieldByName(name)
	case reflect.Map:
		key := reflect.ValueOf(name)
		v = v.MapIndex(key)
	}
	next := index + 1
	if len(names) == next {
		return v
	}
	return pathValue(v, names, next, tagKey)
}

//it is not efficient if set TagKey
func getFieldNameByTag(v reflect.Value, name, tagKey string) string {
	vt := reflect.TypeOf(v.Interface())
	for i := 0; i < vt.NumField(); i++ {
		field := vt.Field(i)
		alias := field.Tag.Get(tagKey)
		if alias == name {
			return field.Name
		}
	}
	return name
}
