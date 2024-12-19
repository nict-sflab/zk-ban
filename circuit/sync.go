package circuit

import (
	"math/big"

	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark-crypto/signature"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

// const RevocationListSize = 223
// const RevocationSpeed = 0.82
const RevocationListSize = 111
const SessionSize = 270

type SyncCircuit struct {
	UserSecretKey  frontend.Variable                                  `gnark:"sk"`
	UserPublicKey  frontend.Variable                                  `gnark:",public"`
	Certificate    eddsa.Signature                                    `gnark:"cert"`
	Period         frontend.Variable                                  `gnark:",public"`
	GroupPublicKey eddsa.PublicKey                                    `gnark:",public"`
	SessionNames   [SessionSize]frontend.Variable                     `gnark:",public"`
	Commits        [SessionSize][RevocationListSize]frontend.Variable `gnark:",public"`
}

func (circuit *SyncCircuit) Define(api frontend.API) error {
	err := auth(api, circuit.Period, circuit.UserSecretKey, circuit.UserPublicKey, circuit.Certificate, circuit.GroupPublicKey)
	if err != nil {
		return err
	}

	for i := 0; i < SessionSize; i++ {
		commit, err := hash(api, circuit.SessionNames[i], circuit.UserSecretKey)
		if err != nil {
			return err
		}

		// max := int(RevocationSpeed * float64(i))
		// for j := 0; j < max; j++ {
		for j := 0; j < RevocationListSize; j++ {
			api.AssertIsDifferent(commit, circuit.Commits[i][j])
		}
	}

	return nil
}

func NewSyncCircuitWitness(commit [SessionSize][RevocationListSize]*big.Int, sessionNames [SessionSize]*big.Int, signer *zkbanw.Signer, gpk signature.PublicKey) *SyncCircuit {
	assign := &SyncCircuit{
		UserSecretKey: signer.UserSecretKey.Number,
		UserPublicKey: signer.UserPublicKey.Number,
		Period:        signer.Period,
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	assign.Certificate.Assign(snark.TwistededwardsCurve, signer.Certificate.Signature)

	for i := 0; i < SessionSize; i++ {
		assign.SessionNames[i] = sessionNames[i]

		for j := 0; j < RevocationListSize; j++ {
			assign.Commits[i][j] = commit[i][j]
		}
	}

	return assign
}
