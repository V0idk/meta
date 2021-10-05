package dag

import (
	"container/heap"
	"log"
	. "meta/utils/container"
	"sync"
)

type NODE_NOT_FOUND struct{}

func (e NODE_NOT_FOUND) Error() string {
	return "NODE_NOT_FOUND"
}

type NodeStatus string

const (
	INIT     NodeStatus = "init"
	RUNNING  NodeStatus = "running"
	FINISHED NodeStatus = "finished"
)

type Node struct {
	Id     string
	Func   func(node *Node) //node -> dag -> Executor
	status NodeStatus
	Parent []*Node
	Child  []*Node
	Dag    *DAGExecutor
}

type DAGStatus string

const (
	STOP   DAGStatus = "stop"
	RESUME DAGStatus = "resume"
)

type DAGExecutor struct {
	NodeMap  map[string]*Node
	Context  map[string]map[string]interface{} //node -> key -> value
	Items    map[string]*Item
	Pq       PriorityQueue
	Status   DAGStatus
	finished chan string
	count    int
	wg       *sync.WaitGroup
}

func NewDAGExecutor() *DAGExecutor {
	return &DAGExecutor{
		NodeMap:  make(map[string]*Node),
		Context:  make(map[string]map[string]interface{}),
		finished: make(chan string, 100),
		Pq:       make(PriorityQueue, 0),
		Items:    make(map[string]*Item, 0),
		wg:       &sync.WaitGroup{},
	}
}

func (d *DAGExecutor) AddNode(node *Node) {
	d.NodeMap[node.Id] = node
	node.Dag = d
}

func (d *DAGExecutor) AddEdge(from, to string) error {
	if _, ok := d.NodeMap[from]; !ok {
		return NODE_NOT_FOUND{}
	}
	if _, ok := d.NodeMap[to]; !ok {
		return NODE_NOT_FOUND{}
	}
	d.NodeMap[from].Child = append(d.NodeMap[from].Child, d.NodeMap[to])
	d.NodeMap[to].Parent = append(d.NodeMap[to].Parent, d.NodeMap[from])
	return nil
}

func (d *DAGExecutor) Init() {
	i := 0
	for k, v := range d.NodeMap {
		d.Items[k] = &Item{
			Value:    k,
			Priority: len(v.Parent),
			Index:    i,
		}
		heap.Push(&d.Pq, d.Items[k])
		i += 1
	}
}

func (d *DAGExecutor) Run() {
	defer d.wg.Done()
	for {
		if d.Status == STOP {
			break
		}
		for d.Pq.Len() > 0 && d.Pq[0].Priority == 0 {
			item := heap.Pop(&d.Pq).(*Item)
			go func(item *Item) {
				node := d.NodeMap[item.Value]
				log.Printf("Start %s", node.Id)
				node.Func(node)
				// 所有的状状态由Func控制，实现最大的自由化
				if node.status == FINISHED {
					for _, child := range node.Child {
						childItem := d.Items[child.Id]
						d.Pq.Update(childItem, childItem.Value, childItem.Priority-1)
					}
				}
				// 度数操作必须在发送chan之前。这样确保接收者下次循环前可以启动新的协程。
				d.finished <- item.Value
				log.Printf("Finished %s", node.Id)
			}(item)
		}
		//必须在同一协程中确定d.finished的收发的次数。否则在外部你将无法确定需要收的次数。
		if d.Pq.Len() <= 0 && d.count == len(d.NodeMap) {
			break
		}
		// 最大实时化并发
		<-d.finished
		d.count += 1
	}
}

func (d *DAGExecutor) Resume() {
	d.Status = RESUME
	d.wg.Add(1)
	go d.Run()
}

func (d *DAGExecutor) Wait() {
	d.wg.Wait()
}

//https://stackoverflow.com/questions/23681871/golang-methods-with-same-name-and-arity-but-different-type
func (d *DAGExecutor) WaitNode(ids []string) {
	for _, id := range ids {
		if d.NodeMap[id].status == FINISHED {
			continue
		}
		<-d.finished
	}
}

func (d *DAGExecutor) Stop() {
	d.Status = STOP
}
