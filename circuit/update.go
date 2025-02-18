package circuit

import (
	"github.com/akakou/zk-ban/commit"
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

type UpdateCircuit struct {
	UserSecretKey  frontend.Variable `gnark:",secret"`
	CurrentInfo    CredentialAuthInfo
	NextInfo       PublicKeyAuthInfo
	RevocationList RevocationList
	GroupPublicKey eddsa.PublicKey `gnark:",public"`
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
		nym, err := snark.CircuitHash(api, revokedPerSession.Period, revokedPerSession.Basename, circuit.UserSecretKey)
		if err != nil {
			return err
		}

		for _, revokedNym := range revokedPerSession.Nyms {
			api.AssertIsDifferent(nym, revokedNym)
		}
	}

	return nil
}

func NewUpdateCircuitWitness(next, last *commit.Signer, revocationList commit.RevocationList) *UpdateCircuit {
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

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, last.GroupPublicKey.Bytes())
	assign.CurrentInfo.Credential.Assign(snark.TwistededwardsCurve, last.Credential.Signature)

	rl := RevocationList{}

	for _, rps := range revocationList {
		nyms := []frontend.Variable{}
		for _, nym := range rps.Nyms {
			nyms = append(nyms, nym)
		}

		rl = append(rl, RevokedNymsPerSession{
			Nyms:     nyms,
			Period:   rps.Period,
			Basename: rps.Basename,
		})
	}

	assign.RevocationList = rl

	return assign
}
