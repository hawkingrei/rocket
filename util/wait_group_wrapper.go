package util

import (
	"reflect"
	"sync"
)

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(f reflect.Value, p []reflect.Value) {
	w.Add(1)
	go func() {
		f.Call(p)
		w.Done()
	}()
}
