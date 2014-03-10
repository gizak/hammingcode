package hammingcode

import "sort"
import "github.com/fatih/set"

type Factor struct {
	scope []int
	data  []float64
}

func orderScp(sc []int){
	sort.IntSlice(sc).Sort()
}

func scpDiff(a,b []int) []int{
	sa := setFromIntSlice(a)
	sb := setFromIntSlice(b)
	sc := set.Difference(sa,sb)
	c := set.IntSlice(sc)
	orderScp(c)
	return  c
}

func scpUnion(a, b []int) []int{
	sa := setFromIntSlice(a)
	sb := setFromIntSlice(b)
	sc := set.Union(sa,sb)
	c := set.IntSlice(sc)
	orderScp(c)
	return  c
}

func scpIntsc(a,b []int) []int{
	sa := setFromIntSlice(a)
	sb := setFromIntSlice(b)
	sc := set.Intersection(sa,sb)
	c := set.IntSlice(sc)
	orderScp(c)
	return  c
}


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


func (f Factor) Get(s []int) float64 {
	if len(s) != len(f.scope) {
		panic("Factor.Get scope len dose not match up")
	}
	return f.data[memIndex(s)]
}


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
	aIdxAddr := addrBook(scopec,a.scope)
	bIdxAddr := addrBook(scopec,b.scope)

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
