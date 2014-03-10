package hammingcode

import "testing"
import "github.com/davecgh/go-spew/spew"

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

	spew.Dump(dc.clusters)
}


func TestGenFactors (t *testing.T) {
	dc := NewDecoder()
	dc.noise = 1

	fc := dc.singletonFactor(0,-1)
	//spew.Dump(fc)
	if fc.data[0] + fc.data[1] != 1 {
		t.Fail()
	}
}
