package rdgo

import (
	"errors"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Filter struct {
	TagKey string
	BasicHandler
	//key is handler
	handlerMap sync.Map
}

func (f *Filter) Input(data Data) error {
	if data == nil {
		return errors.New("data must not be nil")
	}
	f.SetLast(data)
	//broadcast
	f.handlerMap.Range(func(k, v interface{}) bool {
		handler := k.(Handler)
		condition := v.(Condition)
		if condition.Match(data, f.TagKey) {
			go handler.Input(data)
		}
		return true
	})
	return nil
}

func (f *Filter) Select(condition Condition, handler Handler) error {
	if handler == nil {
		return errors.New("handler must not be nil")
	}
	if _, ok := f.handlerMap.Load(handler); ok {
		return errors.New("handler is duplicated")
	}
	//ouput last data from memory
	lastData := f.GetLast()
	if lastData != nil && condition.Match(lastData, f.TagKey) {
		err := handler.Input(lastData)
		if err != nil {
			return err
		}
	}
	handler.OnDestroy(func() {
		f.Unselect(handler)
	})
	f.handlerMap.Store(handler, condition)
	// f.handlerMap.Range(func(k, v interface{}) bool {
	// 	log.Printf("%p\n", k)
	// 	return true
	// })
	return nil
}

func (f *Filter) Unselect(handler Handler) {
	f.handlerMap.Delete(handler)
}

func (f *Filter) Destroy() {
	f.handlerMap.Range(func(k, v interface{}) bool {
		f.handlerMap.Delete(k)
		handler := v.(Handler)
		go handler.Destroy()
		return true
	})
}

type (
	Condition map[Field]Exp

	//Exp : Expression
	Exp map[Operator]Data

	Field string

	Operator string
)

func (c Condition) Match(target Data, tagKey string) bool {
	for field, exp := range c {
		if !exp.Match(target, field, tagKey) {
			return false
		}
	}
	return true
}

func (exp Exp) Match(target Data, f Field, tagKey string) bool {
	value := f.PathValue(target, tagKey)
	for op, v := range exp {
		if !op.Exec(value, v) {
			return false
		}
	}
	return true
}

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

//Operators
const (
	EQ  Operator = "$eq"
	GT  Operator = "$gt"
	GTE Operator = "$gte"
	LT  Operator = "$lt"
	LTE Operator = "$lte"
)

//Exec : execute
func (op Operator) Exec(target, v Data) bool {
	fn := opMap[op]
	return fn(target, v)
}

//Is func returns an eq expression
func Is(v Data) Exp {
	return Exp{EQ: v}
}

//Gt func
func Gt(v Data) Exp {
	return Exp{GT: v}
}

//Gte func
func Gte(v Data) Exp {
	return Exp{GTE: v}
}

//Lt func
func Lt(v Data) Exp {
	return Exp{LT: v}
}

//Lte func
func Lte(v Data) Exp {
	return Exp{LTE: v}
}

var opMap = map[Operator]func(target, v Data) bool{
	EQ: func(target, v Data) bool {
		return target == v
	},
	GT: func(target, v Data) bool {
		target = parseNumber(target)
		v = parseNumber(v)
		switch vt := v.(type) {
		case int64:
			switch tt := target.(type) {
			case int64:
				return tt > vt
			case float64:
				return tt > float64(vt)
			}
		case float64:
			switch tt := target.(type) {
			case int64:
				return float64(tt) > vt
			case float64:
				return tt > vt
			}
		}
		return false
	},
	GTE: func(target, v Data) bool {
		target = parseNumber(target)
		v = parseNumber(v)
		switch vt := v.(type) {
		case int64:
			switch tt := target.(type) {
			case int64:
				return tt >= vt
			case float64:
				return tt >= float64(vt)
			}
		case float64:
			switch tt := target.(type) {
			case int64:
				return float64(tt) >= vt
			case float64:
				return tt >= vt
			}
		}
		return false
	},
	LT: func(target, v Data) bool {
		target = parseNumber(target)
		v = parseNumber(v)
		switch vt := v.(type) {
		case int64:
			switch tt := target.(type) {
			case int64:
				return tt < vt
			case float64:
				return tt < float64(vt)
			}
		case float64:
			switch tt := target.(type) {
			case int64:
				return float64(tt) < vt
			case float64:
				return tt < vt
			}
		}
		return false
	},
	LTE: func(target, v Data) bool {
		target = parseNumber(target)
		v = parseNumber(v)
		switch vt := v.(type) {
		case int64:
			switch tt := target.(type) {
			case int64:
				return tt <= vt
			case float64:
				return tt <= float64(vt)
			}
		case float64:
			switch tt := target.(type) {
			case int64:
				return float64(tt) <= vt
			case float64:
				return tt <= vt
			}
		}
		return false
	},
}

func parseNumber(v Data) interface{} {
	switch vt := v.(type) {
	case time.Time:
		return vt.Unix()
	case int:
		return int64(vt)
	case int8:
		return int64(vt)
	case int16:
		return int64(vt)
	case int32:
		return int64(vt)
	case uint:
		return int64(vt)
	case uint8:
		return int64(vt)
	case uint16:
		return int64(vt)
	case uint32:
		return int64(vt)
	case uint64:
		return int64(vt)
	case float32:
		return float64(vt)
	}
	return v
}
