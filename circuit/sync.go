package circuit

import (
	"math/big"

	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark-crypto/signature"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

const RevocationListSize = 100

type SyncCircuit struct {
	UserSecretKey frontend.Variable `gnark:"sk"`
	// Certificate   eddsa.Signature   `gnark:",public"`
	Nonce  frontend.Variable `gnark:",public"`
	Period frontend.Variable `gnark:",public"`

	GroupPublicKey eddsa.PublicKey                       `gnark:",public"`
	RevocationList [RevocationListSize]frontend.Variable `gnark:",public"`
	Commit         frontend.Variable                     `gnark:",public"`
}

func (circuit *SyncCircuit) Define(api frontend.API) error {
	commit, err := hash(api, circuit.Nonce, circuit.UserSecretKey)
	if err != nil {
		return err
	}

	api.AssertIsEqual(circuit.Commit, commit)

	signature, err := hash(api, circuit.Period, circuit.UserSecretKey)
	if err != nil {
		return err
	}

	for i := 0; i < RevocationListSize; i++ {
		api.AssertIsDifferent(signature, circuit.RevocationList[i])
	}

	return nil
}

func NewSyncCircuitWitness(revocationList [RevocationListSize]*big.Int, nonce, commit *big.Int, signer *zkbanw.Signer, gpk signature.PublicKey) *SyncCircuit {
	assign := &SyncCircuit{
		UserSecretKey: signer.UserSecretKey.Number,
		Period:        signer.Period,
		Nonce:         nonce,
		Commit:        commit,
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())

	for i := 0; i < RevocationListSize; i++ {
		assign.RevocationList[i] = revocationList[i]
	}

	return assign
}
