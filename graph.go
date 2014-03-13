package hammingcode

type vertex struct {
	Factor
	edges  *edge
	id     int
	belief Factor
}

type edge struct {
	nodeId int
	value  Factor
	next   *edge
}

type graph struct {
	nodes []vertex
}

func (g graph) getEdge(i, j int) *Factor {
	// iter edges mounted on i to find edge[i->j]
	for edge := g.nodes[i].edges; edge != nil; edge = edge.next {
		if edge.nodeId == j {
			return &edge.value
		}
	}
	return nil
}

func (g graph) getVertex(i int) *vertex {
	return &g.nodes[i]
}

// add an edge if not exist
func (g graph) setEdge(i, j int, nfc Factor) {
	fc := g.getEdge(i, j)
	if fc == nil {
		// if node i doesnt mount any edges
		if g.nodes[i].edges == nil {
			g.nodes[i].edges = &edge{j, nfc, nil}
			return
		}

		// last final val is a pointer pointing to the last edge in node i's edges
		last := g.nodes[i].edges
		for edgeNode := g.nodes[i].edges; edgeNode != nil; edgeNode = edgeNode.next {
			last = edgeNode
		}
		last.next = &edge{j, nfc, nil}
	} else {
		*fc = nfc
	}
}

func (g graph) setVertex(i int, ncl Factor) {
	cl := g.getVertex(i)
	cl.Factor = ncl
}

func newGraph() *graph {
	return &graph{[]vertex{}}
}

func (g *graph) addVertex(fc Factor) *vertex {
	l := len(g.nodes)
	v := vertex{
		Factor: fc,
		id:     l,
		edges:  nil,
	}
	g.nodes = append(g.nodes, v)
	return &v
}
