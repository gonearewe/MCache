package main

import (
	"time"

	"github.com/gonearewe/MCache/benchmark/client"
)

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
			panic("fail to Get expected cache, please make sure it passed tests")
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
	c.PipelineRun(reqs)
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
