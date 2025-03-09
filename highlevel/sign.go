package highlevel

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"strings"

	core "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
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

	cs := groth16.NewCS(snark.EcCurve)
	csReader := bytes.NewReader(circuitBytes)
	_, err = cs.ReadFrom(csReader)
	if err != nil {
		return NewResult("", err).String()
	}

	proveKeyStruct := groth16.NewProvingKey(snark.EcCurve)
	proveKeyReader := bytes.NewReader(proveKey)
	err = proveKeyStruct.ReadDump(proveKeyReader)
	if err != nil {
		return NewResult("", err).String()
	}

	mBig := big.NewInt(0).SetBytes(m)
	counterBig := big.NewInt(counter)
	periodBig := big.NewInt(jsonSignerStruct.Period)
	upkBig := big.NewInt(0).SetBytes(jsonSignerStruct.UserPublicKey)
	secretBig := big.NewInt(0).SetBytes(jsonSignerStruct.Secret)

	_, gpkStruct, err := zkbanw.RandomGroupKeyPair()
	if err != nil {
		return NewResult("", err).String()
	}

	_, err = gpkStruct.SetBytes(jsonSignerStruct.GroupPublicKey)
	if err != nil {
		return NewResult("", err).String()
	}

	snarkProver := snark.SnarkProver{
		ConstraintSystem: cs,
		ProveKey:         proveKeyStruct,
	}

	signerStruct := zkbanw.Signer{
		UserSecretKey:  &zkbanw.UserSecretKey{Number: secretBig},
		UserPublicKey:  &zkbanw.UserPublicKey{Number: upkBig},
		Credential:     &zkbanw.Credential{Signature: jsonSignerStruct.Credential},
		Period:         periodBig,
		GroupPublicKey: gpkStruct,
	}

	proof, assign, err := core.Sign(mBig, counterBig, &signerStruct, &snarkProver)
	if err != nil {
		return NewResult("", err).String()
	}

	var proofBuf bytes.Buffer
	_, err = proof.WriteTo(&proofBuf)
	if err != nil {
		return NewResult("", err).String()
	}

	pi := big.NewInt(0).SetBytes(proofBuf.Bytes())
	nym := assign.Nym.(*big.Int)
	signature := assign.Signature.(*big.Int)

	result := Signature{
		Proof: pi,
		Nym:   nym,
		Sigma: signature,
	}

	return NewResult(result.String(), nil).String()
}
