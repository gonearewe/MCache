package main

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"

	"github.com/gonearewe/MCache/client"
)

func startBenchmarkRoutine(routineID, count int, ch chan<- *result) {
	var client_ = client.New(Config.Type, Config.ServerName)
	var requests = make([]*client.Request, 0, count)
	var res = &result{countMiss: 0, countGet: 0, countSet: 0, statBuckets: make([]statistic, 0)}

	// prepare requests
	for i := 0; i < count; i++ {
		requests = append(requests, &client.Request{Type: "", Key: "", Val: nil, Error: nil})

		// determine benchmark operation type
		var type_ = Config.Operation
		if type_ == "mixed" {
			if rand.Intn(2) == 1 {
				type_ = "get"
			} else {
				type_ = "set"
			}
		}
		requests[i].Type = client.RequestType(type_)

		// determine key and value
		uniqueID := routineID*count + i
		requests[i].Key = strconv.Itoa(uniqueID)
		k := []byte(requests[i].Key)
		// value is a slice of byte starting with its key and supplied to demanded size with meaningless letters
		requests[i].Val = append(k, bytes.Repeat([]byte{'8'}, Config.ValueSize-len(k))...)
	}

	// handle pipelineRun
	if Config.PipeLen > 1 {
		var reqs []*client.Request
		for i, req := range requests {
			if len(reqs) == Config.PipeLen || i == len(requests)-1 {
				pipelineRun(client_, reqs, res)
			} else {
				reqs = append(reqs, req)
			}
		}

		ch <- res
		return
	}

	// not pipeline
	for _, req := range requests {
		run(client_, req, res)
	}
	ch <- res
	return
}

func run(c client.Client, req *client.Request, res *result) {
	expected := req.Val // only used when type is Get

	// benchmark
	start := time.Now()
	c.Run(req)
	duration := time.Now().Sub(start)

	// when the type is Get, verify response
	type_ := req.Type
	if type_ == client.RequestGet {
		if len(req.Val) == 0 {
			type_ = "miss"
		} else if !byteSliceEqual(req.Val, expected) {
			panic("fail to Get expected cache: expecting " + string(expected) + " : got: " + string(req.Val))
		}
	}

	res.addDuration(duration, type_)
}

func pipelineRun(c client.Client, reqs []*client.Request, res *result) {
	expected := make([][]byte, len(reqs)) // only used when type is Get
	for i, req := range reqs {
		if req.Type == client.RequestGet {
			expected[i] = req.Val
		}
	}

	// benchmark
	start := time.Now()
	c.PipelinedRun(reqs)
	duration := time.Now().Sub(start)

	for i, req := range reqs {
		// when the type is Get, verify response
		type_ := req.Type
		if type_ == client.RequestGet {
			if len(req.Val) == 0 {
				type_ = "miss"
			} else if !byteSliceEqual(req.Val, expected[i]) {
				panic("fail to Get expected cache, please make sure it passed tests")
			}
		}

		res.addDuration(duration, type_) // add for each request
	}
}

// byteSliceEqual tells if a and b (both byte slice) has the exact same elements.
func byteSliceEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i, elem := range a {
		if b[i] != elem {
			return false
		}
	}

	return true
}
