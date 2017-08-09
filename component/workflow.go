package component

type workflow struct {
	start func(interface{}) interface{}
	end   func(interface{}) interface{}
	nodes []func(interface{}) interface{}
}

func (w *workflow) Add(item func(interface{}) interface{}) {
	w.nodes = append(w.nodes, item)
}

func (w *workflow) generate() {

}
