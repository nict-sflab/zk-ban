package load

import (
	"encoding/json"

	gnarkserializable "github.com/akakou/gnark-serializable"
	"github.com/akakou/zk-ban/snark"
)

func DocodeProver(buf []byte) (*snark.SnarkProver, error) {
	prover := snark.SnarkProver{}

	err := json.Unmarshal(buf, &prover)
	return &prover, err
}

func DocodeVerifyingKey(buf []byte) (*gnarkserializable.VerifyingKey, error) {
	verifyingKey := gnarkserializable.VerifyingKey{}

	err := json.Unmarshal(buf, &verifyingKey)
	return &verifyingKey, err
}

func DocodeSizedVerifyingKey(buf []byte) (*snark.SizedSnarkVerifier, error) {
	verifyingKey := snark.SizedSnarkVerifier{}

	err := json.Unmarshal(buf, &verifyingKey)
	return &verifyingKey, err
}
