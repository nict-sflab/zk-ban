package highlevel

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
)

func Verify(signatureBytes string, circuitBytes, verifyKeyBytes []byte) []byte {
	proofs := strings.Split(signatureBytes, ".")

	if len(proofs) != 2 {
		return NewResult("", fmt.Errorf("invalid signature (divide with .)")).Bytes()
	}

	proofString, err := base64.StdEncoding.DecodeString(proofs[0])
	if err != nil {
		return NewResult("", err).Bytes()
	}

	witString, _ := base64.StdEncoding.DecodeString(proofs[1])
	if err != nil {
		return NewResult("", err).Bytes()
	}

	cs := groth16.NewCS(snark.EcCurve)
	_, err = cs.ReadFrom(bytes.NewReader(circuitBytes))
	if err != nil {
		return NewResult("", err).Bytes()
	}

	verifyKeyStruct := groth16.NewVerifyingKey(snark.EcCurve)
	_, err = verifyKeyStruct.ReadFrom(bytes.NewReader(verifyKeyBytes))
	if err != nil {
		return NewResult("", err).Bytes()
	}

	proofStruct := groth16.NewProof(snark.EcCurve)
	_, err = proofStruct.ReadFrom(bytes.NewReader(proofString))
	if err != nil {
		return NewResult("", err).Bytes()
	}

	// schema, err := frontend.NewSchema(&circuit.SignCircuit{})
	// if err != nil {
	// 	return NewResult("", err).Bytes()
	// }
	wit, err := witness.New(snark.EcCurve.ScalarField())
	if err != nil {
		return NewResult("", err).Bytes()
	}

	_, err = wit.ReadFrom(bytes.NewReader(witString))
	if err != nil {
		return NewResult("", err).Bytes()
	}

	fmt.Printf("wit: %+v\n", wit)

	// pubWit, err := wit.Public()
	err = groth16.Verify(proofStruct, verifyKeyStruct, wit)
	fmt.Println("err: ", err)
	if err != nil {
		return NewResult("", err).Bytes()
	}

	return NewResult("ok", nil).Bytes()
}
