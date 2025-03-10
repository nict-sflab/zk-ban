package highlevel

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"strings"

	core "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	"github.com/akakou/zk-ban/witness"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

type Signature struct {
	Proof *big.Int
	Nym   *big.Int
	Sigma *big.Int
}

func (signature *Signature) String() string {
	proof := signature.Proof.Bytes()
	sigma := signature.Sigma.Bytes()
	nym := signature.Nym.Bytes()

	encodedProof := base64.StdEncoding.EncodeToString(proof)
	encodedSigma := base64.StdEncoding.EncodeToString(sigma)
	encodeNym := base64.StdEncoding.EncodeToString(nym)

	return string(encodedProof) + "." + string(encodedSigma) + "." + string(encodeNym)
}

func (signature *Signature) FromString(str string) error {
	strs := strings.Split(str, ".")
	if len(strs) != 3 {
		return errors.New("invalid signature string (expected 3 parts)")
	}

	proof, err := base64.StdEncoding.DecodeString(strs[0])
	if err != nil {
		return err
	}

	sigma, err := base64.StdEncoding.DecodeString(strs[1])
	if err != nil {
		return err
	}

	nym, err := base64.StdEncoding.DecodeString(strs[2])
	if err != nil {
		return err
	}

	signature.Proof = big.NewInt(0).SetBytes(proof)
	signature.Nym = big.NewInt(0).SetBytes(nym)
	signature.Sigma = big.NewInt(0).SetBytes(sigma)

	return nil
}

type Signer struct {
	Credential     []byte
	Secret         []byte
	UserPublicKey  []byte
	GroupPublicKey []byte
	Period         int64
}

func Sign(m []byte, counter int64, signer, circuitBytes, proveKey []byte) string {
	jsonSignerStruct := Signer{}
	err := json.Unmarshal(signer, &jsonSignerStruct)
	if err != nil {
		return NewResult("", err).String()
	}

	cs, err := snark.DecodeCircuit(circuitBytes)
	if err != nil {
		return NewResult("", err).String()
	}

	proverKeyObj, err := snark.DecodeProverKey(proveKey)
	if err != nil {
		return NewResult("", err).String()
	}

	mBig := big.NewInt(0).SetBytes(m)
	counterBig := big.NewInt(counter)
	periodBig := big.NewInt(jsonSignerStruct.Period)
	upkBig := big.NewInt(0).SetBytes(jsonSignerStruct.UserPublicKey)
	secretBig := big.NewInt(0).SetBytes(jsonSignerStruct.Secret)

	gpkObj, err := zkbanw.GroupPublicKeyFromBytes(jsonSignerStruct.GroupPublicKey)
	if err != nil {
		return NewResult("", err).String()
	}

	snarkProver := snark.SnarkProver{
		ConstraintSystem: cs,
		ProveKey:         proverKeyObj,
	}

	signerStruct := zkbanw.Signer{
		UserSecretKey:  &zkbanw.UserSecretKey{Number: secretBig},
		UserPublicKey:  &zkbanw.UserPublicKey{Number: upkBig},
		Credential:     &zkbanw.Credential{Signature: jsonSignerStruct.Credential},
		Period:         periodBig,
		GroupPublicKey: gpkObj,
	}

	proof, assign, err := core.Sign(mBig, counterBig, &signerStruct, &snarkProver)
	if err != nil {
		return NewResult("", err).String()
	}

	proofBytes, err := snark.EncodeProof(proof)
	pi := big.NewInt(0).SetBytes(proofBytes)
	nym := assign.Nym.(*big.Int)
	signature := assign.Signature.(*big.Int)

	result := Signature{
		Proof: pi,
		Nym:   nym,
		Sigma: signature,
	}

	return NewResult(result.String(), nil).String()
}

func Verify(proof string, m []byte, counter, period int64, gpk, circuitBytes, verifyKeyBytes []byte) string {
	sig := Signature{}
	sig.FromString(proof)

	proofObj, err := snark.DecodeProof(sig.Proof.Bytes())
	if err != nil {
		return NewResult("", err).String()
	}

	verifyKeyObj, err := snark.DecodeVerifierKey(verifyKeyBytes)
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

	wit, err := frontend.NewWitness(&assign, snark.EcCurve.ScalarField())
	if err != nil {
		return NewResult("", err).String()
	}

	pubWit, err := wit.Public()
	if err != nil {
		return NewResult("", err).String()
	}

	err = groth16.Verify(proofObj, verifyKeyObj, pubWit)
	if err != nil {
		return NewResult("", err).String()
	}

	return NewResult("ok", nil).String()
}
