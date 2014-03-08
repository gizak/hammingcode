package hammingcode

import "sort"
import "github.com/fatih/set"
import _"fmt"


type Factor struct {
	order []int
	data []float64
}


func setFromIntSlice(sl []int) *set.SetNonTS {
	s := set.NewNonTS()
	for _,v := range sl {
		s.Add(v)
	}
	return s
}

func setFromFactor(a Factor) *set.SetNonTS {
	s := set.NewNonTS()
	for i := 0; i < len(a.order); i++ {
		s.Add(a.order[i])
	}
	return s
}


func memIndex(s []int) int {
	idx := 0
	b := 1
	for i := 0; i < len(s); i++ {
		if s[i] > 1 || s[i] < 0 {
			panic(s)
		}
		idx = idx + b * s[i]
		b = b * 2
	}
	return idx
}


func (f Factor) Get(s []int) float64 {
	return f.data[memIndex(s)]
}



func (f Factor) Set(s []int, v float64) {
	idx := memIndex(s)
	f.data[idx] = v
}

func pow2(x int) int {
	y := 1
	for i:=0; i<x; i++{
		y = y*2
	}
	return y
}


func walk(l int, fn func([]int)) {
	idx := make([]int,l)
	for i := 0; i < pow2(l); i++ {
		fn(idx)
		// simulate binary addition
		for j := 0; j<l; j++ {
			idx[j] = (idx[j] + 1)%2
			if idx[j] == 1 {
				break
			}
		}

	}
}


func getAddrBook(base, a []int) []int{
	aIdxAddr := []int{}
	has := func(sl []int, v int) bool{
		for _,val := range sl {
			if v == val {
				return true
			}
		}
		return false
	}
	
	for i,v := range base {
		if has(a,v) {
			aIdxAddr = append(aIdxAddr,i)
		}
	}
	return aIdxAddr
}


func FactorProduct(a,b Factor) Factor {
	// union set slice, ordered
	sa := setFromFactor(a)
	sb := setFromFactor(b)
	sc := set.Union(sa,sb)
	ordc := set.IntSlice(sc)
	sort.IntSlice(ordc).Sort()

	// address book based on ordc
	aIdxAddr := []int{}
	bIdxAddr := []int{}
	has := func(f Factor, v int) bool{
		for _,val := range f.order {
			if v == val {
				return true
			}
		}
		return false
	}
	for i,v := range ordc {
		if has(a,v) {
			aIdxAddr = append(aIdxAddr,i)
		}
		if has(b,v) {
			bIdxAddr = append(bIdxAddr,i)
		}
	}
	
	// assign c
	c := Factor{ordc,make([]float64,pow2(len(ordc)))}

	walk(len(ordc),func(idx []int){
		aIdxSl := []int{}
		bIdxSl := []int{}
		for _,v := range aIdxAddr {
			aIdxSl = append(aIdxSl,idx[v])
		}
		for _,v := range bIdxAddr {
			bIdxSl = append(bIdxSl,idx[v])
		}
		//fmt.Printf("%+v\n",b.Get(bIdxSl))
		//fmt.Printf("%+v\n",bIdxSl)
		c.Set(idx, a.Get(aIdxSl) * b.Get(bIdxSl))
	})
	
	return c
}

func NewFactor(order []int) Factor{
	l := len(order)
	if l>0 {
		n := pow2(l)
		factor := Factor{order,make([]float64,n)}
		for i,_ := range factor.data {
			factor.data[i] = 1
		}
		return factor
	}
	return Factor{}
}
