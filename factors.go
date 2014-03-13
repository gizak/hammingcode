package hammingcode

import "sort"
import "github.com/fatih/set"

// Factor represents factor function in MRF. It has two properties: scope and data, scope is an assemble of involving random variables; data is a one-dimensional float64 array storing the value matrix.
type Factor struct {
	scope []int
	data  []float64
}

// increasely sort
func orderScp(sc []int) {
	sort.IntSlice(sc).Sort()
}

// a\b
func scpDiff(a, b []int) []int {
	sa := setFromIntSlice(a)
	sb := setFromIntSlice(b)
	sc := set.Difference(sa, sb)
	c := set.IntSlice(sc)
	orderScp(c)
	return c
}

// a + b
func scpUnion(a, b []int) []int {
	sa := setFromIntSlice(a)
	sb := setFromIntSlice(b)
	sc := set.Union(sa, sb)
	c := set.IntSlice(sc)
	orderScp(c)
	return c
}

// ab
func scpIntsc(a, b []int) []int {
	sa := setFromIntSlice(a)
	sb := setFromIntSlice(b)
	sc := set.Intersection(sa, sb)
	c := set.IntSlice(sc)
	orderScp(c)
	return c
}

// Generates a set form given int slice
func setFromIntSlice(sl []int) *set.SetNonTS {
	s := set.NewNonTS()
	for _, v := range sl {
		s.Add(v)
	}
	return s
}

// the index on the memeroy. i.e. the index of factor data slice
func memIndex(s []int) int {
	idx := 0
	b := 1
	for i := 0; i < len(s); i++ {
		if s[i] > 1 || s[i] < 0 {
			panic(s)
		}
		idx = idx + b*s[i]
		b = b * 2
	}
	return idx
}

// Get returns the actual value. e.g.:
//    value := factor.Get([]int{0,1,0,1})
//    //value is the value of P when x0=0,x1=1,x2=0,x3=1
func (f Factor) Get(s []int) float64 {
	if f.data == nil {
		panic("Facotor.data did not init")
	}
	if len(s) != len(f.scope) {
		panic("Factor.Get scope len dose not match up")
	}

	return f.data[memIndex(s)]
}

// Set sets the value of the Factor given a configuration
func (f Factor) Set(s []int, v float64) {
	if len(s) != len(f.scope) {
		panic("Factor.Set scope len dose not match up")
	}
	idx := memIndex(s)
	f.data[idx] = v
}

// y = 2**x
func pow2(x int) int {
	y := 1
	for i := 0; i < x; i++ {
		y = y * 2
	}
	return y
}

// generate and iter bin-like code slice invoking fn
func walk(l int, fn func([]int)) {
	idx := make([]int, l)
	for i := 0; i < pow2(l); i++ {
		fn(idx)
		// simulate binary addition
		for j := 0; j < l; j++ {
			idx[j] = (idx[j] + 1) % 2
			if idx[j] == 1 {
				break
			}
		}
	}
}

// return the index slice of a's el in the base. i.e. base[rt[i]] = a[i]
// a is neccessarily a subset of base
func addrBook(base, a []int) []int {
	aIdxAddr := []int{}
	has := func(sl []int, v int) bool {
		for _, val := range sl {
			if v == val {
				return true
			}
		}
		return false
	}

	for i, v := range base {
		if has(a, v) {
			aIdxAddr = append(aIdxAddr, i)
		}
	}
	return aIdxAddr
}

// return a*b. merge scopes and multiply values
func FactorProduct(a, b Factor) Factor {
	scopec := scpUnion(a.scope, b.scope)

	// address book based on scopec
	aIdxAddr := addrBook(scopec, a.scope)
	bIdxAddr := addrBook(scopec, b.scope)

	// assign c, new factor
	c := NewFactor(scopec)

	walk(len(scopec), func(idx []int) {
		aIdxSl := []int{}
		bIdxSl := []int{}
		for _, v := range aIdxAddr {
			aIdxSl = append(aIdxSl, idx[v])
		}
		for _, v := range bIdxAddr {
			bIdxSl = append(bIdxSl, idx[v])
		}

		c.Set(idx, a.Get(aIdxSl)*b.Get(bIdxSl))
	})

	return c
}

// return a new factor with all value is 1
func NewFactor(order []int) Factor {
	l := len(order)
	n := pow2(l)
	factor := Factor{order, make([]float64, n)}
	for i, _ := range factor.data {
		factor.data[i] = 1
	}
	return factor
}

// concat a,b given a,b's values and scopes
func catScpIdx(aIdx, bIdx, aScp, bScp []int) []int {
	scp := scpUnion(aScp, bScp)
	scpIdx := make([]int, len(scp))

	foo := func(xIdx, xScp []int) {
		xAddr := addrBook(scp, xScp)
		for i, v := range xAddr {
			scpIdx[v] = xIdx[i]
		}
	}
	foo(aIdx, aScp)
	foo(bIdx, bScp)

	return scpIdx
}

// sum-product
func (fc Factor) sumOut(rm []int) Factor {
	keep := scpDiff(fc.scope, rm)
	nfc := NewFactor(keep)

	walk(len(keep), func(oidx []int) {
		sum := 0.0
		var idx []int
		walk(len(rm), func(iidx []int) {
			idx = catScpIdx(oidx, iidx, keep, rm)
			sum += fc.Get(idx)
		})
		nfc.Set(oidx, sum)
	})

	return nfc
}

// max-product
func (fc Factor) maxOut(rm []int) Factor {
	keep := scpDiff(fc.scope, rm)
	nfc := NewFactor(keep)

	walk(len(keep), func(oidx []int) {
		max := 0.0
		var idx []int
		walk(len(rm), func(iidx []int) {
			idx = catScpIdx(oidx, iidx, keep, rm)
			if max < fc.Get(idx) {
				max = fc.Get(idx)
			}
		})
		nfc.Set(oidx, max)
	})

	return nfc
}
