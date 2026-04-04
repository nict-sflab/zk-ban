package witness

import (
	"math/big"
	"sort"

	"github.com/akakou/zk-ban/primitives"
	"github.com/akakou/zk-ban/snark"
	fr_bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	nativeposeidon2 "github.com/consensys/gnark-crypto/ecc/bls12-381/fr/poseidon2"
)

const (
	merkleLeafDomain = int64(1001)
	merkleNodeDomain = int64(1002)
)

type MerkleMembershipProof struct {
	Leaf     *primitives.BigInt
	Index    *primitives.BigInt
	Siblings []*primitives.BigInt
}

type NymNonMembershipProof struct {
	Left  MerkleMembershipProof
	Right MerkleMembershipProof
}

type SessionRevocationAccumulator struct {
	Period        *primitives.BigInt
	Root          *primitives.BigInt
	CounterProofs []NymNonMembershipProof
}

type RevocationAccumulator []SessionRevocationAccumulator

type merkleHashTree struct {
	levels [][]*big.Int
}

func fieldModulus() *big.Int {
	return snark.EcCurve.ScalarField()
}

func normalizeFieldElement(n *big.Int) *big.Int {
	mod := fieldModulus()
	v := new(big.Int).Mod(cloneBigInt(n), mod)
	if v.Sign() < 0 {
		v.Add(v, mod)
	}
	return v
}

func cloneBigInt(n *big.Int) *big.Int {
	if n == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(n)
}

func toPrimitiveBigInt(n *big.Int) *primitives.BigInt {
	return &primitives.BigInt{*cloneBigInt(n)}
}

func canonicalFieldBytes(n *big.Int) []byte {
	value := normalizeFieldElement(n).Bytes()
	if len(value) == fr_bls12381.Bytes {
		return value
	}

	out := make([]byte, fr_bls12381.Bytes)
	copy(out[fr_bls12381.Bytes-len(value):], value)
	return out
}

func circuitCompatibleHash(data ...*big.Int) (*big.Int, error) {
	params := nativeposeidon2.GetDefaultParameters()
	permutation := nativeposeidon2.NewPermutation(params.Width, params.NbFullRounds, params.NbPartialRounds)
	state := canonicalFieldBytes(big.NewInt(0))
	for _, d := range data {
		nextState, err := permutation.Compress(state, canonicalFieldBytes(d))
		if err != nil {
			return nil, err
		}
		state = nextState
	}

	return normalizeFieldElement(new(big.Int).SetBytes(state)), nil
}

func hashMerkleLeaf(value *big.Int) (*big.Int, error) {
	return circuitCompatibleHash(big.NewInt(merkleLeafDomain), value)
}

func hashMerkleNode(left, right *big.Int) (*big.Int, error) {
	return circuitCompatibleHash(big.NewInt(merkleNodeDomain), left, right)
}

func newMerkleHashTree(leaves []*big.Int) (*merkleHashTree, error) {
	if len(leaves) == 0 {
		leaves = []*big.Int{big.NewInt(0)}
	}

	width := 1
	for width < len(leaves) {
		width <<= 1
	}

	level := make([]*big.Int, width)
	lastLeaf := leaves[len(leaves)-1]
	for i := range width {
		value := lastLeaf
		if i < len(leaves) {
			value = leaves[i]
		}

		h, err := hashMerkleLeaf(value)
		if err != nil {
			return nil, err
		}
		level[i] = h
	}

	levels := [][]*big.Int{level}
	for len(level) > 1 {
		next := make([]*big.Int, len(level)/2)
		for i := 0; i < len(level); i += 2 {
			h, err := hashMerkleNode(level[i], level[i+1])
			if err != nil {
				return nil, err
			}
			next[i/2] = h
		}

		levels = append(levels, next)
		level = next
	}

	return &merkleHashTree{
		levels: levels,
	}, nil
}

func (tree *merkleHashTree) root() *big.Int {
	if len(tree.levels) == 0 || len(tree.levels[len(tree.levels)-1]) == 0 {
		return big.NewInt(0)
	}
	return tree.levels[len(tree.levels)-1][0]
}

