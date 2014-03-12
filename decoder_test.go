package hammingcode

import "testing"
//import "github.com/davecgh/go-spew/spew"


func TestGenFactors (t *testing.T) {
	dc := NewDecoder()
	dc.noise = 1

	fc := dc.singletonFactor(0,-1)

	//spew.Dump(fc)
	if fc.data[0] + fc.data[1] != 1 {
		t.Fail()
	}
}

func TestDecoderInit (t *testing.T) {
	//dc := NewDecoder()
	cm := []int{0,0,0,1,1,1,1}
	fc := indicatorFactor(cm)
	//spew.Dump(fc)
	if fc.Get([]int{1,1,0,1}) != 0 || fc.Get([]int{1,1,0,0}) != 1 {
		t.Fail()
	}

	dc := NewDecoder()
	dc.noise = 1
	dc.rcvdCode = []float64{-1,-1,-1,-1,-1,-1,-1}

	dc.initAsCliqueTree([][]int{
		[]int{1,0,1,0,1,0,1},
		[]int{0,0,0,1,1,1,1},
		[]int{0,1,1,0,0,1,1},
	})

	//spew.Dump(dc.clusters.nodes[1])

	if len(dc.clusters.nodes) != 3 ||
		dc.clusters.nodes[1].edges == nil ||
		dc.clusters.nodes[1].edges.nodeId != 0 ||
		dc.clusters.nodes[1].edges.next.nodeId != 2 {
		t.Fail()
	}
}


func TestIsReady(t *testing.T) {
	dc := NewDecoder()
	dc.rcvdCode = []float64{-1,-1,-1,-1,-1,-1,-1}
	dc.initAsCliqueTree([][]int{
		[]int{1,0,1,0,1,0,1},
		[]int{0,0,0,1,1,1,1},
		[]int{0,1,1,0,0,1,1},
	})

	if dc.isReady(1,2) || dc.isReady(1,0) || !dc.isReady(0,1) {
		t.Fail()
	}

	dc.updateMsgs()
	dc.updateBelief()
	//spew.Dump(dc.clusters.nodes)
}


func TestDecode (t *testing.T) {
	dc := NewDecoder()
	dc.rcvdCode = []float64{1.2,1.1,0.8,1,0.7,1,1.2}
	dc.initAsCliqueTree([][]int{
		[]int{1,0,1,0,1,0,1},
		[]int{0,0,0,1,1,1,1},
		[]int{0,1,1,0,0,1,1},
	})


	//code := dc.Decode()
	//spew.Dump(code)

	/*
	dc.updateMsgs()
	dc.updateBelief()

	c0 := dc.clusters.getVertex(0).belief
	c1 := dc.clusters.getVertex(1).belief

	s := scpIntsc(c0.scope,c1.scope)

	fc0 := c0.sumOut(scpDiff(c0.scope,s))
	fc1 := c1.sumOut(scpDiff(c1.scope,s))

	spew.Dump(fc0,fc1)
	*/

}
