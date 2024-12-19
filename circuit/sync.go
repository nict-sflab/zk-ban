package circuit

import (
	"math/big"

	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark-crypto/signature"
	"github.com/consensys/gnark/frontend"
)

// const RevocationListSize = 223
// const RevocationSpeed = 0.82
const RevocationListSize = 111
const PeriodSize = 270

type SyncCircuit struct {
	UserSecretKey frontend.Variable `gnark:"sk"`
	UserPublicKey frontend.Variable `gnark:",public"`
	// Certificate   eddsa.Signature               `gnark:"cert"`
	Period [PeriodSize]frontend.Variable `gnark:",public"`
	// GroupPublicKey eddsa.PublicKey                                   `gnark:",public"`
	Commit [PeriodSize][RevocationListSize]frontend.Variable `gnark:",public"`
}

func (circuit *SyncCircuit) Define(api frontend.API) error {
	// err := auth(api, circuit.Period, circuit.UserSecretKey, circuit.UserPublicKey, circuit.Certificate, circuit.GroupPublicKey)
	// if err != nil {
	// 	return err
	// }

	for i := 0; i < PeriodSize; i++ {
		commit3, err := hash(api, circuit.Period[i], circuit.UserSecretKey)
		if err != nil {
			return err
		}

		// max := int(RevocationSpeed * float64(i))
		// for j := 0; j < max; j++ {
		for j := 0; j < RevocationListSize; j++ {
			api.AssertIsDifferent(commit3, circuit.Commit[i][j])
		}
	}

	return nil
}

func NewSyncCircuitWitness(commit [PeriodSize][RevocationListSize]*big.Int, bsn *big.Int, periods [PeriodSize]*big.Int, signer *zkbanw.Signer, gpk signature.PublicKey) *SyncCircuit {
	assign := &SyncCircuit{
		UserSecretKey: signer.UserSecretKey.Number,
		UserPublicKey: signer.UserPublicKey.Number,
	}

	// assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	// assign.Certificate.Assign(snark.TwistededwardsCurve, signer.Certificate.Signature)

	for i := 0; i < PeriodSize; i++ {
		assign.Period[i] = periods[i]

		for j := 0; j < RevocationListSize; j++ {
			assign.Commit[i][j] = commit[i][j]
		}
	}

	return assign
}
