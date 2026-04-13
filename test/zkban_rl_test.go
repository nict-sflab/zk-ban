package test

import (
	"fmt"
	"math"
	"testing"

	"github.com/akakou/zk-ban/witness"
)

func TestRLSize(t *testing.T) {
	periodNum := 30
	nymNum := 30000
	uRL := EmptyUniformRevocationList(periodNum, nymNum)
	pRL := EmptyProportionalRevocationList(periodNum, nymNum)
	gRL := EmptyGaussianRevocationList(periodNum, nymNum)

	fmt.Printf("u: %v\n", uRL.Sizes())
	fmt.Printf("p: %v\n", pRL.Sizes())
	fmt.Printf("g: %v\n", gRL.Sizes())

	testRLSize(nymNum, uRL, "uniform", 0, t)
	testRLSize(nymNum, pRL, "propotinal", 0, t)
	testRLSize(nymNum, gRL, "gaussian", 0, t)
}

func testRLSize(expected int, rl witness.RevocationList, name string, threshold float64, t *testing.T) {
	t.Run("rl-test-"+name, func(t *testing.T) {
		actual := rl.Size()
		diff := float64(expected - actual)
		if math.Abs(diff) > threshold {
			t.Fatalf("diff is %v (= %v(expected) - %v(actual))", diff, expected, actual)
		}
	})
}
