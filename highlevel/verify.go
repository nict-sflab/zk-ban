package highlevel

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
)

func Verify(signatureBytes, circuitBytes, verifyKeyBytes []byte) []byte {
	var signature Signature
	if err := json.Unmarshal(signatureBytes, &signature); err != nil {
		return NewResult([]byte{}, err).Bytes()
	}

	cs := groth16.NewCS(snark.EcCurve)
	_, err := cs.ReadFrom(bytes.NewReader(circuitBytes))
	if err != nil {
		return NewResult([]byte{}, err).Bytes()
	}

	verifyKeyStruct := groth16.NewVerifyingKey(snark.EcCurve)
	_, err = verifyKeyStruct.ReadFrom(bytes.NewReader(verifyKeyBytes))
	if err != nil {
		return NewResult([]byte{}, err).Bytes()
	}

	proofStruct := groth16.NewProof(snark.EcCurve)
	_, err = proofStruct.ReadFrom(bytes.NewReader(signature.Proof))
	if err != nil {
		return NewResult([]byte{}, err).Bytes()
	}

	schema, err := frontend.NewSchema(&circuit.SignCircuit{})
	if err != nil {
		return NewResult([]byte{}, err).Bytes()
	}
	wit, err := witness.New(snark.EcCurve.ScalarField())
	if err != nil {
		return NewResult([]byte{}, err).Bytes()
	}
	err = wit.FromJSON(schema, signature.Circuit)
	if err != nil {
		return NewResult([]byte{}, err).Bytes()
	}

	fmt.Printf("wit: %+v\n", wit)

	// pubWit, err := wit.Public()
	err = groth16.Verify(proofStruct, verifyKeyStruct, wit)
	fmt.Println("err: ", err)
	if err != nil {
		return NewResult([]byte{}, err).Bytes()
	}

	return NewResult([]byte("ok"), nil).Bytes()
}
