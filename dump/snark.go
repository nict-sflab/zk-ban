package dump

import (
	"encoding/json"

	gnarkserializable "github.com/akakou/gnark-serializable"
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

	verifierBuf, err := json.Marshal(&gnarkserializable.VerifyingKey{cc.VerifyKey})
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
	assignedAccumulator, err := circuit.NewRevocationAccumulatorTemplateAssigned(rl, circuit.MAX)
	if err != nil {
		return nil, nil, err
	}

	cc, err := snark.InitSNARK(&precomputes.UpdateCircuit{
		UpdateCircuit: circuit.UpdateCircuit{
			RevocationAccumulator: assignedAccumulator,
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
