package circuit

import (
	nativeposeidon2 "github.com/consensys/gnark-crypto/ecc/bls12-381/fr/poseidon2"
	"github.com/consensys/gnark/frontend"
	stdhash "github.com/consensys/gnark/std/hash"
	"github.com/consensys/gnark/std/hash/mimc"

	stdposeidon2 "github.com/consensys/gnark/std/permutation/poseidon2"
)

func HashMaker(newHasher func(frontend.API) (stdhash.FieldHasher, error)) func(api frontend.API, data ...frontend.Variable) (frontend.Variable, error) {
	return func(api frontend.API, data ...frontend.Variable) (frontend.Variable, error) {
		hasher, err := newHasher(api)
		if err != nil {
			return nil, err
		}

		for _, d := range data {
			hasher.Write(d)
		}

		commit := hasher.Sum()

		return commit, nil
	}
}

func NewPoseidon2(api frontend.API) (stdhash.FieldHasher, error) {
	p := nativeposeidon2.GetDefaultParameters()

	poseidon, err := stdposeidon2.NewPoseidon2FromParameters(api, p.Width, p.NbFullRounds, p.NbPartialRounds)
	if err != nil {
		return nil, err
	}
	return stdhash.NewMerkleDamgardHasher(api, poseidon, 0), nil
}

func NewMIMC(api frontend.API) (stdhash.FieldHasher, error) {
	hasher, err := mimc.NewMiMC(api)
	return &hasher, err
}
