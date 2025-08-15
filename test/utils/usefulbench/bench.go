package usefulbench

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type UsefulBenchmaker struct {
	Count  int
	Max    int
	Result map[string]map[string]int64
}

func New(max int) *UsefulBenchmaker {
	benchmarker := &UsefulBenchmaker{
		Max:    max,
		Result: make(map[string]map[string]int64),
	}

	benchmarker.Result["env"] = make(map[string]int64)

	return benchmarker

}

func (b *UsefulBenchmaker) Loop() bool {
	res := b.Count < b.Max
	b.Count++
	return res
}

func (b *UsefulBenchmaker) Run(tag string, target func(b Benchmarker)) bool {
	fmt.Printf("Now measure latency of %s...\n", tag)

	// set up tag
	tags := strings.Split(tag, ":")
	if len(tags) != 2 {
		log.Fatalf("Error: length of %v is not 2", tags)
	}

	familyTag := tags[0]
	nameTag := tags[1]

	// run test
	start := time.Now()
	target(b)
	end := time.Since(start)
	latency := end / time.Duration(b.Max)

	b.Count = 0

	// process result
	_, hasFamily := b.Result[familyTag]
	if !hasFamily {
		b.Result[familyTag] = make(map[string]int64)
	}
	b.Result[familyTag][nameTag] = int64(latency)

	fmt.Printf("%s takes %d ns\n", tag, latency)
	fmt.Printf("%s\n\n", b.ResultJson())

	return true
}
