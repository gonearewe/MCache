package main

import (
	"time"

	"github.com/gonearewe/MCache/client"
)

// statistic is the element of one bucket, it contains sum of counts and durations
// within a certain benchmark duration range.
type statistic struct {
	count    int
	duration time.Duration
}

type result struct {
	countGet  int
	countSet  int
	countMiss int
	// statBuckets is used for classifying a number of benchmark results(duration),
	// index of statBuckets plus one indicates the maximum duration(measured by millisecond) this bucket can hold,
	// the larger the index is, the more time the benchmark counts inside it cost.
	// EXAMPLE: duration 3.56 ms (in range [3,4)) belongs to statBuckets[3]
	statBuckets []statistic
}

func (r *result) plus(src *result) {
	r.countGet += src.countGet
	r.countSet += src.countSet
	r.countMiss += src.countMiss

	for i, stat := range src.statBuckets {
		r.updateStatBuckets(i, stat)
	}
	// gap:=len(src.statBuckets)- len(r.statBuckets)
	//	// if gap>0{
	//	// 	src.statBuckets = append(src.statBuckets, make([]statistic,gap)...)
	//	// }
	//	// for i:=range r.statBuckets{
	//	// 	if i>= len(src.statBuckets){
	//	// 		break
	//	// 	}
	//	//
	//	// 	r.statBuckets[i].count+=src.statBuckets[i].count
	//	// 	r.statBuckets[i].duration+=src.statBuckets[i].duration
	//	// }
}

// addDuration increases both count of specific type and count in statBuckets by one,
// and apart from one count, it also updates statBuckets with duration.
func (r *result) addDuration(d time.Duration, type_ client.RequestType) {
	switch type_ {
	case client.RequestSet:
		r.countSet++
	case client.RequestGet:
		r.countGet++
	default:
		r.countMiss++
	}

	bucketID := int(d / time.Millisecond)                           // classify this duration to determine which bucket it belongs to
	r.updateStatBuckets(bucketID, statistic{count: 1, duration: d}) // also increase count in buckets
}

// updateStatBuckets updates result's statBuckets, filling new statistic into given
// index of bucket, if the bucket of given index doesn't exist, it enlarges the size of buckets.
func (r *result) updateStatBuckets(bucketID int, stat statistic) {
	if bucketID > len(r.statBuckets)-1 { // enlarge the size of buckets
		extraBuckets := make([]statistic, bucketID-len(r.statBuckets)+1)
		r.statBuckets = append(r.statBuckets, extraBuckets...)
	}

	bucket := r.statBuckets[bucketID] // bucket is just a copy, not the element itself
	bucket.count += stat.count
	bucket.duration += stat.duration
	r.statBuckets[bucketID] = bucket // NOTICE: this statement is required
}
