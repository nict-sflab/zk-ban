package zkban

import (
	"math/big"

	zkbanc "github.com/akakou/zk-ban/circuit"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
)

func CertificateRequest(period *big.Int, usk *zkbanw.UserSecretKey, pk groth16.ProvingKey, ccs constraint.ConstraintSystem) (groth16.Proof, witness.Witness, error) {
	upk, err := usk.PublicKey(period)
	if err != nil {
		return nil, nil, err
	}

	witness := zkbanc.CertificateRequestWitness(usk.UserSecretKey, upk.UserPublicKey, period)

	snark, pubWit, _, err := proveSNARK(witness, pk, ccs)
	if err != nil {
		return nil, nil, err
	}

	return snark, pubWit, nil
}
