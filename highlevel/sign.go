package highlevel

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func Sign(m []byte, counter, period int64, upk, cred, secret, gpk, circuitBytes, proveKey []byte) []byte {
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
	periodBig := big.NewInt(period)
	upkBig := big.NewInt(0).SetBytes(upk)
	secretBig := big.NewInt(0).SetBytes(secret)

	_, gpkStruct, err := zkbanw.RandomGroupKeyPair()
	if err != nil {
		return NewResult([]byte{}, err).Bytes()
	}
	if _, err = gpkStruct.SetBytes(gpk); err != nil {
		return NewResult([]byte{}, err).Bytes()
	}

	snarkProver := snark.SnarkProver{
		ConstraintSystem: cs,
		ProveKey:         proveKeyStruct,
	}

	signerStruct := zkbanw.Signer{
		UserSecretKey:  &zkbanw.UserSecretKey{Number: secretBig},
		UserPublicKey:  &zkbanw.UserPublicKey{Number: upkBig},
		Credential:     &zkbanw.Credential{Signature: cred},
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

	fmt.Printf("Sign: %v\n", string(witBytes))

	result, err := json.Marshal(signature)
	if err != nil {
		return NewResult([]byte{}, err).Bytes()
	}

	return NewResult(result, nil).Bytes()
}
