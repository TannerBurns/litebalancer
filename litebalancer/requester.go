package litebalancer

import (
	"math/rand"
	"time"
)

type Requester struct {
	Fn    func() interface{}
	RetFn func(interface{})
	Work  chan Request
}

func (rq *Requester) MakeRequest(work chan<- Request) {
	c := make(chan interface{})
	for {
		time.Sleep(time.Duration(rand.Int63n(int64(time.Millisecond))))
		work <- Request{rq.Fn, c}
		res := <-c
		rq.RetFn(res)
	}
}

func NewRequester(fn func() interface{}, rfn func(interface{})) *Requester {
	r := Requester{fn, rfn, make(chan Request)}
	return &r
}
