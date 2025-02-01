package circuit

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/liyue201/gnark-circomlib/circuits"
)

func HashMaker(newHasher func(frontend.API) (hash.FieldHasher, error)) func(api frontend.API, data ...frontend.Variable) (frontend.Variable, error) {
	return func(api frontend.API, data ...frontend.Variable) (frontend.Variable, error) {
		hasher, _ := newHasher(api)
		for _, d := range data {
			hasher.Write(d)
		}

		commit := hasher.Sum()

		return commit, nil
	}
}

func NewPoseidon(api frontend.API) (hash.FieldHasher, error) {
	hasher := circuits.NewPoseidonHash(api)
	return hasher, nil
}

func NewMIMC(api frontend.API) (hash.FieldHasher, error) {
	hasher, err := mimc.NewMiMC(api)
	return &hasher, err
}
