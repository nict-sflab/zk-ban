package circuit

import (
	"github.com/akakou/zk-ban/snark"
	"github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/frontend"
)

const (
	merkleLeafDomain = int64(1001)
	merkleNodeDomain = int64(1002)
)

type MerkleMembershipProof struct {
	Leaf     frontend.Variable
	Index    frontend.Variable
	Siblings []frontend.Variable
}

type NymNonMembershipProof struct {
	Left  MerkleMembershipProof
	Right MerkleMembershipProof
}

type SessionRevocationAccumulator struct {
	Root          frontend.Variable `gnark:",public"`
	Period        frontend.Variable `gnark:",public"`
	CounterProofs []NymNonMembershipProof
}
type RevocationAccumulator []SessionRevocationAccumulator

func NewRevocationAccumulatorAssigned(revocationAccumulator witness.RevocationAccumulator) RevocationAccumulator {
	assigned := RevocationAccumulator{}

	for _, session := range revocationAccumulator {
		counterProofs := make([]NymNonMembershipProof, 0, len(session.CounterProofs))
		for _, proof := range session.CounterProofs {
			leftSiblings := make([]frontend.Variable, 0, len(proof.Left.Siblings))
			for _, sibling := range proof.Left.Siblings {
				leftSiblings = append(leftSiblings, sibling.Int)
			}

			rightSiblings := make([]frontend.Variable, 0, len(proof.Right.Siblings))
			for _, sibling := range proof.Right.Siblings {
				rightSiblings = append(rightSiblings, sibling.Int)
			}

			counterProofs = append(counterProofs, NymNonMembershipProof{
				Left: MerkleMembershipProof{
					Leaf:     proof.Left.Leaf.Int,
					Index:    proof.Left.Index.Int,
					Siblings: leftSiblings,
				},
				Right: MerkleMembershipProof{
					Leaf:     proof.Right.Leaf.Int,
					Index:    proof.Right.Index.Int,
					Siblings: rightSiblings,
				},
			})
		}

		assigned = append(assigned, SessionRevocationAccumulator{
			Root:          session.Root.Int,
			Period:        session.Period.Int,
			CounterProofs: counterProofs,
		})
	}

	return assigned
}

func NewRevocationAccumulatorTemplateAssigned(revocationList witness.RevocationList, counterSize int) (RevocationAccumulator, error) {
	accumulator, err := witness.NewRevocationAccumulatorSkeleton(revocationList, counterSize)
	if err != nil {
		return nil, err
	}
	return NewRevocationAccumulatorAssigned(accumulator), nil
}

func verifyMerkleMembership(
	api frontend.API,
	root frontend.Variable,
	proof MerkleMembershipProof,
) error {
	digest, err := snark.CircuitHash(api, merkleLeafDomain, proof.Leaf)
	if err != nil {
		return err
	}

	indexBits := api.ToBinary(proof.Index, len(proof.Siblings))
	for depth, sibling := range proof.Siblings {
		left := api.Select(indexBits[depth], sibling, digest)
		right := api.Select(indexBits[depth], digest, sibling)

		digest, err = snark.CircuitHash(api, merkleNodeDomain, left, right)
		if err != nil {
			return err
		}
	}

	api.AssertIsEqual(digest, root)
	return nil
}

func verifyNymNonMembership(
	api frontend.API,
	root frontend.Variable,
	nym frontend.Variable,
	proof NymNonMembershipProof,
) error {
	if err := verifyMerkleMembership(api, root, proof.Left); err != nil {
		return err
	}
	if err := verifyMerkleMembership(api, root, proof.Right); err != nil {
		return err
	}

	api.AssertIsEqual(proof.Right.Index, api.Add(proof.Left.Index, 1))

	api.AssertIsLessOrEqual(proof.Left.Leaf, nym)
	api.AssertIsLessOrEqual(nym, proof.Right.Leaf)
	api.AssertIsDifferent(nym, proof.Left.Leaf)
	api.AssertIsDifferent(nym, proof.Right.Leaf)

	return nil
}
