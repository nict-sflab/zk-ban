package zkbantest

import (
	"fmt"
	"math"
	"testing"

	"github.com/akakou/zk-ban/witness"
)

const THREASHOLD = -1

func TestRLSize(t *testing.T) {
	uRL := EmptyUniformRevocationList(10, 100)
	pRL := EmptyProportionalRevocationList(10, 100)
	gRL := EmptyGaussianRevocationList(180, 90000)

	testRLSize(100, uRL, "uniform", t)
	testRLSize(100, pRL, "propotinal", t)
	testRLSize(100, gRL, "gaussian", t)
}

func testRLSize(size int, rl witness.RevocationList, name string, t *testing.T) {
	t.Run("rl-test-"+name, func(t *testing.T) {
		sum := 0
		for i := range rl {
			l := len(rl[i].Nyms)
			sum += l
			fmt.Printf("%d, ", l)
		}
		fmt.Println("")

		diff := float64(size - sum)
		if math.Abs(diff) > THREASHOLD {
			t.Fatalf("diff is %v (= %v(size) - %v(sum))", diff, size, sum)
		}
	})
}
