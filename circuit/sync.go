package circuit

import (
	zkbanw "github.com/akakou/zk-ban/witness"
	tw "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark-crypto/signature"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

const RevocationListSize = 1000

type SyncCircuit struct {
	UserSecretKey frontend.Variable `gnark:"sk"`
	UserPublicKey frontend.Variable `gnark:"pk"`
	Certificate   eddsa.Signature   `gnark:"cert"`
	Period        frontend.Variable `gnark:",public"`

	GroupPublicKey eddsa.PublicKey                       `gnark:",public"`
	Commit2        [RevocationListSize]frontend.Variable `gnark:",public"`
	Commit3        [RevocationListSize]frontend.Variable `gnark:",public"`
}

func (circuit *SyncCircuit) Define(api frontend.API) error {
	err := auth(api, circuit.Period, circuit.UserSecretKey, circuit.UserPublicKey, circuit.Certificate, circuit.GroupPublicKey)
	if err != nil {
		return err
	}

	for i := 0; i < RevocationListSize; i++ {
		commit3, err := mimcHash(api, circuit.Commit2[i], circuit.UserSecretKey)
		if err != nil {
			return err
		}

		api.AssertIsEqual(commit3, circuit.Commit3[i])
	}

	return nil
}

func SyncCircuitWitness(commit2, commit3 [RevocationListSize][]byte, signer *zkbanw.Signer, gpk signature.PublicKey) *SyncCircuit {
	assign := &SyncCircuit{
		UserSecretKey: signer.UserSecretKey.Number,
		UserPublicKey: signer.UserPublicKey.Buffer,
		Period:        signer.Period,
	}

	assign.GroupPublicKey.Assign(tw.BN254, gpk.Bytes())
	assign.Certificate.Assign(tw.BN254, signer.Certificate.Signature)

	for i := 0; i < RevocationListSize; i++ {
		assign.Commit2[i] = commit2[i]
		assign.Commit3[i] = commit3[i]
	}

	return assign
}
