package highlevel

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"

	core "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
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
		return NewResult("", err).Bytes()
	}

	cs := groth16.NewCS(snark.EcCurve)
	csReader := bytes.NewReader(circuitBytes)
	_, err = cs.ReadFrom(csReader)
	if err != nil {
		return NewResult("", err).Bytes()
	}

	proveKeyStruct := groth16.NewProvingKey(snark.EcCurve)
	proveKeyReader := bytes.NewReader(proveKey)
	err = proveKeyStruct.ReadDump(proveKeyReader)
	if err != nil {
		return NewResult("", err).Bytes()
	}

	mBig := big.NewInt(0).SetBytes(m)
	counterBig := big.NewInt(counter)
	periodBig := big.NewInt(jsonSignerStruct.Period)
	upkBig := big.NewInt(0).SetBytes(jsonSignerStruct.UserPublicKey)
	secretBig := big.NewInt(0).SetBytes(jsonSignerStruct.Secret)

	_, gpkStruct, err := zkbanw.RandomGroupKeyPair()
	if err != nil {
		return NewResult("", err).Bytes()
	}

	_, err = gpkStruct.SetBytes(jsonSignerStruct.GroupPublicKey)
	if err != nil {
		return NewResult("", err).Bytes()
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
		return NewResult("", err).Bytes()
	}

	// schema, err := frontend.NewSchema(&zkbanc.SignCircuit{})
	// if err != nil {
	// 	return NewResult([]byte{}, err).Bytes()
	// }

	var witBuf bytes.Buffer
	_, err = pubWit.WriteTo(&witBuf)
	if err != nil {
		return NewResult("", err).Bytes()
	}
	witString := base64.StdEncoding.EncodeToString(witBuf.Bytes())

	var proofBuf bytes.Buffer
	_, err = proof.WriteTo(&proofBuf)
	if err != nil {
		return NewResult("", err).Bytes()
	}

	proofString := base64.StdEncoding.EncodeToString(proofBuf.Bytes())
	result := fmt.Sprintf("%s.%s", proofString, witString)

	return NewResult(result, nil).Bytes()
}
