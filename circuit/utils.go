package snark

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

func mimcHash(api frontend.API, data ...frontend.Variable) (frontend.Variable, error) {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return nil, err
	}

	for _, d := range data {
		mimc.Write(d)
	}

	commit := mimc.Sum()

	return commit, nil
}
