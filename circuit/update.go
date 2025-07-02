package circuit

import (
	"github.com/akakou/zk-ban/snark"
	"github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

var UpdatePreparableIndex = 5

type UpdateCircuit struct {
	UserSecretKey  frontend.Variable `gnark:",secret"`
	CurrentInfo    CredentialAuthInfo
	NextInfo       PublicKeyAuthInfo
	GroupPublicKey eddsa.PublicKey `gnark:",public"`
	RevocationList RevocationList
}

func (circuit *UpdateCircuit) Define(api frontend.API) error {
	err := authCredential(api, circuit.CurrentInfo, circuit.UserSecretKey, circuit.GroupPublicKey)
	if err != nil {
		return err
	}

	err = authPubKey(api, circuit.NextInfo, circuit.UserSecretKey)
	if err != nil {
		return err
	}

	for _, revokedPerSession := range circuit.RevocationList {
		nym, err := snark.CircuitHash(api, revokedPerSession.SessionTag, circuit.UserSecretKey)
		if err != nil {
			return err
		}

		for _, revokedNym := range revokedPerSession.Nyms {
			api.AssertIsDifferent(nym, revokedNym)
		}
	}

	return nil
}

func (circuit *UpdateCircuit) PreparableIndex() int {
	return UpdatePreparableIndex
}

func NewUpdateCircuitWitness(next, last *witness.Signer, revocationList witness.RevocationList, gpk *witness.GroupPublicKey) *UpdateCircuit {
	assign := &UpdateCircuit{
		UserSecretKey: last.UserSecretKey.Number,
		CurrentInfo: CredentialAuthInfo{
			Period: last.Period,
		},
		NextInfo: PublicKeyAuthInfo{
			Period:        next.Period,
			UserPublicKey: next.UserPublicKey.Number,
		},
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	assign.CurrentInfo.Credential.Assign(snark.TwistededwardsCurve, last.Credential.Signature)

	rl := RevocationList{}

	for _, rps := range revocationList {
		nyms := []frontend.Variable{}
		for _, nym := range rps.Nyms {
			nyms = append(nyms, nym)
		}

		rl = append(rl, RevokedNymsPerSession{
			Nyms:       nyms,
			SessionTag: rps.SessionTag,
		})
	}

	assign.RevocationList = rl

	return assign
}
