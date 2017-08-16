package util

import (
	"reflect"
	"sync"
)

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) ReflectWrap(f reflect.Value, p []reflect.Value) {
	w.Add(1)
	go func() {
		f.Call(p)
		w.Done()
	}()
}

func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
