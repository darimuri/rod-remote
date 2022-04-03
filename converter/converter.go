package converter

import (
	"fmt"

	"github.com/awalterschulze/gographviz"
	flow "github.com/s8sg/goflow/flow/v1"
)

func Convert(g *gographviz.Graph) (*flow.Dag, error) {
	dag := flow.NewDag()

	for _, n := range g.Nodes.Nodes {
		dag.Node(n.Name, doNothingWorkload)
	}

	for _, e := range g.Edges.Edges {
		dag.Edge(e.Src, e.Dst)
	}

	return dag, nil
}

func doNothingWorkload(data []byte, option map[string][]string) ([]byte, error) {
	return []byte(fmt.Sprintf("you said \"%s\"", string(data))), nil
}
