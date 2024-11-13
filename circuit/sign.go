package snark

import (
	"math/big"

	tw "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark-crypto/signature"

	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/signature/eddsa"

	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

type SignCircuit struct {
	UserPublicKey  frontend.Variable `gnark:"pk"`
	Certificate    eddsa.Signature   `gnark:"cert"`
	Nonce          frontend.Variable `gnark:"nonce"`
	UserSecretKey  frontend.Variable `gnark:"sk"`
	GroupPublicKey eddsa.PublicKey   `gnark:",public"`
	Message        frontend.Variable `gnark:",public"`
	Commit         frontend.Variable `gnark:",public"`
}

func (circuit *SignCircuit) Define(api frontend.API) error {
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

	api.AssertIsEqual(circuit.Commit, h)

	return nil
}

func NewSignWitness(m, r *big.Int, proof []byte, signer *zkbanw.Signer, gpk signature.PublicKey) *SignCircuit {
	assign := &SignCircuit{
		Nonce:         r,
		UserSecretKey: signer.UserSecretKey.Number,
		Message:       m,
		Commit:        proof,
		UserPublicKey: signer.UserPublicKey.Buffer,
	}

	assign.GroupPublicKey.Assign(tw.BN254, gpk.Bytes())
	assign.Certificate.Assign(tw.BN254, signer.Certificate.Signature)

	return assign
}
