package rdgo

import (
	"errors"
	"log"
	"reflect"
	"strings"
	"sync"
)

type Filter struct {
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
		if condition.Match(data) {
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
	if lastData != nil && condition.Match(lastData) {
		err := handler.Input(lastData)
		if err != nil {
			return err
		}
	}
	handler.OnDestroy(func() {
		f.Unselect(handler)
	})
	f.handlerMap.Store(handler, condition)
	f.handlerMap.Range(func(k, v interface{}) bool {
		log.Printf("%p\n", k)
		return true
	})
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

	//Exp means Expression
	Exp map[Operator]Data

	Field string

	Operator string
)

func (c Condition) Match(target Data) bool {
	for field, exp := range c {
		if !exp.Match(target, field) {
			return false
		}
	}
	return true
}

func (exp Exp) Match(target Data, f Field) bool {
	value := f.PathValue(target)
	for op, v := range exp {
		if !op.Exec(value, v) {
			return false
		}
	}
	return true
}

func (f Field) PathValue(target Data) Data {
	names := strings.Split(string(f), ".")
	v := reflect.ValueOf(target)
	v = pathValue(v, names, 0)
	if v == (reflect.Value{}) {
		return v
	}
	r := v.Interface()
	return r
}

func pathValue(v reflect.Value, names []string, index int) reflect.Value {
	v = reflect.Indirect(v)
	name := names[index]
	switch v.Kind() {
	case reflect.Struct:
		v = v.FieldByName(name)
	case reflect.Map:
		key := reflect.ValueOf(name)
		v = v.MapIndex(key)
	}
	next := index + 1
	if len(names) == next {
		return v
	}
	return pathValue(v, names, next)
}

const (
	Eq  Operator = "$eq"
	Gt  Operator = "$gt"
	Gte Operator = "$gte"
	Lt  Operator = "$lt"
	Lte Operator = "$lte"
)

var opMap = map[Operator]func(target, v Data) bool{
	Eq: EqFn,
}

func (op Operator) Exec(target, v Data) bool {
	fn := opMap[op]
	return fn(target, v)
}

func EqFn(target, v Data) bool {
	return target == v
}

//Is func returns an eq expression
func Is(v Data) Exp {
	return Exp{Eq: v}
}
