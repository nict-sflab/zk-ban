package witness

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"sort"

	"github.com/akakou/zk-ban/primitives"
	"github.com/akakou/zk-ban/snark"
)

var (
	ErrInvalidRevocationList  = errors.New("invalid revocation list")
	ErrUnsignedRevocationList = errors.New("revocation list has no signed intervals")
	ErrRevokedNym             = errors.New("nym is revoked")
)

type RevocationList []RevokedNymsPerPeriod

type RevokedNymsPerPeriod struct {
	Period *primitives.BigInt
	Nyms   []*primitives.BigInt

	// SignedIntervals authenticates every adjacent pair in the sorted
	// revocation list, including the field-boundary sentinels. The signed
	// message is H(period, H(lower, upper)).
	SignedIntervals []SignedInterval
}

type SignedInterval struct {
	Lower     *primitives.BigInt
	Upper     *primitives.BigInt
	Signature []byte
}

type NymNonMembershipProof struct {
	Lower     *primitives.BigInt
	Upper     *primitives.BigInt
	Signature []byte
}

type NonMembershipProofsPerPeriod struct {
	Period *primitives.BigInt
	Proofs []NymNonMembershipProof
}

type NonMembershipProofList []NonMembershipProofsPerPeriod

func (rl RevocationList) Size() int {
	return rl.Sizes().Size()
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
	return &primitives.BigInt{Int: *n}
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
		exp1 := x - u
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
		rl = append(rl, RevokedNymsPerPeriod{
			Nyms:   nyms,
			Period: primitives.NewBigInt(0),
		})
	}

	return rl
}

// IntervalMessage maps t || H(b || c) to the single field element accepted by
// the gnark EdDSA gadget: H(t, H(b, c)).
func IntervalMessage(period, lower, upper *primitives.BigInt) (*primitives.BigInt, error) {
	if period == nil || lower == nil || upper == nil {
		return nil, fmt.Errorf("%w: nil interval message input", ErrInvalidRevocationList)
	}

	intervalDigest, err := snark.CommitHash(&lower.Int, &upper.Int)
	if err != nil {
		return nil, err
	}

	message, err := snark.CommitHash(&period.Int, intervalDigest)
	if err != nil {
		return nil, err
	}

	return &primitives.BigInt{Int: *message}, nil
}

// SignInterval signs H(t, H(b, c)). The circuit recomputes this message,
// verifies it with the group public key, and proves b < nym < c.
func (gsk *GroupSecretKey) SignInterval(period, lower, upper *primitives.BigInt) (*SignedInterval, error) {
	if gsk == nil || gsk.Signer == nil {
		return nil, fmt.Errorf("%w: missing signing key", ErrInvalidRevocationList)
	}
	if period == nil || period.Sign() < 0 || !period.IsInt64() {
		return nil, fmt.Errorf("%w: period must be a non-negative int64", ErrInvalidRevocationList)
	}
	if lower == nil || upper == nil || lower.Cmp(&upper.Int) >= 0 {
		return nil, fmt.Errorf("%w: interval bounds must satisfy lower < upper", ErrInvalidRevocationList)
	}

	message, err := IntervalMessage(period, lower, upper)
	if err != nil {
		return nil, err
	}

	signature, err := gsk.Sign(message.Bytes(), snark.NewCommitHash())
	if err != nil {
		return nil, err
	}

	return &SignedInterval{
		Lower:     cloneBigInt(lower),
		Upper:     cloneBigInt(upper),
		Signature: append([]byte(nil), signature...),
	}, nil
}

// AuthenticateRevocationList sorts and deduplicates each period's revoked nyms,
// adds 0 and (field modulus - 1) sentinels, then signs every adjacent pair.
// Only the selected signed pair is placed in the circuit witness, so circuit
// cost is independent of the number of revoked nyms in a period.
func (gsk *GroupSecretKey) AuthenticateRevocationList(rl RevocationList) (RevocationList, error) {
	result := make(RevocationList, len(rl))
	fieldMax := new(big.Int).Sub(snark.EcCurve.ScalarField(), big.NewInt(1))

	for periodIndex, revokedPerPeriod := range rl {
		if revokedPerPeriod.Period == nil || revokedPerPeriod.Period.Sign() < 0 || !revokedPerPeriod.Period.IsInt64() {
			return nil, fmt.Errorf("%w: period at index %d must be a non-negative int64", ErrInvalidRevocationList, periodIndex)
		}

		values := make([]*big.Int, 0, len(revokedPerPeriod.Nyms)+2)
		values = append(values, big.NewInt(0), new(big.Int).Set(fieldMax))

		nyms := make([]*primitives.BigInt, 0, len(revokedPerPeriod.Nyms))
		for nymIndex, nym := range revokedPerPeriod.Nyms {
			if nym == nil {
				return nil, fmt.Errorf("%w: nil nym at period %d, index %d", ErrInvalidRevocationList, periodIndex, nymIndex)
			}
			if nym.Sign() < 0 || nym.Cmp(fieldMax) > 0 {
				return nil, fmt.Errorf("%w: nym outside scalar field at period %d, index %d", ErrInvalidRevocationList, periodIndex, nymIndex)
			}

			copyNym := cloneBigInt(nym)
			nyms = append(nyms, copyNym)
			values = append(values, new(big.Int).Set(&nym.Int))
		}

		sort.Slice(values, func(i, j int) bool { return values[i].Cmp(values[j]) < 0 })
		values = deduplicateBigInts(values)

		intervals := make([]SignedInterval, 0, len(values)-1)
		for i := 0; i+1 < len(values); i++ {
			lower := &primitives.BigInt{Int: *new(big.Int).Set(values[i])}
			upper := &primitives.BigInt{Int: *new(big.Int).Set(values[i+1])}

			interval, err := gsk.SignInterval(revokedPerPeriod.Period, lower, upper)
			if err != nil {
				return nil, err
			}
			intervals = append(intervals, *interval)
		}

		result[periodIndex] = RevokedNymsPerPeriod{
			Period:          cloneBigInt(revokedPerPeriod.Period),
			Nyms:            nyms,
			SignedIntervals: intervals,
		}
	}

	return result, nil
}

