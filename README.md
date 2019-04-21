# LiteBalancer

    Simple load balancer using goroutines and channels
    
    Requester -> Balancer -> Request -> Requester

# How to use

    Create a start function that receives a slice of interface and returns an interface
    

# Example

    This example is how to run the load balancer for N amount of time, maxWork being the N value. 
    Notice the call of NewBalancer has also changed and has an extra variable. 
    In this example we only want to print 1000 random numbers.

``` go
package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/TannerBurns/litebalancer/litebalancer"
)

// function that returns type interface to generate data
func randoNumbers(args []interface{}) interface{} {
	return rand.Intn((args[1].(int) - args[0].(int)) + args[0].(int))
}

// function that receives type interface to display the data
func printNumbers(num interface{}) {
	fmt.Println(num.(int))
}

func main() {
	// initialize number of requesters, workers and maxwork allowed
	const numRequesters = 10
	const numWorkers = 10
	// if maxWork is not supplied a value of -1 will be applied
	// if value of maxWork is -1, the balancer will run forever
	const maxWork = 1000

	// initialize arguments for start function
	args := make([]interface{}, 2)
	args[0] = 1
	args[1] = 50

    rq, err := litebalancer.NewRequester(randoNumbers, printNumbers)
    // check to make sure requester is ready
	if err != nil {
		log.Fatal(err)
	}
	// initialize work for requesters, arguments can be sent to start function
	// when make request is made to link the arguments and start function
	for i := 0; i < numRequesters; i++ {
		go rq.MakeRequest(rq.Work, args)
	}
	// run a new balancer to handle work
	litebalancer.NewBalancer(
		numRequesters,
		numWorkers,
		maxWork,
	).Balance(
		rq.Work,
	)
}
```