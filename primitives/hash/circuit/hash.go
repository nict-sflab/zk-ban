package circuit

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/poseidon2"
)

var HashNew = poseidon2.NewMerkleDamgardHasher

func Hash(api frontend.API, data ...frontend.Variable) (frontend.Variable, error) {
	hasher, err := HashNew(api)
	if err != nil {
		return nil, err
	}

	for _, d := range data {
		hasher.Write(d)
	}

	sum := hasher.Sum()

	return sum, nil
}
