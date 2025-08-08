package snark_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

type Circuit struct {
}

func (circuit *Circuit) Define(api frontend.API) error {
	return nil
}

type Json struct {
	Proof snark.Proof
	A     int
}

func TestEncode(t *testing.T) {
	var circuit Circuit
	r1cs, err := frontend.Compile(
		ecc.BLS12_381.ScalarField(),
		r1cs.NewBuilder,
		&circuit)
	if err != nil {
		log.Fatalf("compile: %v", err)
	}

	pk, _, err := groth16.Setup(r1cs)
	if err != nil {
		log.Fatalf("setup: %v", err)
	}

	assignment := &Circuit{}
	witness, err := frontend.NewWitness(assignment, ecc.BLS12_381.ScalarField())
	if err != nil {
		log.Fatalf("witness: %v", err)
	}

	proof, err := groth16.Prove(r1cs, pk, witness)
	if err != nil {
		log.Fatalf("prove: %v", err)
	}

	proof2 := snark.Proof{proof}
	j := Json{}

	buf, err := json.Marshal(Json{
		Proof: proof2,
		A:     100,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", string(buf))

	err = json.Unmarshal(buf, &j)
	if err != nil {
		panic(err)
	}
}
