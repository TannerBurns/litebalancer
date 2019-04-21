package litebalancer

type Request struct {
	Fn       func([]interface{}) interface{}
	Args     []interface{}
	Response chan interface{}
}

type Worker struct {
	requests chan Request
	pending  int
	index    int
}

func (w *Worker) Work(done chan *Worker) {
	for {
		req := <-w.requests
		req.Response <- req.Fn(req.Args)
		done <- w
	}
}
