package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/awalterschulze/gographviz"
)

func main() {
	// Create an Abstract Syntax Tree (AST) from a DOT representation.
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	ast, err := gographviz.Parse(b)
	if err != nil {
		log.Fatal(err)
	}

	// Render the AST into a graph.
	g := gographviz.NewGraph()
	gographviz.Analyse(ast, g)

	// For each node, find the cycles it participates in. Keep all nodes
	// that participate in cylces in r.
	r := map[string]bool{}
	for node := range g.Nodes.Lookup {
		union(r, cycles(node, g.Edges.SrcToDsts))
	}

	// To render a graph consisting only of cycles:
	//
	// 1. Create a new graph.
	// 2. Loop over all edges in the input graph.
	// 3. Add an edge to the new graph if both the source and the
	//    destination are in a cycle.
	fg := gographviz.NewGraph()
	fg.SetName(g.Name)
	fg.SetDir(true)
	for _, edge := range g.Edges.Edges {
		if r[edge.Src] && r[edge.Dst] {
			fg.AddNode("", edge.Src, nil)
			fg.AddNode("", edge.Dst, nil)
			fg.AddEdge(edge.Src, edge.Dst, true, attrsToMap(edge.Attrs))
		}
	}

	io.WriteString(os.Stdout, fg.String())
}

func attrsToMap(attrs gographviz.Attrs) map[string]string {
	ret := make(map[string]string)
	for k, v := range attrs {
		ret[string(k)] = v

	}
	return ret
}

// Cycles returns a set of all nodes that are involved in cycles starting in
// start. It is horribly inefficient, both from an implementation as well as
// an algorithmic point of view. It uses DFS to find cycles. I didn't feel
// like implementing Johnson's algorithm. If you want that, Python has a
// library called networkx, which has a simple_cycles() function that uses
// an efficient algorithm.
func cycles(start string, adj map[string]map[string][]*gographviz.Edge) map[string]bool {
	return dfs(start, start, nil, adj)
}

// dfs returns all nodes that are in a simple cycle starting in start.
func dfs(start, cur string, stack []string, adj map[string]map[string][]*gographviz.Edge) map[string]bool {
	// Check for cycle or cross.
	if in(stack, cur) {
		if cur == start {
			return set(stack) // Found a cycle, return the current stack as a set.
		}
		return nil // Found a cross, just return.
	}

	r := map[string]bool{}
	childStack := append(stack, cur)
	// Loop over all possible destinations of cur.
	for dst := range adj[cur] {
		// Add all nodes that are in a cycle as found by the recursive call
		// to dfs.
		union(r, dfs(start, dst, childStack, adj))
	}
	return r
}

// Check if elem is in stack.
func in(stack []string, elem string) bool {
	for _, s := range stack {
		if s == elem {
			return true
		}
	}
	return false
}

// Adds b's elements to a, in-place
func union(a, b map[string]bool) {
	for k, v := range b {
		if v {
			a[k] = v
		}
	}
}

// Turn a list into a set.
func set(a []string) map[string]bool {
	r := map[string]bool{}
	for _, k := range a {
		r[k] = true
	}
	return r
}
