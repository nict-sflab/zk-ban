package snark

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

type CertificateRequestCircuit struct {
	UserSecretKey frontend.Variable `gnark:"sk"`
	UserPublicKey frontend.Variable `gnark:",public"`
	Period        frontend.Variable `gnark:",public"`
}

func (circuit *CertificateRequestCircuit) Define(api frontend.API) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	mimc.Write(circuit.Period, circuit.UserSecretKey)
	h := mimc.Sum()

	api.AssertIsEqual(circuit.UserPublicKey, h)

	return nil
}
