package hammingcode

import "testing"
import "github.com/davecgh/go-spew/spew"

var dc *Decoder

func TestNew(t *testing.T) {
	H := [][]int{
		[]int{0,0,0,1,1,1,1},
		[]int{0,1,1,0,0,1,1},
		[]int{1,0,1,0,1,0,1},
	}
	dc = NewDecoder(MaxProduct)
	dc.SetNoiseLevel(1)
	dc.SetCheckMatrix(H)
	dc.Accept([]float64{-1,-1,-1,-1,-1,-1,-1})
	dc.InitCliqueTree()
	spew.Dump(dc.clusters[1].neighbours[0].scope)
}

func TestElOut(t *testing.T) {
	//fc := Factor{[]int{0,1},[]float64{1,10,0,1}}
	//nfc := elOut(fc, []int{0}, MaxProduct)
	//spew.Dump(nfc)
	/*
	res := []int{}
	
	dc.Next()

	//spew.Dump(dc.edges[6][7])
	for i:=0; i<7; i++ {
		if dc.clusters[i].belief.data[0] > dc.clusters[i].belief.data[1] {
			res = append(res,0)
		} else {
			res = append(res,1)
		}
	}
	spew.Dump(res)
*/
}

