package dag

import (
	"log"
	"testing"
	"time"
)

func a(node *Node) {
	log.Printf("a")
	time.Sleep(1 * time.Second)
	node.status = FINISHED
}

func b(node *Node) {
	log.Printf("b")
	time.Sleep(1 * time.Second)
	node.status = FINISHED
}

func c(node *Node) {
	log.Printf("c")
	time.Sleep(1 * time.Second)
	node.status = FINISHED
}

func d(node *Node) {
	log.Printf("d")
	time.Sleep(1 * time.Second)
	node.status = FINISHED
}

func printThing(a string) func(node *Node) {
	return func(node *Node) {
		log.Printf(a)
		time.Sleep(1 * time.Second)
		node.status = FINISHED
	}
}

func TestDAGExecutor_Resume(t *testing.T) {
	dag := NewDAGExecutor()
	var foo = "abcdefghijklmnopqrstuvwxyz"
	for _, c := range foo {
		dag.AddNode(&Node{
			Id:   string(c),
			Func: printThing(string(c)),
		})
	}
	dag.AddEdge("a", "b")
	dag.AddEdge("a", "c")
	dag.AddEdge("b", "d")
	dag.AddEdge("c", "d")
	dag.AddEdge("y", "z")
	dag.AddEdge("g", "a")
	dag.AddEdge("d", "x")
	dag.AddEdge("x", "m")
	dag.AddEdge("x", "n")
	dag.Init()
	dag.Resume()
	dag.Wait()
}
