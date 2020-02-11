package main

import (
	"fmt"
	"time"
)

// NOTE: first, call init function in config.go to get configs

func main() {
	var benchmarkResult = &result{countSet: 0, countMiss: 0, countGet: 0, statBuckets: make([]statistic, 0)}
	var resultChan = make(chan *result, Config.Threads)

	// benchmark
	start := time.Now()
	for i := 0; i < Config.Threads; i++ {
		go startBenchmarkRoutine(i, Config.Total/Config.Threads, resultChan)
	}
	for i := 0; i < Config.Threads; i++ {
		benchmarkResult.plus(<-resultChan)
	}
	duration := time.Now().Sub(start) // total benchmark time

	totalCount := benchmarkResult.countGet + benchmarkResult.countMiss + benchmarkResult.countSet
	fmt.Printf("%d records get\n", benchmarkResult.countGet)
	fmt.Printf("%d records miss\n", benchmarkResult.countMiss)
	fmt.Printf("%d records set\n", benchmarkResult.countSet)
	fmt.Printf("%f seconds total\n", duration.Seconds())

	statCountSum := 0
	statTimeSum := time.Duration(0)
	for i, bucket := range benchmarkResult.statBuckets {
		if bucket.count == 0 {
			continue
		}
		statCountSum += bucket.count
		statTimeSum += bucket.duration
		fmt.Printf("%d%% requests < %d ms\n", statCountSum*100/totalCount, i+1)
	}
	fmt.Printf("%d usec average for each request\n", int64(statTimeSum/time.Microsecond)/int64(statCountSum))
	fmt.Printf("throughput is %f MB/bucket\n",
		float64((benchmarkResult.countSet+benchmarkResult.countGet)*Config.ValueSize)/1e6/duration.Seconds())
	fmt.Printf("rps is %f\n", float64(totalCount)/float64(duration.Seconds()))
}
