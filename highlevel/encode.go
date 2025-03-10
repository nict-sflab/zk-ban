package highlevel

import (
	"github.com/akakou/zk-ban/snark"
)

type HighLevelSnarkProver struct {
	ConstraintSystem []byte
	ProveKey         []byte
}

func (prover HighLevelSnarkProver) ToSnarkProver() (*snark.SnarkProver, error) {
	cs, err := snark.DecodeCircuit(prover.ConstraintSystem)
	if err != nil {
		return nil, err
	}

	proveKeyObj, err := snark.DecodeProverKey(prover.ProveKey)
	if err != nil {
		return nil, err
	}

	proverObj := snark.SnarkProver{
		ConstraintSystem: cs,
		ProveKey:         proveKeyObj,
	}

	return &proverObj, nil
}

func (highl *HighLevelSnarkProver) FromSnarkProver(prover *snark.SnarkProver) error {
	circuitBytes, err := snark.EncodeCircuit(prover.ConstraintSystem)
	if err != nil {
		return err
	}

	proverKeyBytes, err := snark.EncodeProverKey(prover.ProveKey)
	if err != nil {
		return err
	}

	highl.ProveKey = proverKeyBytes
	highl.ConstraintSystem = circuitBytes

	return nil
}
