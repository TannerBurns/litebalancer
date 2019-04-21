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
	pool      Pool
	done      chan *Worker
	Completed int
	Max       int
}

func NewBalancer(nums ...int) *Balancer {
	var b *Balancer
	done := make(chan *Worker, nums[1])
	// create nWorker WOK channels
	if len(nums) > 2 {
		b = &Balancer{make(Pool, 0, nums[1]), done, 0, nums[2]}
	} else {
		b = &Balancer{make(Pool, 0, nums[1]), done, 0, -1}
	}
	for i := 0; i < nums[1]; i++ {
		w := &Worker{requests: make(chan Request, nums[0])}
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
			b.Completed++
			if b.Max == -1 {
				continue
			} else if b.Completed == b.Max {
				return
			}
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
