# LiteBalancer

    Simple load balancer using goroutines and channels
    
    Requester -> Balancer -> Request -> Requester

# How to use

    Implement a function that returns a type of interface to generate data
    Implement a function that receives a type of interface to handle data
    Create a new requester and pass your newly created functions
    Initialize the 100 requesters
    Run a new balancer to handle the work and return the value

# Example

    package main

    import (
        "fmt"
        "math/rand"

        "github.com/TannerBurns/litebalancer/litebalancer"
    )

    // function that returns type interface to generate data
    func randoNumbers() interface{} {
        min := 1
        max := 50
        return rand.Intn((max - min) + min)
    }

    // function that receives type interface to display the data
    func printNumbers(num interface{}) {
        fmt.Println(num.(int))
    }

    func main() {
        // initialize number of requesters and workers
        const numRequesters = 100
        const numWorkers = 10

        // create new requester and send the two functions
        rq := litebalancer.NewRequester(randoNumbers, printNumbers)
        // initialize work for requesters
        for i := 0; i < numRequesters; i++ {
            go rq.MakeRequest(rq.Work)
        }
        // run a new balancer to handle work
        litebalancer.NewBalancer(numRequesters, numWorkers).Balance(rq.Work)
    }

