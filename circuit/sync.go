package circuit

import (
	"math/big"

	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark-crypto/signature"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

const RevocationListSize = 10000

type SyncCircuit struct {
	UserSecretKey frontend.Variable `gnark:"sk"`
	Certificate   eddsa.Signature   `gnark:"cert"`
	Period        frontend.Variable `gnark:",public"`
	Basename      frontend.Variable `gnark:",public"`

	GroupPublicKey eddsa.PublicKey                       `gnark:",public"`
	Commit         [RevocationListSize]frontend.Variable `gnark:",public"`
}

func (circuit *SyncCircuit) Define(api frontend.API) error {
	err := auth(api, circuit.Period, circuit.UserSecretKey, circuit.Certificate, circuit.GroupPublicKey)
	if err != nil {
		return err
	}

	commit3, err := hash(api, circuit.Period, circuit.Basename, circuit.UserSecretKey)
	if err != nil {
		return err
	}

	for i := 0; i < RevocationListSize; i++ {
		api.AssertIsDifferent(commit3, circuit.Commit[i])
	}

	return nil
}

func NewSyncCircuitWitness(commit [RevocationListSize]*big.Int, bsn *big.Int, signer *zkbanw.Signer, gpk signature.PublicKey) *SyncCircuit {
	assign := &SyncCircuit{
		UserSecretKey: signer.UserSecretKey.Number,
		Period:        signer.Period,
		Basename:      bsn,
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	assign.Certificate.Assign(snark.TwistededwardsCurve, signer.Certificate.Signature)

	for i := 0; i < RevocationListSize; i++ {
		assign.Commit[i] = commit[i]
	}

	return assign
}
