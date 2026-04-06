package dump

import (
	"encoding/json"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/precomputes"
	"github.com/akakou/zk-ban/snark"
	"github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/frontend"
)

func Prepare[T frontend.Circuit](c T) ([]byte, []byte, error) {
	cc, err := snark.InitSNARK(c)
	if err != nil {
		return nil, nil, err
	}

	verifierBuf, err := json.Marshal(&cc.VerifyKey)
	if err != nil {
		return nil, nil, err
	}

	prover := snark.SnarkProver{
		ConstraintSystem: cc.ConstraintSystem,
		ProveKey:         cc.ProveKey,
	}

	proverBuf, err := json.Marshal(&prover)
	if err != nil {
		return nil, nil, err
	}

	return proverBuf, verifierBuf, nil
}

func JoinRequestCircuit() ([]byte, []byte, error) {
	return Prepare(&circuit.JoinRequestCircuit{})
}

func SignCircuit() ([]byte, []byte, error) {
	return Prepare(&circuit.SignCircuit{})
}

func MakeUpdateCircuit(rlSize witness.RevocationListSize) ([]byte, []byte, error) {
	rl := witness.EmptyRevocationList(rlSize)

	cc, err := snark.InitSNARK(&precomputes.UpdateCircuit{
		UpdateCircuit: circuit.UpdateCircuit{
			RevocationList: circuit.NewRevocationListAssigned(rl),
		},
	})
	if err != nil {
		return nil, nil, err
	}

	verifier := snark.SizedSnarkVerifier{
		VerifyKey: &cc.VerifyKey,
		RLSize:    rlSize,
	}

	verifierBuf, err := json.Marshal(&verifier)
	if err != nil {
		return nil, nil, err
	}

	prover := snark.SnarkProver{
		ConstraintSystem: cc.ConstraintSystem,
		ProveKey:         cc.ProveKey,
	}

	proverBuf, err := json.Marshal(&prover)
	if err != nil {
		return nil, nil, err
	}

	return proverBuf, verifierBuf, nil
}
