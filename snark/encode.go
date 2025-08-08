package snark

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/consensys/gnark/backend/groth16"
)

type Proof struct{ Proof groth16.Proof }

func (proof Proof) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	_, err := proof.Proof.WriteTo(&buffer)
	if err != nil {
		return nil, err
	}

	raw := buffer.Bytes()
	enc := base64.URLEncoding.EncodeToString(raw)
	res := strconv.Quote(enc)

	return []byte(res), nil
}

func (proof *Proof) UnmarshalJSON(buf []byte) error {
	if string(buf) == "null" {
		return nil
	}

	if proof == nil {
		proof = &Proof{}
	}

	fmt.Printf("\n\nthis is msg \n-----\n%v\n-----\n\n", string(buf))

	unq, err := strconv.Unquote(string(buf))
	if err != nil {
		return err
	}

	raw, err := base64.URLEncoding.DecodeString(unq)
	if err != nil {
		return err
	}

	z := groth16.NewProof(EcCurve)
	reader := bytes.NewReader(raw)
	_, err = z.ReadFrom(reader)

	proof.Proof = z

	return err
	// return err
}

// func EncodeCircuit(circuit constraint.ConstraintSystem) ([]byte, error) {
// 	return encode.EncodeWithWriteTo(circuit)
// }

// func DecodeCircuit(circuitBytes []byte) (constraint.ConstraintSystem, error) {
// 	cs := groth16.NewCS(EcCurve)
// 	err := encode.DecodeWithReadFrom(circuitBytes, cs)
// 	return cs, err
// }

// func EncodeProverKey(proveKey groth16.ProvingKey) ([]byte, error) {
// 	return encode.EncodeWithWriteDump(proveKey)
// }

// func DecodeProverKey(proveKeyBytes []byte) (groth16.ProvingKey, error) {
// 	proveKey := groth16.NewProvingKey(EcCurve)
// 	err := encode.DecodeWithReadDump(proveKeyBytes, proveKey)
// 	return proveKey, err
// }

// func EncodeVerifierKey(verifyKey groth16.VerifyingKey) ([]byte, error) {
// 	return encode.EncodeWithWriteTo(verifyKey)
// }

// func DecodeVerifierKey(verifyKeyBytes []byte) (groth16.VerifyingKey, error) {
// 	verifyKey := groth16.NewVerifyingKey(EcCurve)
// 	err := encode.DecodeWithReadFrom(verifyKeyBytes, verifyKey)
// 	return verifyKey, err

// }

// func EncodeProof(proof groth16.Proof) ([]byte, error) {
// 	return encode.EncodeWithWriteTo(proof)
// }

// func DecodeProof(proofBytes []byte) (groth16.Proof, error) {
// 	proof := groth16.NewProof(EcCurve)
// 	err := encode.DecodeWithReadFrom(proofBytes, proof)
// 	return proof, err
// }
