package component

import (
	"fmt"
	"reflect"

	"github.com/hawkingrei/emitter/util"
	"github.com/hawkingrei/emitter/util/inject"
)

type Node struct {
	Function  interface{}
	Parameter inject.Injector
	ChanType  reflect.Type
}

type Input interface{}
type Output interface{}

//type InputChan reflect.Value
//type OutputChan reflect.Value

//func InputOutputNode()

type Workflow struct {
	channels []chan interface{}
	nodes    []Node
	wg       util.WaitGroupWrapper
}

func (w *Workflow) OutputNode(output reflect.Value, node Node) {

	var n Node = node
	n.Parameter.Map(output.Interface())
	n.Parameter.Invoke(n.Function)
}

func (w *Workflow) InputNode(input reflect.Value, node Node) {
	selectcase := make([]reflect.SelectCase, 1)
	selectcase[0].Chan = input
	selectcase[0].Dir = reflect.SelectRecv
	for {
		chosen, recv, recvOk := reflect.Select(selectcase)
		if recvOk {
			switch chosen {
			case 0:
				var n Node = node
				n.Parameter.Map(recv.Interface())
				n.Parameter.Invoke(n.Function)
			}
		} else {
			break
		}
	}
}

func (w *Workflow) InputOnputNode(input reflect.Value, output reflect.Value, node Node) {
	selectcase := make([]reflect.SelectCase, 2)
	selectcase[0].Chan = input
	selectcase[0].Dir = reflect.SelectDefault
	for {
		chosen, recv, recvOk := reflect.Select(selectcase)
		if recvOk {
			switch chosen {
			case 0:
				var n Node = node
				n.Parameter.Map(recv.Interface())
				n.Parameter.Invoke(n.Function)
				val, err := node.Parameter.Invoke(node.Function)
				if err != nil {
					fmt.Println(err.Error())
				}
				output.Send(val[0])

			}
		} else {
			break
		}
	}
}

func NewWorkflow() Workflow {
	w := Workflow{}
	return w
}

func makeChannel(t reflect.Type, chanDir reflect.ChanDir, buffer int) reflect.Value {
	ctype := reflect.ChanOf(chanDir, t)
	return reflect.MakeChan(ctype, buffer)
}

func (w *Workflow) Add(f interface{}, p inject.Injector, t reflect.Type) {
	//if p.HasType((*Input)(nil)) || p.HasType((*Output)(nil)) {
	w.nodes = append(w.nodes, Node{Function: f, Parameter: p, ChanType: t})
	//	return
	//}
	//panic("lock of Input/Output parameter.")
}

func (w *Workflow) assembly(begin int, process int, end int, input reflect.Value) {
	n := w.nodes[process]
	if begin == process {
		if !n.Parameter.HasType((*Output)(nil)) && n.Parameter.HasType((*Input)(nil)) {
			panic("first workflow node need Output,but not Input")
		}
		output := makeChannel(n.ChanType, reflect.BothDir, 1)
		w.wg.Wrap(func() {
			w.OutputNode(output, n)
			output.Close()
		})
		w.assembly(begin, process+1, end, output)
	} else if process == end {
		if !n.Parameter.HasType((*Input)(nil)) && n.Parameter.HasType((*Output)(nil)) {
			panic("first workflow node need Input,but not Output")
		}
		w.wg.Wrap(func() { w.InputNode(input, n) })
		return
	} else {
		input := makeChannel(n.ChanType, reflect.BothDir, 1)
		output := makeChannel(n.ChanType, reflect.BothDir, 1)
		w.wg.Wrap(func() { w.InputOnputNode(output, input, n) })
		w.assembly(begin, process+1, end, output)
	}
}

func (w *Workflow) Run() {
	w.assembly(0, 0, len(w.nodes)-1, reflect.ValueOf((*Output)(nil)))
	w.wg.Wait()
}
