package main

import (
	"fmt"
	"reflect"

	"github.com/hawkingrei/emitter/component"
	"github.com/hawkingrei/emitter/util/inject"
)

type SpecialString interface{}
type Year interface{}

type Other struct {
	Id int
}

func (o *Other) One(output component.Output) {
	out := reflect.ValueOf(output)
	for a := range []int{2, 3, 4} {
		fmt.Println("one ", a)
		v := reflect.ValueOf(a)
		out.Send(v)
	}
}

func (o *Other) Two(a int) {
	fmt.Println(a + o.Id)
}

func main() {
	wf := component.NewWorkflow()
	one := Other{Id: 1}
	inj := inject.New()
	wf.Add(one.One, inj, reflect.TypeOf(1))
	wf.Add(one.Two, inj, reflect.TypeOf(1))
	wf.Run()
}
