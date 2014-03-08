package hammingcode

import "testing"
import _"github.com/davecgh/go-spew/spew"


func TestData(t *testing.T) {
	//f := Factor{[]int{1,2},[]float64{0.1,0.2,0.3,0.4}}
	f1 := Factor{[]int{1,2},[]float64{1,1,1,1,1}}
	f0 := Factor{[]int{1,3},[]float64{1,1,1,2}}
	 FactorProduct(f1,f0)
	//s := setFromFactor(f)
	//t.Error(1)
	//spew.Dump(f)
}



















