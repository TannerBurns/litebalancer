# LiteBalancer

    Simple load balancer using goroutines and channels
    
    Requester -> Balancer -> Request -> Requester

# How to use

    Create a start function that receives a slice of interface and returns an interface.
	Create a end function if needed the receives a type of interface that will come 
	from the start function.
	If there is a need for dynamic arguments to the start function, they can be provided
	in the creation of the new request.
    

# Example

    This example is how to run the load balancer for N amount of time, maxWork being the N value.
	No value has to be given to the balancer and -1 will be set as default.

	Start function: randoNumbers
	End function: printNumbers
	Requesters: 10
	Workers: 10
	N: 1000
    

``` go
package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/TannerBurns/litebalancer/litebalancer"
)

// function that returns type interface to generate data
// args will be a slice of the arguments provided
// args[0] will need to be typecasted into the appropriate datatype
func randoNumbers(args []interface{}) interface{} {
	nums := args[0].([]int)
	min := nums[0]
	max := nums[1]
	return rand.Intn((max - min) + min)
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
	const maxWork = 100

	// initialize arguments for start function
	args := make([]interface{}, 2)
	args[0] = 1
	args[1] = 50

	rq, err := litebalancer.NewRequester(randoNumbers, printNumbers)
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