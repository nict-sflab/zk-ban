package highlevel

import (
	"bytes"
	"encoding/json"
	"math/big"

	core "github.com/akakou/zk-ban"
	zkbanc "github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

type Signature struct {
	Circuit []byte
	Proof   []byte
}

type Signer struct {
	Credential     []byte
	Secret         []byte
	UserPublicKey  []byte
	GroupPublicKey []byte
	Period         int64
}

func Sign(m []byte, counter int64, signer, circuitBytes, proveKey []byte) []byte {
	jsonSignerStruct := Signer{}
	err := json.Unmarshal(signer, &jsonSignerStruct)
	if err != nil {
		return NewResult([]byte{}, err).Bytes()
	}

	cs := groth16.NewCS(snark.EcCurve)
	csReader := bytes.NewReader(circuitBytes)
	if _, err := cs.ReadFrom(csReader); err != nil {
		return NewResult([]byte{}, err).Bytes()
	}

	proveKeyStruct := groth16.NewProvingKey(snark.EcCurve)
	proveKeyReader := bytes.NewReader(proveKey)
	if err := proveKeyStruct.ReadDump(proveKeyReader); err != nil {
		return NewResult([]byte{}, err).Bytes()
	}

	mBig := big.NewInt(0).SetBytes(m)
	counterBig := big.NewInt(counter)
	periodBig := big.NewInt(jsonSignerStruct.Period)
	upkBig := big.NewInt(0).SetBytes(jsonSignerStruct.UserPublicKey)
	secretBig := big.NewInt(0).SetBytes(jsonSignerStruct.Secret)

	_, gpkStruct, err := zkbanw.RandomGroupKeyPair()
	if err != nil {
		return NewResult([]byte{}, err).Bytes()
	}

	_, err = gpkStruct.SetBytes(jsonSignerStruct.GroupPublicKey)
	if err != nil {
		return NewResult([]byte{}, err).Bytes()
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

	proof, pubWit, err := core.Sign(mBig, counterBig, &signerStruct, &snarkProver)
	if err != nil {
		return NewResult([]byte{}, err).Bytes()
	}

	schema, err := frontend.NewSchema(&zkbanc.SignCircuit{})
	if err != nil {
		return NewResult([]byte{}, err).Bytes()
	}

	witBytes, err := pubWit.ToJSON(schema)
	if err != nil {
		return NewResult([]byte{}, err).Bytes()
	}

	var buf bytes.Buffer
	if _, err = proof.WriteRawTo(&buf); err != nil {
		return NewResult([]byte{}, err).Bytes()
	}

	var signature = Signature{
		Circuit: witBytes,
		Proof:   buf.Bytes(),
	}

	result, err := json.Marshal(signature)
	if err != nil {
		return NewResult([]byte{}, err).Bytes()
	}

	return NewResult(result, nil).Bytes()
}
