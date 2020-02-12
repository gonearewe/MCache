package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

type config struct {
	Type, ServerName, Operation        string
	Total, ValueSize, Threads, PipeLen int
}

var Config = &config{}

func init() {
	flag.StringVar(&Config.Type, "type", "tcp", "cache server type")
	flag.StringVar(&Config.ServerName, "h", "localhost", "cache server address")
	flag.IntVar(&Config.Total, "n", 1000, "total number of requests")
	flag.IntVar(&Config.ValueSize, "d", 10, "data size of SET/GET value in bytes")
	flag.IntVar(&Config.Threads, "c", 1, "number of parallel connections")
	flag.StringVar(&Config.Operation, "t", "mixed", "test set, could be get/set/mixed")
	flag.IntVar(&Config.PipeLen, "P", 100, "pipeline length")

	flag.Parse()

	fmt.Println("type is", Config.Type)
	fmt.Println("server is", Config.ServerName)
	fmt.Println("total", Config.Total, "requests")
	fmt.Println("data size is", Config.ValueSize)
	fmt.Println("we have", Config.Threads, "connections")
	fmt.Println("operation is", Config.Operation)
	fmt.Println("pipeline length is", Config.PipeLen)

	rand.Seed(time.Now().UnixNano()) // prepare for random
}
