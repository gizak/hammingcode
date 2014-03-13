package hammingcode

import "testing"

//import "github.com/davecgh/go-spew/spew"

func TestNewGraph(t *testing.T) {
	newGraph()
}

func TestGraphOperations(t *testing.T) {
	g := newGraph()
	fc0 := NewFactor([]int{0})
	fc1 := NewFactor([]int{0, 1})
	fc2 := NewFactor([]int{1})
	g.addVertex(fc0)
	g.addVertex(fc1)
	g.addVertex(fc2)

	g.setEdge(0, 1, NewFactor([]int{0}))
	g.setEdge(0, 2, NewFactor([]int{0, 1}))
	g.setEdge(0, 1, NewFactor([]int{1}))

	g.setVertex(0, fc2)
	//spew.Dump(g)
	if g.getVertex(0).scope[0] != 1 ||
		g.nodes[0].edges.next.nodeId != 2 ||
		g.nodes[0].edges.value.scope[0] != 1 {
		t.Fail()
	}
}
