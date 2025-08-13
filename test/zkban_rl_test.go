package zkbantest

import (
	"fmt"
	"math"
	"testing"
)

func TestRLSize(t *testing.T) {
	rl := EmptyProportionalRevocationList(10, 100)

	s := 0
	for i := range rl {
		s += len(rl[i].Nyms)
		fmt.Printf("%v, ", s)
	}
	fmt.Println("")

	if math.Abs(float64(s)-100.) < 5 {
		panic("")
	}
}
