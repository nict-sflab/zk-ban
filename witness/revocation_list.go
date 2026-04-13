package witness

import (
	"log"
	"math"
	"math/big"

	"github.com/akakou/zk-ban/primitives"
	"github.com/akakou/zk-ban/snark"
)

type RevocationList []RevokedNymsPerPeriod
type RevokedNymsPerPeriod struct {
	Period *primitives.BigInt
	Nyms   []*primitives.BigInt
}

func (rl RevocationList) Size() int {
	sum := rl.Sizes().Size()

	return sum
}

func (rl RevocationList) Sizes() RevocationListSize {
	result := RevocationListSize{}
	for _, r := range rl {
		result = append(result, len(r.Nyms))
	}

	return result
}

type RevocationListSize []int

func (rl RevocationListSize) Size() int {
	sum := 0
	for _, r := range rl {
		sum += r
	}

	return sum
}

var InitBigInt = ZeroInitBigInt
var ZeroInitBigInt = func() *primitives.BigInt { return primitives.NewBigInt(0) }
var MimcInitBigInt = func() *primitives.BigInt {
	n, err := snark.CommitHash(big.NewInt(202506252))
	if err != nil {
		log.Fatal(err)
	}
	return &primitives.BigInt{*n}
}

func MakeUniformRLSize(periodNumber, nymsNumberPerPeriod int) RevocationListSize {
	rlSize := []int{}
	for range periodNumber {
		rlSize = append(rlSize, nymsNumberPerPeriod)
	}

	return rlSize
}

func MakeLinearRLSize(periodNumber, max, min int) RevocationListSize {
	a := float64(max-min) / float64(periodNumber)

	rlSize := []int{}
	for i := range periodNumber {
		v := float64(i)*a + float64(min)
		if v < 1 {
			v = 1
		}

		rlSize = append(rlSize, int(v))
	}

	return rlSize
}

func MakeGaussianRLSize(periodNumber, nymNum int, sigma float64) RevocationListSize {
	u := float64(periodNumber) - 1

	var gaussianFunc = func(x float64) float64 {
		exp1 := (x - u)
		exp1 *= exp1
		exp2 := 2 * math.Pow(sigma, 2)
		exp := -1 * exp1 / exp2

		area := math.Sqrt(2*math.Pi) * sigma
		result := math.Exp(exp) / area

		return result
	}

	rlSize := []int{}
	for i := range periodNumber {
		v := gaussianFunc(float64(i)) * 2 * float64(nymNum)
		if v < 1 {
			v = 1
		}
		rlSize = append(rlSize, int(v))
	}

	return rlSize
}

func MakeUniformRLSizeFromTotal(periodSize, nymNum int) RevocationListSize {
	if periodSize > nymNum {
		log.Fatalf("wrong period size and nym num: %v(period) > %v(nym)\n", periodSize, nymNum)
	}
	nymNumPerPeriod := nymNum / periodSize

	rlSize := MakeUniformRLSize(periodSize, nymNumPerPeriod)

	for i := range nymNum - rlSize.Size() {
		rlSize[i] += 1
	}

	return rlSize
}

func MakeProportionalRLSizeFromTotal(periodSize, nymNum int) RevocationListSize {
	max := nymNum * 2 / periodSize

	rlSize := MakeLinearRLSize(periodSize, max, 1)
	rlSize = AjustRLSize(rlSize, nymNum)

	return rlSize
}

func MakeGaussianRLSizeFromTotal(periodSize, nymNum int, sd float64) RevocationListSize {
	rlSize := MakeGaussianRLSize(periodSize, nymNum, sd)
	rlSize = AjustRLSize(rlSize, nymNum)

	return rlSize
}

func AjustRLSize(rlSize RevocationListSize, nymNum int) RevocationListSize {
	sum := 0
	for _, r := range rlSize {
		sum += r
	}

	diff := sum - nymNum

	for {
		oneCount := 0
		for i, r := range rlSize {
			if diff == 0 {
				return rlSize
			}

			if r == 1 {
				oneCount++
			} else if diff > 0 {
				rlSize[i]--
				diff--
			} else {
				rlSize[i]++
				diff++
			}
		}
		if oneCount == len(rlSize) {
			panic("at least one element in rlSize is non-one")
		}
	}
}

func EmptyRevocationList(rlSize RevocationListSize) RevocationList {
	rl := RevocationList{}

	for _, nymPerPeriod := range rlSize {
		nyms := []*primitives.BigInt{}

		for range nymPerPeriod {
			n := InitBigInt()
			nyms = append(nyms, n)
		}

		// Period is encoded in 64 bits in the circuit; default to 0 for empty entries.
		n := primitives.NewBigInt(0)
		rl = append(rl, RevokedNymsPerPeriod{
			Nyms:   nyms,
			Period: n,
		})
	}

	return rl
}
