package zkban

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

type AuthCircuit struct {
	// PK frontend.Variable `gnark:"pk"`
	// Cert  frontend.Variable `gnark:"cert"`
	Nonce     frontend.Variable `gnark:"nonce"`
	SecretKey frontend.Variable `gnark:"sk"`
	Message   frontend.Variable `gnark:",public"`
	Hash      frontend.Variable `gnark:",public"`
}

func (circuit *AuthCircuit) Define(api frontend.API) error {
	mimc, _ := mimc.NewMiMC(api)

	mimc.Write(circuit.Message, circuit.Nonce, circuit.SecretKey)
	h := mimc.Sum()

	api.AssertIsEqual(circuit.Hash, h)

	return nil
}
