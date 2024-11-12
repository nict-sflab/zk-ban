package zkban

import "github.com/consensys/gnark/frontend"

type AuthCircuit struct {
	SK frontend.Variable `gnark:"sk"`
}

func (circuit *AuthCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(1, 1)
	return nil
}
