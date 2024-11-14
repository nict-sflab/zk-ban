package circuit

import (
	poseidon "github.com/AlpinYukseloglu/poseidon-gnark/circuits"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

type hashT func(frontend.API, ...frontend.Variable) (frontend.Variable, error)

var hash = mimcHash

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

func poseidonHash(api frontend.API, data ...frontend.Variable) (frontend.Variable, error) {
	return poseidon.Poseidon(api, data), nil
}
