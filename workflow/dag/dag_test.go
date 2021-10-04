package dag

import (
	"log"
	"testing"
	"time"
)

func a(node *Node) {
	log.Printf("A")
	time.Sleep(1 * time.Second)
	node.status = FINISHED
}

func b(node *Node) {
	log.Printf("B")
	time.Sleep(1 * time.Second)
	node.status = FINISHED
}

func c(node *Node) {
	log.Printf("C")
	time.Sleep(1 * time.Second)
	node.status = FINISHED
}

func d(node *Node) {
	log.Printf("D")
	time.Sleep(1 * time.Second)
	node.status = FINISHED
}

func TestDAGExecutor_Resume(t *testing.T) {
	dag := NewDAGExecutor()
	dag.AddNode(&Node{
		Id:   "b",
		Func: b,
	})
	dag.AddNode(&Node{
		Id:   "c",
		Func: c,
	})

	dag.AddNode(&Node{
		Id:   "a",
		Func: a,
	})
	dag.AddNode(&Node{
		Id:   "d",
		Func: d,
	})
	dag.AddEdge("a", "b")
	dag.AddEdge("a", "c")
	dag.AddEdge("b", "d")
	dag.AddEdge("c", "d")
	dag.Init()
	dag.Resume()
	dag.Wait()
}
