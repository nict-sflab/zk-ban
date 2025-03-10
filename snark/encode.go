package snark

import (
	"github.com/akakou/zk-ban/encode"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
)

func EncodeCircuit(circuit constraint.ConstraintSystem) ([]byte, error) {
	return encode.EncodeWithWriteTo(circuit)
}

func DecodeCircuit(circuitBytes []byte) (constraint.ConstraintSystem, error) {
	cs := groth16.NewCS(EcCurve)
	err := encode.DecodeWithReadFrom(circuitBytes, cs)
	return cs, err
}

func EncodeProverKey(proveKey groth16.ProvingKey) ([]byte, error) {
	return encode.EncodeWithWriteDump(proveKey)
}

func DecodeProverKey(proveKeyBytes []byte) (groth16.ProvingKey, error) {
	proveKey := groth16.NewProvingKey(EcCurve)
	err := encode.DecodeWithReadDump(proveKeyBytes, proveKey)
	return proveKey, err
}

func EncodeVerifierKey(verifyKey groth16.VerifyingKey) ([]byte, error) {
	return encode.EncodeWithWriteTo(verifyKey)
}

func DecodeVerifierKey(verifyKeyBytes []byte) (groth16.VerifyingKey, error) {
	verifyKey := groth16.NewVerifyingKey(EcCurve)
	err := encode.DecodeWithReadFrom(verifyKeyBytes, verifyKey)
	return verifyKey, err

}

func EncodeProof(proof groth16.Proof) ([]byte, error) {
	return encode.EncodeWithWriteTo(proof)
}

func DecodeProof(proofBytes []byte) (groth16.Proof, error) {
	proof := groth16.NewProof(EcCurve)
	err := encode.DecodeWithReadFrom(proofBytes, proof)
	return proof, err
}
