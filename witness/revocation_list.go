package witness

import (
	"log"
	"math"
	"math/big"

	"github.com/akakou/zk-ban/primitives"
	"github.com/akakou/zk-ban/snark"
)

type RevocationList []RevokedNymsPerSession
type RevokedNymsPerSession struct {
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

func MakeUniformRLSize(sessionNumber, nymsNumberPerSession int) RevocationListSize {
	rlSize := []int{}
	for range sessionNumber {
		rlSize = append(rlSize, nymsNumberPerSession)
	}

	return rlSize
}

func MakeLinearRLSize(sessionNumber, max, min int) RevocationListSize {
	a := float64(max-min) / float64(sessionNumber)

	rlSize := []int{}
	for i := range sessionNumber {
		v := float64(i)*a + float64(min)
		if v < 1 {
			v = 1
		}

		rlSize = append(rlSize, int(v))
	}

	return rlSize
}

func MakeGaussianRLSize(sessionNumber, nymNum int, sigma float64) RevocationListSize {
	u := float64(sessionNumber) - 1

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
	for i := range sessionNumber {
		v := gaussianFunc(float64(i)) * 2 * float64(nymNum)
		if v < 1 {
			v = 1
		}
		rlSize = append(rlSize, int(v))
	}

	return rlSize
}

func MakeUniformRLSizeFromTotal(sessionSize, nymNum int) RevocationListSize {
	if sessionSize > nymNum {
		log.Fatalf("wrong session size and nym num: %v(sess) > %v(nym)\n", sessionSize, nymNum)
	}
	nymNumPerSession := nymNum / sessionSize

	rlSize := MakeUniformRLSize(sessionSize, nymNumPerSession)

	for i := range nymNum - rlSize.Size() {
		rlSize[i] += 1
	}

	return rlSize
}

func MakeProportionalRLSizeFromTotal(sessionSize, nymNum int) RevocationListSize {
	max := nymNum * 2 / sessionSize

	rlSize := MakeLinearRLSize(sessionSize, max, 1)
	rlSize = AjustRLSize(rlSize, nymNum)

	return rlSize
}

func MakeGaussianRLSizeFromTotal(sessionSize, nymNum int, sd float64) RevocationListSize {
	rlSize := MakeGaussianRLSize(sessionSize, nymNum, sd)
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

	for _, nymPerSession := range rlSize {
		nyms := []*primitives.BigInt{}

		for range nymPerSession {
			n := InitBigInt()
			nyms = append(nyms, n)
		}

		// Period is encoded in 64 bits in the circuit; default to 0 for empty entries.
		n := primitives.NewBigInt(0)
		rl = append(rl, RevokedNymsPerSession{
			Nyms:   nyms,
			Period: n,
		})
	}

	return rl
}