func (tree *merkleHashTree) siblings(index int) []*big.Int {
	siblings := make([]*big.Int, 0, len(tree.levels)-1)
	current := index
	for depth := 0; depth < len(tree.levels)-1; depth++ {
		siblingIndex := current ^ 1
		siblings = append(siblings, cloneBigInt(tree.levels[depth][siblingIndex]))
		current /= 2
	}
	return siblings
}

func merkleUpperSentinel() *big.Int {
	return new(big.Int).Sub(snark.EcCurve.ScalarField(), big.NewInt(1))
}

func sortedAugmentedSessionLeaves(rps RevokedNymsPerSession) []*big.Int {
	sorted := make([]*big.Int, 0, len(rps.Nyms))
	for _, nym := range rps.Nyms {
		sorted = append(sorted, normalizeFieldElement(&nym.Int))
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Cmp(sorted[j]) < 0
	})

	leaves := make([]*big.Int, 0, len(sorted)+2)
	leaves = append(leaves, big.NewInt(0))
	leaves = append(leaves, sorted...)
	leaves = append(leaves, merkleUpperSentinel())

	return leaves
}

func findNonMembershipPairIndices(leaves []*big.Int, target *big.Int) (int, int) {
	i := sort.Search(len(leaves), func(i int) bool {
		return leaves[i].Cmp(target) >= 0
	})

	if i == len(leaves) {
		return len(leaves) - 2, len(leaves) - 1
	}

	if leaves[i].Cmp(target) == 0 {
		if i+1 < len(leaves) {
			return i, i + 1
		}
		return i - 1, i
	}

	if i == 0 {
		return 0, 1
	}

	return i - 1, i
}

func newMembershipProof(leaves []*big.Int, tree *merkleHashTree, index int) MerkleMembershipProof {
	siblings := tree.siblings(index)

	proofSiblings := make([]*primitives.BigInt, 0, len(siblings))
	for _, sibling := range siblings {
		proofSiblings = append(proofSiblings, toPrimitiveBigInt(sibling))
	}

	return MerkleMembershipProof{
		Leaf:     toPrimitiveBigInt(leaves[index]),
		Index:    primitives.NewBigInt(int64(index)),
		Siblings: proofSiblings,
	}
}

func sessionNym(signer *Signer, period *primitives.BigInt, counter int) (*big.Int, error) {
	// The circuit encodes period/counter in 64 bits. Use the same truncation.
	sessionTag := SessionTag(int64(counter), int64(period.Int.Uint64()))
	return circuitCompatibleHash(&sessionTag.Int, &signer.UserSecretKey.Number.Int)
}

func buildRevocationAccumulator(
	rl RevocationList,
	maxCounter int,
	nymAt func(period *primitives.BigInt, counter int) (*big.Int, error),
) (RevocationAccumulator, error) {
	accumulator := RevocationAccumulator{}

	for _, rps := range rl {
		leaves := sortedAugmentedSessionLeaves(rps)
		tree, err := newMerkleHashTree(leaves)
		if err != nil {
			return nil, err
		}

		counterProofs := make([]NymNonMembershipProof, 0, maxCounter)
		for counter := range maxCounter {
			nym, err := nymAt(rps.Period, counter)
			if err != nil {
				return nil, err
			}

			leftIndex, rightIndex := findNonMembershipPairIndices(leaves, nym)
			counterProofs = append(counterProofs, NymNonMembershipProof{
				Left:  newMembershipProof(leaves, tree, leftIndex),
				Right: newMembershipProof(leaves, tree, rightIndex),
			})
		}

		accumulator = append(accumulator, SessionRevocationAccumulator{
			Period:        &primitives.BigInt{*new(big.Int).Set(&rps.Period.Int)},
			Root:          toPrimitiveBigInt(tree.root()),
			CounterProofs: counterProofs,
		})
	}

	return accumulator, nil
}

func NewRevocationAccumulator(rl RevocationList, signer *Signer, maxCounter int) (RevocationAccumulator, error) {
	return buildRevocationAccumulator(rl, maxCounter, func(period *primitives.BigInt, counter int) (*big.Int, error) {
		return sessionNym(signer, period, counter)
	})
}

func NewRevocationAccumulatorSkeleton(rl RevocationList, maxCounter int) (RevocationAccumulator, error) {
	return buildRevocationAccumulator(rl, maxCounter, func(_ *primitives.BigInt, _ int) (*big.Int, error) {
		return big.NewInt(0), nil
	})
}
