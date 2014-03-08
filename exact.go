package hammingcode

func (dc *Decoder) InitCliqueTree() {
	dc.addSingletons()
	dc.addIndictors()

	newCl := dc.clusters[7:]
	for i:=0; i<7; i++ {
		for index,v := range newCl {
			for _,vv := range v.scope {
				if vv == i {
					newCl[index].factor =
						FactorProduct(newCl[index].factor,dc.clusters[i].factor)
					goto contd
				}
			}
		}
	contd:
	}
	dc.clusters = newCl

	mat := make([][]Factor,3)
	for i:=0; i<3; i++ {
		mat[i] = make([]Factor,3)
	}
	dc.edges = mat

	link := func(i,j int) {
		dc.clusters[i].neighbours = append(dc.clusters[i].neighbours, &dc.clusters[j])
		dc.clusters[j].neighbours = append(dc.clusters[j].neighbours, &dc.clusters[i])
	}

	link(0,1)
	link(1,2)
}

// fucked up solution
func (dc *Decoder) ProcessMess() []int {
	dc.Init()
	out := []int{}
	for i:=0; i<7; i++ {
		if dc.clusters[i].factor.data[0] > dc.clusters[i].factor.data[1] {
			out = append(out,0)
		} else {
			out = append(out,1)
		}
	}
	return out
}
//func (dc *Decoder) {}
// func (dc *Decoder) {}

