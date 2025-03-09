package highlevel

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	"github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

func Verify(proof string, m []byte, counter, period int64, gpk, circuitBytes, verifyKeyBytes []byte) string {
	sig := Signature{}
	sig.FromString(proof)

	proofStruct := groth16.NewProof(snark.EcCurve)
	_, err := proofStruct.ReadFrom(bytes.NewReader(sig.Proof.Bytes()))
	if err != nil {
		return NewResult("", err).String()
	}

	cs := groth16.NewCS(snark.EcCurve)
	_, err = cs.ReadFrom(bytes.NewReader(circuitBytes))
	if err != nil {
		return NewResult("", err).String()
	}

	verifyKeyStruct := groth16.NewVerifyingKey(snark.EcCurve)
	_, err = verifyKeyStruct.ReadFrom(bytes.NewReader(verifyKeyBytes))
	if err != nil {
		return NewResult("", err).String()
	}

	mBig := big.NewInt(0).SetBytes(m)
	counterBig := big.NewInt(counter)
	periodBig := big.NewInt(period)

	publicKey := eddsa.PublicKey{}
	publicKey.Assign(snark.TwistededwardsCurve, gpk)

	dummy := eddsa.Signature{}
	dummyBuf := [32]byte{}
	dummy.Assign(snark.TwistededwardsCurve, dummyBuf[:])

	assign := circuit.SignCircuit{
		UserSecretKey: big.NewInt(0),
		CredentialAuthInfo: circuit.CredentialAuthInfo{
			Credential: dummy,
			Period:     period,
		},
		GroupPublicKey: publicKey,
		SessionTag:     witness.SessionTag(counterBig, periodBig),
		Message:        mBig,
		Signature:      sig.Sigma,
		Nym:            sig.Nym,
	}

	fmt.Print(
		big.NewInt(0),
		circuit.CredentialAuthInfo{
			Credential: eddsa.Signature{},
			Period:     period,
		},
		publicKey,
		witness.SessionTag(counterBig, periodBig),
		mBig,
		sig.Nym,
		sig.Sigma,
	)

	wit, err := frontend.NewWitness(&assign, snark.EcCurve.ScalarField())
	if err != nil {
		return NewResult("", err).String()
	}

	pubWit, err := wit.Public()
	if err != nil {
		return NewResult("", err).String()
	}

	err = groth16.Verify(proofStruct, verifyKeyStruct, pubWit)
	if err != nil {
		return NewResult("", err).String()
	}

	return NewResult("ok", nil).String()
}
