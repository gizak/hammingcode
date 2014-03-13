package hammingcode

import "testing"

//import "github.com/davecgh/go-spew/spew"

func TestGenFactors(t *testing.T) {
	dc := NewDecoder()
	dc.noise = 1

	fc := dc.singletonFactor(0, -1)

	//spew.Dump(fc)
	if fc.data[0]+fc.data[1] != 1 {
		t.Fail()
	}
}

func TestDecoderInit(t *testing.T) {
	//dc := NewDecoder()
	cm := []int{0, 0, 0, 1, 1, 1, 1}
	fc := indicatorFactor(cm)
	//spew.Dump(fc)
	if fc.Get([]int{1, 1, 0, 1}) != 0 || fc.Get([]int{1, 1, 0, 0}) != 1 {
		t.Fail()
	}

	dc := NewDecoder()
	dc.noise = 1
	dc.rcvdCode = []float64{-1, -1, -1, -1, -1, -1, -1}

	dc.initAsCliqueTree([][]int{
		[]int{1, 0, 1, 0, 1, 0, 1},
		[]int{0, 0, 0, 1, 1, 1, 1},
		[]int{0, 1, 1, 0, 0, 1, 1},
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
	dc.rcvdCode = []float64{-1, -1, -1, -1, -1, -1, -1}
	dc.initAsCliqueTree([][]int{
		[]int{1, 0, 1, 0, 1, 0, 1},
		[]int{0, 0, 0, 1, 1, 1, 1},
		[]int{0, 1, 1, 0, 0, 1, 1},
	})

	if dc.isReady(1, 2) || dc.isReady(1, 0) || !dc.isReady(0, 1) {
		t.Fail()
	}

	dc.updateMsgs()
	dc.updateBelief()
	//spew.Dump(dc.clusters.nodes)
}

func TestSumProdDecode(t *testing.T) {
	dc := NewDecoder()
	dc.rcvdCode = []float64{1.2, 1.1, 0.8, 1, 0.7, 1, 1.2}
	dc.initAsCliqueTree([][]int{
		[]int{1, 0, 1, 0, 1, 0, 1},
		[]int{0, 0, 0, 1, 1, 1, 1},
		[]int{0, 1, 1, 0, 0, 1, 1},
	})

	//spew.Dump(dc.Decode())

	// adjacent belif agreement
	dc.Decode()

	fc0 := dc.clusters.getVertex(0)
	fc1 := dc.clusters.getVertex(1)

	s := scpIntsc(fc0.scope, fc1.scope)

	mu0 := fc0.belief.sumOut(scpDiff(fc0.scope, s))
	mu1 := fc1.belief.sumOut(scpDiff(fc1.scope, s))

	for i, _ := range mu0.data {
		delta := mu0.data[i] - mu1.data[i]
		if -1e-8 > delta || delta > 1e-8 {
			t.Fail()
		}
	}
}

func TestMaxProdDecode(t *testing.T) {
	dc := NewDecoder()
	dc.SetArch(CliqueTree)
	dc.SetMode(MaxProduct)
	dc.SetRcvCode([]float64{1.2, 1.1, 0.8, 1, 0.7, 1, 1.2})
	dc.SetNoiseLevel(1)
	dc.Init([][]int{
		[]int{1, 0, 1, 0, 1, 0, 1},
		[]int{0, 0, 0, 1, 1, 1, 1},
		[]int{0, 1, 1, 0, 0, 1, 1},
	})

	// adjacent belif agreement
	dc.Decode()

	fc0 := dc.clusters.getVertex(0)
	fc1 := dc.clusters.getVertex(1)

	//spew.Dump(fc0, fc1)

	s := scpIntsc(fc0.scope, fc1.scope)

	mu0 := fc0.belief.maxOut(scpDiff(fc0.scope, s))
	mu1 := fc1.belief.maxOut(scpDiff(fc1.scope, s))

	for i, _ := range mu0.data {
		delta := mu0.data[i] - mu1.data[i]
		if -1e-8 > delta || delta > 1e-8 {
			t.Fail()
		}
	}
}
