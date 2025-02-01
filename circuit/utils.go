package circuit

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

type hashT func(frontend.API, ...frontend.Variable) (frontend.Variable, error)

func hash(api frontend.API, data ...frontend.Variable) (frontend.Variable, error) {
	// hasher := circuits.NewPoseidonHash(api)
	hasher, _ := mimc.NewMiMC(api)

	for _, d := range data {
		hasher.Write(d)
	}

	commit := hasher.Sum()

	return commit, nil
}
