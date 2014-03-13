// Copyleft 2014 winter. All rights you dont need to care about ;)
// Build for the course ELG5131 Graphical Models, Prof. Yongyi Mao

// Package hammingcode implements hamming (7,4) error control decoding task using graphical model method.
package hammingcode

import "math"
import "errors"

const SumProduct = 0        //mode
const MaxProduct = 1 << 1   //mode
const CliqueTree = 1 << 2   // arch
const LoopyCluster = 1 << 3 // arch

// Decoder stores all the graph info used to compute.
type Decoder struct {
	clusters *graph
	mode     int
	arch     int
	rcvdCode []float64
	iterCnt  int
	noise    float64
}

// NewDecoder returns a pointer pointing to Decoder.
func NewDecoder() *Decoder {
	return &Decoder{
		clusters: newGraph(),
		mode:     SumProduct,
		noise:    1,
		arch:     CliqueTree,
	}
}

// Calc P(x|z)
func (dc Decoder) proirProb(z float64) (float64, float64) {
	foo := func(x int) float64 {
		return 1 / (1 + math.Exp((1-2*float64(x))*4*z/math.Pow(dc.noise, 2)))
	}
	return foo(0), foo(1)
}

// Generate the membership function factor e.g. delta(x0 + x2 + x4 + x6).
func indicatorFactor(chk []int) Factor {
	scope := []int{}
	for i, v := range chk {
		if v == 1 {
			scope = append(scope, i)
		}
	}
	fc := NewFactor(scope)

	walk(len(scope), func(idx []int) {
		sum := 1
		for _, v := range idx {
			sum += v
		}
		fc.Set(idx, float64(sum%2))
	})

	return fc
}

// Generate singletons probability factor. e.g. F1 := P(x1|z1)
func (dc Decoder) singletonFactor(idx int, z float64) Factor {
	x0, x1 := dc.proirProb(z)
	fc := Factor{[]int{idx}, []float64{x0, x1}}
	return fc
}

// Init build the model structure, should be invoked after all neccessary vars being set
func (dc Decoder) Init(chk [][]int) error {
	switch dc.arch {

	case CliqueTree:
		return dc.initAsCliqueTree(chk)
	case LoopyCluster:
	default:
		errors.New("No designated cluster struct is assigned")
	}
	return nil
}

// Init as a clique tree
func (dc Decoder) initAsCliqueTree(chk [][]int) error {
	singletons := []Factor{}
	for i, v := range dc.rcvdCode {
		singletons = append(singletons, dc.singletonFactor(i, v))
	}

	indicators := []Factor{}
	for _, v := range chk {
		indicators = append(indicators, indicatorFactor(v))
	}

	// merge singletons into indicators
	for i, v := range singletons {
		for indi, indv := range indicators {
			for _, scpv := range indv.scope {
				if scpv == i {
					indicators[indi] = FactorProduct(indv, v)
					goto last
				}
			}
		}
	last:
	}

	// init vertex
	for _, v := range indicators {
		dc.clusters.addVertex(v)
	}

	// init edges
	dc.clusters.setEdge(0, 1, Factor{})
	dc.clusters.setEdge(1, 0, Factor{})
	dc.clusters.setEdge(1, 2, Factor{})
	dc.clusters.setEdge(2, 1, Factor{})

	return nil
}

// msg from i to j
func (dc Decoder) msg(i, j int) Factor {
	fc := dc.clusters.getVertex(i).Factor
	for nb := dc.clusters.nodes[i].edges; nb != nil; nb = nb.next {
		if nb.nodeId != j {
			fc = FactorProduct(fc, *dc.clusters.getEdge(nb.nodeId, i))
		}
	}

	rm := scpDiff(dc.clusters.nodes[i].scope, dc.clusters.nodes[j].scope)

	if dc.mode == MaxProduct {
		fc = fc.maxOut(rm)
	} else {
		fc = fc.sumOut(rm)
	}

	return fc
}

// Check whether msg_{i \to j} is ready
func (dc Decoder) isReady(i, j int) bool {

	for nb := dc.clusters.nodes[i].edges; nb != nil; nb = nb.next {
		if nb.nodeId != j {
			edge := dc.clusters.getEdge(nb.nodeId, i)
			if edge == nil || edge.data == nil {
				return false
			}
		}
	}
	return true
}

// Update all edges msgs iterativly until all msgs are computed
func (dc Decoder) updateMsgs() {
	status := make(map[[2]int]bool)

	for computable := true; computable; {
		computable = false
		for i, _ := range dc.clusters.nodes {
			for nb := dc.clusters.nodes[i].edges; nb != nil; nb = nb.next {
				j := nb.nodeId
				if dc.isReady(i, j) {
					if !status[[2]int{i, j}] {
						computable = true
						dc.clusters.setEdge(i, j, dc.msg(i, j))
						status[[2]int{i, j}] = true
					}
				}
			}
		}
	}
}

// Compute every vertex's belief, should be called after updateMsgs
func (dc Decoder) updateBelief() {
	for i, _ := range dc.clusters.nodes {
		fc := dc.clusters.getVertex(i).Factor
		for nb := dc.clusters.nodes[i].edges; nb != nil; nb = nb.next {
			fc = FactorProduct(*dc.clusters.getEdge(nb.nodeId, i), fc)
		}
		dc.clusters.getVertex(i).belief = fc
	}
}

// Decode when using clique tree struct
func (dc Decoder) decodeAsCliqueTree() []int {
	dc.updateMsgs()
	dc.updateBelief()
	code := make([]int, 7)
	foo := func(i int) Factor {
		for _, v := range dc.clusters.nodes {
			for _, vv := range v.scope {
				if vv == i {
					return v.sumOut(scpDiff(v.scope, []int{i}))
				}
			}
		}
		return Factor{}
	}

	// switch
	if dc.mode == MaxProduct {
		goto MAXP
	}

	for i, _ := range code {
		fc := foo(i)
		if fc.data[1] > fc.data[0] {
			code[i] = 1
		}
	}
	return code

MAXP:
	for _, v := range dc.clusters.nodes {
		// find max
		maxIdx := make([]int, 4)
		maxVal := -1.0
		walk(len(v.scope), func(idx []int) {
			if maxVal < v.belief.Get(idx) {
				maxVal = v.belief.Get(idx)
				copy(maxIdx, idx) // ! simply assign maxId = idx dosen't work, stack var issue?
			}
		})

		for i, _ := range maxIdx {
			code[v.scope[i]] = maxIdx[i]
		}
	}
	return code
}

// Decode decodes the code from a received code. The choice of method depends on the .arch property.
func (dc Decoder) Decode() []int {
	switch dc.arch {
	case CliqueTree:
		return dc.decodeAsCliqueTree()
	default:
		return nil //errors.New("No valid struct is assigned")
	}
}

// Set received code into Docoder. The code should be []float64 formatted.
func (dc *Decoder) SetRcvCode(c []float64) {
	dc.rcvdCode = c
}

// Set the method used in decoding (i.e. Max-Produt or Sum-Product)
func (dc *Decoder) SetMode(m int) {
	dc.mode = m
}

// Set the Structure used in decoding (i.e. Cycle-free: CliquTree or LoopyCluster)
func (dc *Decoder) SetArch(a int) {
	dc.arch = a
}

// Set the guassian noise's std dev
func (dc *Decoder) SetNoiseLevel(n float64) {
	dc.noise = n
}
