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
	SessionTag *primitives.BigInt
	Nyms       []*primitives.BigInt
}

func (rl RevocationList) Size() int {
	sum := 0
	for _, r := range rl {
		sum += len(r.Nyms)
	}

	return sum
}

func (rl RevocationList) Sizes() []int {
	result := []int{}
	for _, r := range rl {
		result = append(result, len(r.Nyms))
	}

	return result
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

func MakeUniformRLSize(sessionNumber, nymsNumberPerSession int) []int {
	rlSize := []int{}
	for range sessionNumber {
		rlSize = append(rlSize, nymsNumberPerSession)
	}

	return rlSize
}

func MakeLinearRLSize(sessionNumber, max, min int) []int {
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

func MakeGaussianRLSize(sessionNumber, nymNum int, sigma float64) []int {
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

func AjustRLSize(rlSize []int, nymNum int) []int {
	sum := 0
	for _, r := range rlSize {
		sum += r
	}

	diff := sum - nymNum

	for {
		for i, r := range rlSize {
			if diff == 0 {
				return rlSize
			}

			if r == 1 {
			} else if diff > 0 {
				rlSize[i]--
				diff--
			} else {
				rlSize[i]++
				diff++
			}
		}
	}
}

func EmptyRevocationList(nymSizePerSessions []int) RevocationList {
	rl := RevocationList{}

	for _, nymPerSession := range nymSizePerSessions {
		nyms := []*primitives.BigInt{}

		for range nymPerSession {
			n := InitBigInt()
			nyms = append(nyms, n)
		}

		n := InitBigInt()
		rl = append(rl, RevokedNymsPerSession{
			Nyms:       nyms,
			SessionTag: n,
		})
	}

	return rl
}
