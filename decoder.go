package hammingcode

import "math"

const MaxProduct = 1
const SumProduct = 0

type Decoder struct {
	clusters *graph
	mode     int
	rcvdCode []float64
	iterCnt  int
	chkMat   [][]int
	noise    float64
}

func NewDecoder() Decoder {
	return Decoder{}
}

func (dc Decoder) proirProb(z float64) (float64, float64) {
	foo := func(x int) float64 {
		return 1 / (1 + math.Exp((1-2*float64(x))*4*z/math.Pow(dc.noise, 2)))
	}
	return foo(0), foo(1)
}

func indicatorFactor(chk []int) Factor {
	scope := []int{}
	for _, v := range chk {
		if v == 1 {
			scope = append(scope, v)
		}
	}
	fc := NewFactor(scope)

	walk(len(scope), func(idx []int) {
		sum := 0
		for _, v := range idx {
			sum += v
		}
		fc.Set(idx, float64(sum%2))
	})

	return fc
}


func (dc Decoder) singletonFactor(idx int, z float64) Factor {
	x0, x1 := dc.proirProb(z)
	fc := Factor{[]int{idx},[]float64{x0,x1}}
	return fc
}


func (dc Decoder) Init() error {

	return nil
}


func (dc Decoder) initAsCliqueTree(chk [][]int) error {
	singletons := []Factor{}
	for i,v := range dc.rcvdCode {
		singletons = append(singletons,dc.singletonFactor(i,v))
	}

	indicators := []Factor{}
	for _,v := range chk {
		indicators = append(indicators, indicatorFactor(v))
	}

	// merge singletons into indicators
	for i,v := range singletons {
		for indi,indv := range indicators {
			for _, scpv := range indv.scope {
				if scpv == i {
					indicators[indi] = FactorProduct(indv,v)
					goto last
				}
			}
		last:
		}
	}
	for _,v := range indicators {
		dc.clusters.addVertex(v)
	}

	return nil
}
