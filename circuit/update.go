package circuit

import (
	"github.com/akakou/zk-ban/primitives"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

var UpdatePreparableIndex = 4

type UpdateCircuit struct {
	UserSecretKey  frontend.Variable `gnark:",secret"`
	Credential     eddsa.Signature   `gnark:",secret"`
	NextInfo       PublicKeyAuthInfo
	CurrentInfo    PublicKeyAuthInfo
	GroupPublicKey eddsa.PublicKey `gnark:",public"`
	RevocationList RevocationList
}

func (circuit *UpdateCircuit) Define(api frontend.API) error {
	err := authCredential(api, circuit.UserSecretKey, circuit.Credential, circuit.CurrentInfo.Period, circuit.GroupPublicKey)
	if err != nil {
		return err
	}

	err = authPubKey(api, circuit.CurrentInfo, circuit.UserSecretKey, ONE_TIME_TICKET)
	if err != nil {
		return err
	}

	err = authPubKey(api, circuit.NextInfo, circuit.UserSecretKey, PUBLIC_KEY)
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

func NewUpdateCircuitWitness(
	nextPeriod int64,
	nextPublicKey *zkbanw.UserPublicKey,
	ticket *zkbanw.OneTimeTicket,
	signer *zkbanw.Signer,
	revocationList zkbanw.RevocationList,
	gpk *zkbanw.GroupPublicKey,
) (witness.Witness, error) {
	assign := &UpdateCircuit{
		UserSecretKey: signer.UserSecretKey.Number.Int,
		CurrentInfo: PublicKeyAuthInfo{
			Period:        signer.Period,
			UserPublicKey: ticket.Number.Int,
		},
		NextInfo: PublicKeyAuthInfo{
			Period:        nextPeriod,
			UserPublicKey: nextPublicKey.Number.Int,
		},
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	assign.Credential.Assign(snark.TwistededwardsCurve, signer.Credential.Signature)

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

	return frontend.NewWitness(assign, snark.EcCurve.ScalarField())
}

func NewPublicUpdateCircuitWitness(
	nextPeriod int64,
	nextPublicKey *zkbanw.UserPublicKey,
	updateTicket *zkbanw.OneTimeTicket,
	lastPeriod int64,
	revocationList zkbanw.RevocationList,
	gpk *zkbanw.GroupPublicKey,
) (witness.Witness, error) {
	wit, err := NewUpdateCircuitWitness(nextPeriod, nextPublicKey, updateTicket, &zkbanw.Signer{
		UserSecretKey: &zkbanw.UserSecretKey{
			primitives.NewBigInt(0),
		},
		Credential: &zkbanw.Credential{
			Signature: make([]byte, 32),
		},
		Period: lastPeriod,
	}, revocationList, gpk)

	if err != nil {
		return nil, err
	}

	return wit.Public()
}
