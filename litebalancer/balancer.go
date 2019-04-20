package litebalancer

import "container/heap"

type Pool []*Worker

func (p Pool) Len() int           { return len(p) }
func (p Pool) Less(i, j int) bool { return p[i].pending < p[j].pending }
func (p Pool) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p *Pool) Push(x interface{}) {
	*p = append(*p, x.(*Worker))
}
func (p *Pool) Pop() interface{} {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[0 : n-1]
	return x
}

type Balancer struct {
	pool Pool
	done chan *Worker
}

func NewBalancer(nRequester int, nWorker int) *Balancer {
	done := make(chan *Worker, nWorker)
	// create nWorker WOK channels
	b := &Balancer{make(Pool, 0, nWorker), done}
	for i := 0; i < nWorker; i++ {
		w := &Worker{requests: make(chan Request, nRequester)}
		// put them in heap
		heap.Push(&b.pool, w)
		go w.Work(b.done)
	}
	return b
}

func (b *Balancer) Balance(work chan Request) {
	for {
		select {
		case req := <-work:
			b.dispatch(req)
		case w := <-b.done:
			b.completed(w)
		}
	}
}

func (b *Balancer) dispatch(req Request) {
	//find least loaded worker
	w := heap.Pop(&b.pool).(*Worker)
	//send that worker the task
	w.requests <- req
	//increment que of work for that worker
	w.pending++
	//add to heap
	heap.Push(&b.pool, w)
}

func (b *Balancer) completed(w *Worker) {
	//remove from queue
	w.pending--
	//remove from heap
	heap.Remove(&b.pool, w.index)
	//put back into available workers
	heap.Push(&b.pool, w)
}
