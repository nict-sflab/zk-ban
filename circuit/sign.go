package snark

import (
	"math/big"

	tw "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark-crypto/signature"

	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/signature/eddsa"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

type ProofCircuit struct {
	UserPublicKey  frontend.Variable `gnark:"pk"`
	Certificate    eddsa.Signature   `gnark:"cert"`
	Nonce          frontend.Variable `gnark:"nonce"`
	UserSecretKey  frontend.Variable `gnark:"sk"`
	GroupPublicKey eddsa.PublicKey   `gnark:",public"`
	Message        frontend.Variable `gnark:",public"`
	Hash           frontend.Variable `gnark:",public"`
}

func (circuit *ProofCircuit) Define(api frontend.API) error {
	mimc1, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	mimc2, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	curve, err := twistededwards.NewEdCurve(api, tw.BN254)
	if err != nil {
		return err
	}

	mimc1.Write(circuit.Message, circuit.Nonce, circuit.UserSecretKey)
	h := mimc1.Sum()

	err = eddsa.Verify(curve, circuit.Certificate, circuit.UserPublicKey, circuit.GroupPublicKey, &mimc2)
	if err != nil {
		return err
	}

	api.AssertIsEqual(circuit.Hash, h)

	return nil
}

func ProofWitness(m, r, usk *big.Int, upk, proof, cert []byte, gpk signature.PublicKey) *ProofCircuit {
	assign := &ProofCircuit{
		Nonce:         r,
		UserSecretKey: usk,
		Message:       m,
		Hash:          proof,
		UserPublicKey: upk,
	}

	assign.GroupPublicKey.Assign(tw.BN254, gpk.Bytes())
	assign.Certificate.Assign(tw.BN254, cert)

	return assign
}
