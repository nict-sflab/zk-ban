package circuit

import (
	"fmt"

	"github.com/akakou/zk-ban/snark"
	"github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/signature/eddsa"
)

type SignedNymInterval struct {
	Lower     frontend.Variable `gnark:",secret"`
	Upper     frontend.Variable `gnark:",secret"`
	Signature eddsa.Signature   `gnark:",secret"`
}

type RevokedNymsPerPeriod struct {
	Period        frontend.Variable `gnark:",public"`
	CounterProofs []SignedNymInterval
}

type RevocationList []RevokedNymsPerPeriod

// NewRevocationListAssigned creates only the circuit shape. Its size depends on
// the number of periods and MaxSession, not on the number of revoked nyms.
func NewRevocationListAssigned(revocationList witness.RevocationList) RevocationList {
	rl := make(RevocationList, 0, len(revocationList))

	for _, revokedPerPeriod := range revocationList {
		period := frontend.Variable(0)
		if revokedPerPeriod.Period != nil {
			period = revokedPerPeriod.Period.Int
		}

		proofs := make([]SignedNymInterval, MaxSession)
		for counter := range proofs {
			proofs[counter] = emptySignedNymInterval()
		}

		rl = append(rl, RevokedNymsPerPeriod{
			Period:        period,
			CounterProofs: proofs,
		})
	}

	return rl
}

func NewRevocationProofAssigned(proofList witness.NonMembershipProofList) (RevocationList, error) {
	rl := make(RevocationList, 0, len(proofList))

	for periodIndex, proofsPerPeriod := range proofList {
		if proofsPerPeriod.Period == nil {
			return nil, fmt.Errorf("period index %d is nil", periodIndex)
		}
		if len(proofsPerPeriod.Proofs) != MaxSession {
			return nil, fmt.Errorf("period index %d has %d proofs; expected %d", periodIndex, len(proofsPerPeriod.Proofs), MaxSession)
		}

		proofs := make([]SignedNymInterval, 0, len(proofsPerPeriod.Proofs))
		for counter, proof := range proofsPerPeriod.Proofs {
			if proof.Lower == nil || proof.Upper == nil || len(proof.Signature) == 0 {
				return nil, fmt.Errorf("invalid non-membership proof at period index %d, counter %d", periodIndex, counter)
			}

			assigned := SignedNymInterval{
				Lower: proof.Lower.Int,
				Upper: proof.Upper.Int,
			}
			assigned.Signature.Assign(snark.TwistededwardsCurve, proof.Signature)
			proofs = append(proofs, assigned)
		}

		rl = append(rl, RevokedNymsPerPeriod{
			Period:        proofsPerPeriod.Period.Int,
			CounterProofs: proofs,
		})
	}

	return rl, nil
}

func (revocationList RevocationList) CheckRevocation(
	usk frontend.Variable,
	publicKey eddsa.PublicKey,
	api frontend.API,
) error {
	curve, err := twistededwards.NewEdCurve(api, snark.TwistededwardsCurve)
	if err != nil {
		return err
	}

	for _, revokedPerPeriod := range revocationList {
		for counter, proof := range revokedPerPeriod.CounterProofs {
			sessionTag := SessionTag(api, revokedPerPeriod.Period, counter)
			nym, err := snark.CircuitHash(api, sessionTag, usk)
			if err != nil {
				return err
			}

			// The EdDSA gadget accepts one field element as its message, so the
			// requested t || H(b || c) is encoded as H(t, H(b, c)).
			intervalDigest, err := snark.CircuitHash(api, proof.Lower, proof.Upper)
			if err != nil {
				return err
			}
			signedMessage, err := snark.CircuitHash(api, revokedPerPeriod.Period, intervalDigest)
			if err != nil {
				return err
			}

			hash, err := snark.NewCircuitHash(api)
			if err != nil {
				return err
			}
			if err := eddsa.Verify(curve, proof.Signature, signedMessage, publicKey, hash); err != nil {
				return err
			}

			// Strict interval membership is the non-membership proof: b < a < c.
			api.AssertIsEqual(api.Cmp(proof.Lower, nym), -1)
			api.AssertIsEqual(api.Cmp(nym, proof.Upper), -1)
		}
	}

	return nil
}

func emptySignedNymInterval() SignedNymInterval {
	return SignedNymInterval{
		Lower: 0,
		Upper: 0,
		Signature: eddsa.Signature{
			R: twistededwards.Point{
				X: 0,
				Y: 0,
			},
			S: 0,
		},
	}
}
