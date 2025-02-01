package circuit

import (
	"github.com/consensys/gnark/frontend"
	"github.com/liyue201/gnark-circomlib/circuits"
)

type hashT func(frontend.API, ...frontend.Variable) (frontend.Variable, error)

var hash = poseidonHash

func poseidonHash(api frontend.API, data ...frontend.Variable) (frontend.Variable, error) {
	poseidon := circuits.NewPoseidonHash(api)

	for _, d := range data {
		poseidon.Write(d)
	}

	commit := poseidon.Sum()

	return commit, nil
}
