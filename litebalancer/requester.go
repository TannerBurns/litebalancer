package litebalancer

import (
	"errors"
	"math/rand"
	"time"
)

type Requester struct {
	Fn    func([]interface{}) interface{}
	RetFn func(interface{})
	Work  chan Request
}

func (rq *Requester) MakeRequest(work chan<- Request, args ...interface{}) {
	c := make(chan interface{})
	for {
		time.Sleep(time.Duration(rand.Int63n(int64(time.Millisecond))))
		work <- Request{rq.Fn, args, c}
		res := <-c
		if rq.RetFn != nil {
			rq.RetFn(res)
		}
	}
}

func NewRequester(
	fns ...interface{},
) (r *Requester, err error) {
	if len(fns) == 1 {
		r = &Requester{
			Fn:    fns[0].(func([]interface{}) interface{}),
			RetFn: nil,
			Work:  make(chan Request),
		}
	} else if len(fns) == 2 {
		r = &Requester{
			Fn:    fns[0].(func([]interface{}) interface{}),
			RetFn: fns[1].(func(interface{})),
			Work:  make(chan Request),
		}
	} else {
		err = errors.New(
			"failed to create new requester, invalid number of arguments",
		)
	}
	return
}
