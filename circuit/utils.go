package circuit

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/vocdoni/gnark-crypto-primitives/poseidon"
)

type hashT func(frontend.API, ...frontend.Variable) (frontend.Variable, error)

var hash = poseidonHash

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
	// fmt.Print("circuit:", data, "\n")
	commit := poseidon.Hash(api, data...)
	return commit, nil
}
