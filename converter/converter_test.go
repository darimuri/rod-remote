package converter

import (
	"github.com/awalterschulze/gographviz"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("converter", func() {
	Context("convert to goflow", func() {
		It("", func() {
			graphAst, _ := gographviz.ParseString(`digraph G {}`)
			graph := gographviz.NewGraph()
			if err := gographviz.Analyse(graphAst, graph); err != nil {
				panic(err)
			}
			graph.AddNode("G", "a", nil)
			graph.AddNode("G", "b", nil)
			graph.AddEdge("a", "b", true, nil)

			dag, err := Convert(graph)
			Expect(err).Should(BeNil())
			Expect(dag).ShouldNot(BeNil())
		})
	})
})
