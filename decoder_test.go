package hammingcode

import "testing"
import "github.com/davecgh/go-spew/spew"


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
	if fc.Get([]int{1,1,0,1}) != 1 || fc.Get([]int{1,1,0,0}) != 0 {
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

	//spew.Dump(dc.clusters.nodes[1].edges)

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

	//spew.Dump(dc.clusters.nodes)
}


func TestDecode (t *testing.T) {
	dc := NewDecoder()
	dc.rcvdCode = []float64{-1.3,-0.8,-0.9,-1,-1,-1.2,-1}
	dc.initAsCliqueTree([][]int{
		[]int{1,0,1,0,1,0,1},
		[]int{0,0,0,1,1,1,1},
		[]int{0,1,1,0,0,1,1},
	})

	code := dc.Decode()
	spew.Dump(code)
}
