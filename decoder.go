package hammingcode

import "math"
import "github.com/fatih/set"
import "sort"
//import "fmt"

type Decoder struct {
	chkMatrix [][]int
	clusters []cluster
	edges [][]Factor
	noiseSigma float64
	rcvdCode []float64
	iterCount int
	selectedProd int
}


type cluster struct {
	factor  Factor
	scope []int
	belief Factor
	neighbours []*cluster
	index int
}

const MaxProduct = 1
const SumProduct = 0

func NewDecoder(choice int) *Decoder {
	return &Decoder{selectedProd:choice}
}



func (dc *Decoder) SetCheckMatrix(m [][]int) {
	dc.chkMatrix = m
}

func (dc *Decoder) SetNoiseLevel(sigma float64) {
	dc.noiseSigma = sigma
}

func (dc *Decoder) cmpPostProb(z float64) (x0, x1 float64) {
	foo := func(x int) float64 {
		return 1/(1+math.Exp((1-2*float64(x))*4*z/math.Pow(dc.noiseSigma,2))) 
	}
	return foo(0),foo(1)
}


func (dc *Decoder) addSingletons() {
	for i,v := range dc.rcvdCode {
		factor := Factor{[]int{i},make([]float64,2)}
		factor.data[0],factor.data[1] = dc.cmpPostProb(v)
		clst := cluster{
			factor:factor,
			scope: factor.order,
			//belief: 0,
			neighbours: []*cluster{},
			index: i,
		}
		dc.clusters = append(dc.clusters,clst)
	}
}

func (dc *Decoder) addIndictors() {
	for i:=0; i<len(dc.chkMatrix); i++ {
		scope := []int{}
		factor := Factor{[]int{},make([]float64,pow2(7))}
		clst := cluster{
			factor:factor,
			scope: factor.order,
			//belief: 0,
			neighbours: []*cluster{},
			index: 7+i,
		}

		// compute idx scope
		for j:=0; j<7; j++ {
			if dc.chkMatrix[i][j] == 1 {
				scope = append(scope,j)
				}
		}
		
		// compute factor.data
		walk(7,func(idx []int){
			xor := 0
			for _,v := range scope {
				xor = xor + idx[v]
			}
			xor = xor % 2
			factor.Set(idx,float64(xor))
			})
		
		clst.factor.order = scope
		clst.scope = scope
		// add them in clusters
		dc.clusters = append(dc.clusters,clst)
	}
}


func (dc *Decoder) linkClusters() {
	for i:=7; i<7+len(dc.chkMatrix); i++ {
		for _,v := range dc.clusters[i].scope {
			dc.clusters[i].neighbours =
				append(dc.clusters[i].neighbours,&dc.clusters[v])
			dc.clusters[v].neighbours =
				append(dc.clusters[v].neighbours,&dc.clusters[i])
		}
	}
}


// need to improve
func (dc *Decoder) initEdges() {
	n := 7+len(dc.chkMatrix)
	dc.edges = make([][]Factor,n)
	for i:=0; i<n; i++ {
		msgVec := make([]Factor,n)
		for j,_ := range msgVec {
			if i == j {
				continue
			}
			incSet := set.Intersection(
				setFromIntSlice(dc.clusters[i].scope),
				setFromIntSlice(dc.clusters[j].scope),
			)
			incSl := set.IntSlice(incSet)
			sort.IntSlice(incSl).Sort()
			msgVec[j] = NewFactor(incSl)
		}
		dc.edges[i] = msgVec
	}
}


// Init
func (dc *Decoder) Init() error{
	// generate factors
	// generate singleton factors
	dc.addSingletons()
	
	// generate parity factors
	dc.addIndictors()
	
	// generate graph
	dc.linkClusters()
	
	// init all edges messages with 1
	dc.initEdges()

	return	nil
}



func (dc *Decoder) Accept(c []float64) {
	dc.rcvdCode = c
}


func (dc *Decoder) Next() (bool, error){
	dc.iterCount += 1

	// update messages
	n := 7+len(dc.chkMatrix)
	cache := make([][]Factor, n);
	for i:=0; i<n; i++ {
		cache[i] = make([]Factor,n)
	}
	
	for _,v := range dc.clusters {
		i := v.index
		for _,vv := range v.neighbours {
			j := vv.index
			cache[i][j] = dc.msgTo(i,j)
		}
	}
	dc.edges = cache

	// update belief
	for i,_ := range dc.clusters {
		dc.clusters[i].belief = dc.belief(i)
	}
	return true,nil
}


func (dc *Decoder) msgTo(i,j int) Factor {
	setI := setFromIntSlice(dc.clusters[i].scope)
	setJ := setFromIntSlice(dc.clusters[j].scope)
	diffSet := set.Difference(setI,setJ)
	
	phi := dc.clusters[i].factor
	for _,v := range dc.clusters[i].neighbours {
		if v.index != j {
			phi = FactorProduct(phi,dc.edges[v.index][i])
			sum := 0.0
			for _,v := range phi.data {
				sum += v
			}
			for i,_ := range phi.data {
				phi.data[i] = phi.data[i] / sum
			}
		}
	}

	return elOut(phi,set.IntSlice(diffSet),dc.selectedProd)
}


func (dc *Decoder) belief(i int) Factor {
	phi := dc.clusters[i].factor
	for _,v := range dc.clusters[i].neighbours {
		phi = FactorProduct(phi,dc.edges[v.index][i])
		sum := 0.0
		for _,v := range phi.data {
			sum += v
		}
		for i,_ := range phi.data {
			phi.data[i] = phi.data[i] / sum
		}
	}
	return phi
}


func elOut(fc Factor, rm []int, selectFn int) Factor {
	elOutFn := sumOut
	if selectFn == MaxProduct {
		elOutFn = maxOut
	}
	
	base := fc.order

	keepSet := set.Difference(
		setFromIntSlice(base),
		setFromIntSlice(rm),
	)
	keep := set.IntSlice(keepSet)
	
	sort.IntSlice(keep).Sort()
	sort.IntSlice(rm).Sort()

	keepAddr := getAddrBook(base,keep)
	rmAddr := getAddrBook(base,rm)

	tmpl := make([]int,len(base))
	nfc := Factor{keep,make([]float64,pow2(len(keep)))}

	walk(len(keep),func(outter []int){
		tail := 0.0
		walk(len(rm),func(inner []int){
			for i,v := range keepAddr {
				tmpl[v] = outter[i]
			}
			for i,v := range rmAddr {
				tmpl[v] = inner[i]
			}
			value := fc.Get(tmpl)
			tail = elOutFn(value,tail)
		})
		//fmt.Printf("%+v,%+v\n",tmpl,tail)
		nfc.Set(outter,tail)
	})

	return nfc
}


func sumOut(v, prev float64) float64 {
	return prev + v
}



func maxOut(v, prev float64) float64 {
	if v > prev {
		return v
	}
	return prev
}


func (dc *Decoder) isIterable() bool{
	if dc.iterCount > 5 {
		return false
	}
	return true
}
