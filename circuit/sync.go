package circuit

import (
	"math/big"

	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark-crypto/signature"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

const RevocationListSize = 1000

type SyncCircuit struct {
	UserSecretKey     frontend.Variable                     `gnark:"sk"`
	LastCertificate   eddsa.Signature                       `gnark:"cert"`
	LastPeriod        frontend.Variable                     `gnark:",public"`
	GroupPublicKey    eddsa.PublicKey                       `gnark:",public"`
	NextUserPublicKey frontend.Variable                     `gnark:",public"`
	NextPeriod        frontend.Variable                     `gnark:",public"`
	RevocationList    [RevocationListSize]frontend.Variable `gnark:",public"`
}

func (circuit *SyncCircuit) Define(api frontend.API) error {
	err := certAuth(api, circuit.LastPeriod, circuit.UserSecretKey, circuit.LastCertificate, circuit.GroupPublicKey)
	if err != nil {
		return err
	}

	err = pubKeyAuth(api, circuit.NextPeriod, circuit.UserSecretKey, circuit.NextUserPublicKey)
	if err != nil {
		return err
	}

	for i := 0; i < RevocationListSize; i++ {
		commit, err := hash(api, circuit.RevocationList[i], circuit.UserSecretKey)
		if err != nil {
			return err
		}

		api.AssertIsDifferent(commit, circuit.RevocationList[i])
	}

	return nil
}

func NewSyncCircuitWitness(commit [RevocationListSize]*big.Int, sessionNames [RevocationListSize]*big.Int, next, last *zkbanw.Signer, gpk signature.PublicKey) *SyncCircuit {
	assign := &SyncCircuit{
		UserSecretKey:     last.UserSecretKey.Number,
		LastPeriod:        last.Period,
		NextPeriod:        next.Period,
		NextUserPublicKey: next.UserPublicKey.Number,
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	assign.LastCertificate.Assign(snark.TwistededwardsCurve, last.Certificate.Signature)

	for i := 0; i < RevocationListSize; i++ {
		assign.RevocationList[i] = sessionNames[i]
	}

	return assign
}
