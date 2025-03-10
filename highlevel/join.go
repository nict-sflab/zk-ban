package highlevel

import (
	"encoding/json"
	"math/big"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/snark"
)

func JoinRequest(period int64, circuitBytes, proveKey []byte) ([]byte, []byte, error) {
	periodBig := big.NewInt(period)

	cs, err := snark.DecodeCircuit(circuitBytes)
	if err != nil {
		return nil, nil, err
	}

	proveKeyObj, err := snark.DecodeProverKey(proveKey)
	if err != nil {
		return nil, nil, err
	}

	prover := snark.SnarkProver{
		ConstraintSystem: cs,
		ProveKey:         proveKeyObj,
	}

	proof, _, err := zkban.JoinRequest(periodBig, &prover)
	if err != nil {
		return nil, nil, err
	}

	proofBytes, err := snark.EncodeProof(proof)
	if err != nil {
		return nil, nil, err
	}

	assignBytes, err := json.Marshal(proofBytes)

	return proofBytes, assignBytes, err
}
