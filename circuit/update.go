package circuit

import (
	"math/big"

	"github.com/akakou/zk-ban/commit"
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark-crypto/signature"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

// const RevocationListSize = 223
// const RevocationSpeed = 0.82

const RevokedNymsPerSession = 130
const SessionSize = 270

type UpdateCircuit struct {
	UserSecretKey     frontend.Variable                                     `gnark:"sk"`
	LastCredential    eddsa.Signature                                       `gnark:"cert"`
	LastPeriod        frontend.Variable                                     `gnark:",public"`
	GroupPublicKey    eddsa.PublicKey                                       `gnark:",public"`
	NextUserPublicKey frontend.Variable                                     `gnark:",public"`
	NextPeriod        frontend.Variable                                     `gnark:",public"`
	SessionTags       [SessionSize]frontend.Variable                        `gnark:",public"`
	RevokedNyms       [SessionSize][RevokedNymsPerSession]frontend.Variable `gnark:",public"`
}

func (circuit *UpdateCircuit) Define(api frontend.API) error {
	err := authCert(api, circuit.LastPeriod, circuit.UserSecretKey, circuit.LastCredential, circuit.GroupPublicKey)
	if err != nil {
		return err
	}

	err = authPubKey(api, circuit.NextPeriod, circuit.UserSecretKey, circuit.NextUserPublicKey)
	if err != nil {
		return err
	}

	for i := 0; i < SessionSize; i++ {
		nym, err := snark.CircuitHash(api, circuit.SessionTags[i], circuit.UserSecretKey)
		if err != nil {
			return err
		}

		// max := int(RevocationSpeed * float64(i))
		// for j := 0; j < max; j++ {
		for j := 0; j < RevokedNymsPerSession; j++ {
			api.AssertIsDifferent(nym, circuit.RevokedNyms[i][j])
		}
	}

	return nil
}

func NewUpdateCircuitWitness(commit [SessionSize][RevokedNymsPerSession]*big.Int, sessionNames [SessionSize]*big.Int, next, last *commit.Signer, gpk signature.PublicKey) *UpdateCircuit {
	assign := &UpdateCircuit{
		UserSecretKey:     last.UserSecretKey.Number,
		LastPeriod:        last.Period,
		NextPeriod:        next.Period,
		NextUserPublicKey: next.UserPublicKey.Number,
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	assign.LastCredential.Assign(snark.TwistededwardsCurve, last.Credential.Signature)

	for i := 0; i < SessionSize; i++ {
		assign.SessionTags[i] = sessionNames[i]

		for j := 0; j < RevokedNymsPerSession; j++ {
			assign.RevokedNyms[i][j] = commit[i][j]
		}
	}

	return assign
}
