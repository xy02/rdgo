package rdgo

import (
	"errors"
	"sync"
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
		expr := v.(Expr)
		if expr.Match(data, f) {
			go handler.Input(data)
		}
		return true
	})
	return nil
}
func (f *Filter) Select(expr Expr, handler Handler) error {
	if handler == nil {
		return errors.New("handler must not be nil")
	}
	if _, ok := f.handlerMap.Load(handler); ok {
		return errors.New("handler is duplicated")
	}
	//ouput last data from memory
	lastData := f.GetLast()
	if lastData != nil && expr.Match(lastData, f) {
		err := handler.Input(lastData)
		if err != nil {
			return err
		}
	}
	handler.OnDestroy(func() {
		f.Unselect(handler)
	})
	f.handlerMap.Store(handler, expr)
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
