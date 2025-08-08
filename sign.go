package zkban

import (
	"math/big"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/primitives"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
)

func max() *big.Int {
	var i, e = big.NewInt(2), big.NewInt(32)
	i.Exp(i, e, nil)

	return i
}

type Signature struct {
	Commit *zkbanw.SignCommit
	Proof  groth16.Proof
}

func Sign(m *primitives.BigInt, counter int64, signer *zkbanw.Signer, gpk *zkbanw.GroupPublicKey, prover *snark.SnarkProver) (*Signature, error) {
	sessionTag := signer.SessionTag(counter)

	comm, err := signer.CommitSign(m, sessionTag)
	if err != nil {
		return nil, err
	}

	wit, err := circuit.NewSignWitness(m, sessionTag, comm, signer, gpk)
	if err != nil {
		return nil, err
	}

	proof, err := groth16.Prove(prover.ConstraintSystem, prover.ProveKey, wit)
	if err != nil {
		return nil, err
	}

	signature := Signature{
		Commit: comm,
		Proof:  proof,
	}
	return &signature, nil
}

func (signature *Signature) Verify(m *primitives.BigInt, counter, period int64, gpk *zkbanw.GroupPublicKey, verifyKey groth16.VerifyingKey) error {
	sessionTag := zkbanw.SessionTag(counter, period)

	pubWit, err := circuit.NewPublicSignWitness(m, sessionTag, signature.Commit, period, gpk)
	if err != nil {
		return err
	}

	err = groth16.Verify(signature.Proof, verifyKey, pubWit)
	return err
}
