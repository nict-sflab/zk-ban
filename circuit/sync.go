package circuit

import (
	"math/big"

	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark-crypto/signature"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

const RevocationListSize = 30000

type SyncCircuit struct {
	UserSecretKey frontend.Variable `gnark:"sk"`
	UserPublicKey frontend.Variable `gnark:"pk"`
	Certificate   eddsa.Signature   `gnark:"cert"`
	Period        frontend.Variable `gnark:",public"`

	GroupPublicKey eddsa.PublicKey                       `gnark:",public"`
	Commit2        frontend.Variable                     `gnark:",public"`
	Commit3        [RevocationListSize]frontend.Variable `gnark:",public"`
}

func (circuit *SyncCircuit) Define(api frontend.API) error {
	err := auth(api, circuit.Period, circuit.UserSecretKey, circuit.UserPublicKey, circuit.Certificate, circuit.GroupPublicKey)
	if err != nil {
		return err
	}

	commit3, err := hash(api, circuit.Commit2, circuit.UserSecretKey)
	if err != nil {
		return err
	}

	for i := 0; i < RevocationListSize; i++ {
		api.AssertIsDifferent(commit3, circuit.Commit3[i])
	}

	return nil
}

func NewSyncCircuitWitness(commit2 *big.Int, commit3 [RevocationListSize]*big.Int, signer *zkbanw.Signer, gpk signature.PublicKey) *SyncCircuit {
	assign := &SyncCircuit{
		UserSecretKey: signer.UserSecretKey.Number,
		UserPublicKey: signer.UserPublicKey.Number,
		Period:        signer.Period,
		Commit2:       commit2,
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	assign.Certificate.Assign(snark.TwistededwardsCurve, signer.Certificate.Signature)

	for i := 0; i < RevocationListSize; i++ {
		assign.Commit3[i] = commit3[i]
	}

	return assign
}
