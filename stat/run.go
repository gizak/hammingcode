package main

import "github.com/gizak/hammingcode"
import "math/rand"
import "math"
import "os"
import "fmt"

func main() {
	// test Sum-Product
	test(hammingcode.SumProduct,"sum-p.txt")
	test(hammingcode.MaxProduct,"max-p.txt")
}


func test(mod int, fn string) {
	input := []float64{}
	output := []float64{}
	for i:=0.1; i<50; i = i+ 0.05 {
		errCnt := 0
		for j:=0;j<10;j++{
			dc := hammingcode.NewDecoder(mod)
			in := genCode(i)
			dc.Accept(in)
			res := dc.ProcessMess()
			if chkErr([]int{0,0,0,0,0,0,0},res) {
				errCnt += 1
			}
		}
		input = append(input,i)
		output = append(output,float64(errCnt+1)/(7*10))
	}


	f,err := os.Create(fn)

	if err != nil {
		panic(err)
	}
	defer f.Close()
	for i,_ := range input {
		fmt.Fprintf(f, "%f \t %f\n",math.Log10(input[i]),math.Log10(output[i]))
	}
}

// all zero
func genCode(n float64) []float64 {
	res := make([]float64,7)
	for i:=0; i<7;i++ {
		res[i] = -1+rand.NormFloat64()*n
	}
	return res
}


func chkErr(a,b []int) bool {
	for i,v := range a {
		if v != b[i] {
			return true
		}
	}
	return false
}