// NonMembershipProofs selects, for every period and session counter, the signed
// adjacent interval that contains the signer's derived nym.
func (rl RevocationList) NonMembershipProofs(signer *Signer, counterSize int) (NonMembershipProofList, error) {
	if signer == nil || signer.UserSecretKey == nil || signer.UserSecretKey.Number == nil {
		return nil, fmt.Errorf("%w: missing signer secret key", ErrInvalidRevocationList)
	}
	if counterSize < 0 {
		return nil, fmt.Errorf("%w: negative counter size", ErrInvalidRevocationList)
	}

	proofList := make(NonMembershipProofList, 0, len(rl))
	for periodIndex, revokedPerPeriod := range rl {
		period, err := validPeriod(revokedPerPeriod.Period, periodIndex)
		if err != nil {
			return nil, err
		}
		if len(revokedPerPeriod.SignedIntervals) == 0 {
			return nil, fmt.Errorf("%w: period index %d", ErrUnsignedRevocationList, periodIndex)
		}

		proofs := make([]NymNonMembershipProof, 0, counterSize)
		for counter := 0; counter < counterSize; counter++ {
			sessionTag := SessionTag(int64(counter), period)
			nym, err := snark.CommitHash(&sessionTag.Int, &signer.UserSecretKey.Number.Int)
			if err != nil {
				return nil, err
			}

			interval, ok := findContainingInterval(revokedPerPeriod.SignedIntervals, nym)
			if !ok {
				return nil, fmt.Errorf("%w: period=%d counter=%d", ErrRevokedNym, period, counter)
			}

			proofs = append(proofs, copyNonMembershipProof(interval))
		}

		proofList = append(proofList, NonMembershipProofsPerPeriod{
			Period: cloneBigInt(revokedPerPeriod.Period),
			Proofs: proofs,
		})
	}

	return proofList, nil
}

func validPeriod(period *primitives.BigInt, periodIndex int) (int64, error) {
	if period == nil || period.Sign() < 0 || !period.IsInt64() {
		return 0, fmt.Errorf("%w: period at index %d must be a non-negative int64", ErrInvalidRevocationList, periodIndex)
	}
	return period.Int64(), nil
}

func findContainingInterval(intervals []SignedInterval, nym *big.Int) (*SignedInterval, bool) {
	index := sort.Search(len(intervals), func(i int) bool {
		return intervals[i].Upper != nil && intervals[i].Upper.Cmp(nym) > 0
	})
	if index >= len(intervals) {
		return nil, false
	}

	interval := &intervals[index]
	if !validSignedInterval(interval) {
		return nil, false
	}
	if interval.Lower.Cmp(nym) >= 0 || interval.Upper.Cmp(nym) <= 0 {
		return nil, false
	}

	return interval, true
}

func validSignedInterval(interval *SignedInterval) bool {
	return interval != nil &&
		interval.Lower != nil &&
		interval.Upper != nil &&
		interval.Lower.Cmp(&interval.Upper.Int) < 0 &&
		len(interval.Signature) > 0
}

func copyNonMembershipProof(interval *SignedInterval) NymNonMembershipProof {
	return NymNonMembershipProof{
		Lower:     cloneBigInt(interval.Lower),
		Upper:     cloneBigInt(interval.Upper),
		Signature: append([]byte(nil), interval.Signature...),
	}
}

func cloneBigInt(value *primitives.BigInt) *primitives.BigInt {
	if value == nil {
		return nil
	}
	return &primitives.BigInt{Int: *new(big.Int).Set(&value.Int)}
}

func deduplicateBigInts(values []*big.Int) []*big.Int {
	if len(values) < 2 {
		return values
	}

	result := values[:1]
	for _, value := range values[1:] {
		if result[len(result)-1].Cmp(value) != 0 {
			result = append(result, value)
		}
	}
	return result
}
