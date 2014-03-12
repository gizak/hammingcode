package hammingcode

import "testing"
//import "github.com/davecgh/go-spew/spew"



func TestData(t *testing.T) {
	fc := NewFactor([]int{1,2,3})
	if fc.Get([]int{0,1,0}) != 1 || len(fc.scope) != 3 {
		t.Fail()
	}
}

func TestFactorOperations(t *testing.T) {
	fc0 := NewFactor([]int{1,2})
	fc0.Set([]int{1,0},2)
	fc1 := NewFactor([]int{2,3})
	fc1.Set([]int{0,1},3)
	fc2 := FactorProduct(fc0,fc1)
	//spew.Dump(fc2)
	if fc2.Get([]int{1,1,0}) != 1 ||
		fc2.Get([]int{0,0,1}) != 3 ||
		fc2.Get([]int{1,0,0}) != 2 {
		t.Fail()
	}

	// chk with len(scope) = 0 (constant)
	fc3 := NewFactor([]int{})
	fc4 := FactorProduct(fc0,fc3)
	if fc4.Get([]int{1,0}) != 2 ||
		fc4.Get([]int{0,0}) != 1 {
		t.Fail()
	}
}

func TestSumOut (t *testing.T) {
	fc := NewFactor([]int{1,2})
	fc.Set([]int{0,0},2)
	nfc := fc.sumOut([]int{1})
	if nfc.Get([]int{0}) != 3 || nfc.Get([]int{1}) != 2 {
		t.Fail()
	}
	nfc = fc.maxOut([]int{2})
	if nfc.Get([]int{0}) != 2 || nfc.Get([]int{1}) != 1 {
		t.Fail()
	}
}
