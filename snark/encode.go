package snark

import (
	"bytes"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
)

func (proof Proof) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	_, err := proof.WriteTo(&buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (proof *Proof) UnmarshalJSON(buf []byte) error {
	if string(buf) == "null" {
		return nil
	}

	var z groth16.Proof
	reader := bytes.NewReader(buf)
	_, err := z.ReadFrom(reader)

	if err != nil {
		proof.Proof = z
	}

	return err
}

func (ccs ConstraintSystem) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	_, err := ccs.WriteTo(&buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (ccs *ConstraintSystem) UnmarshalJSON(buf []byte) error {
	if string(buf) == "null" {
		return nil
	}

	var z constraint.ConstraintSystem
	reader := bytes.NewReader(buf)
	_, err := z.ReadFrom(reader)

	if err != nil {
		ccs.ConstraintSystem = z
	}

	return err
}

func (pk ProveKey) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	_, err := pk.WriteTo(&buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (pk *ProveKey) UnmarshalJSON(buf []byte) error {
	if string(buf) == "null" {
		return nil
	}

	var z groth16.ProvingKey
	reader := bytes.NewReader(buf)
	_, err := z.ReadFrom(reader)

	if err != nil {
		pk.ProvingKey = z
	}

	return err
}

func (vk VerifyKey) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	_, err := vk.WriteTo(&buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (pk *VerifyKey) UnmarshalJSON(buf []byte) error {
	if string(buf) == "null" {
		return nil
	}

	var z groth16.VerifyingKey
	reader := bytes.NewReader(buf)
	_, err := z.ReadFrom(reader)

	if err != nil {
		pk.VerifyingKey = z
	}

	return err
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
