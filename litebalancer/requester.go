package litebalancer

import (
	"fmt"
	"math/rand"
	"time"
)

type Requester struct {
	Fn   func() interface{}
	Work chan Request
}

func (rq *Requester) MakeRequest(work chan<- Request) {
	c := make(chan interface{})
	for {
		time.Sleep(time.Duration(rand.Int63n(int64(time.Millisecond))))
		work <- Request{rq.Fn, c}
		res := <-c
		fmt.Println(res)
	}
}

func NewRequester(fn func() interface{}) *Requester {
	r := Requester{fn, make(chan Request)}
	return &r
}
