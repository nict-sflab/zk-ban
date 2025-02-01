package circuit

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/liyue201/gnark-circomlib/circuits"
)

type hashT func(frontend.API, ...frontend.Variable) (frontend.Variable, error)

type Hasher interface {
	Write(data ...frontend.Variable)
	Sum() frontend.Variable
}

var hash = hashMaker(newMIMC)

func hashMaker(newHasher func(frontend.API) (Hasher, error)) func(api frontend.API, data ...frontend.Variable) (frontend.Variable, error) {
	return func(api frontend.API, data ...frontend.Variable) (frontend.Variable, error) {
		hasher, _ := newHasher(api)
		for _, d := range data {
			hasher.Write(d)
		}

		commit := hasher.Sum()

		return commit, nil
	}
}

func newPoseidon(api frontend.API) (Hasher, error) {
	hasher := circuits.NewPoseidonHash(api)
	return hasher, nil
}

func newMIMC(api frontend.API) (Hasher, error) {
	hasher, err := mimc.NewMiMC(api)
	return &hasher, err
}
