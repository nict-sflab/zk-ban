package load

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/akakou/zk-ban/dump"
	"github.com/akakou/zk-ban/snark"
	"github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
)

// Loading key safe or unsafe (without checking consistency)
// If you are sure about the key correctness, you can set it to false for performance.
var KeyChecked = false

func LoadUserKey(name, protocol string) (*snark.SnarkProver, error) {
	csFileName := dump.FileName(name, protocol, dump.CircuitFileNameFormat)
	keyFileName := dump.FileName(name, protocol, dump.ProverKeyFileNameFormat)

	pkBin, err := os.ReadFile(keyFileName)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(pkBin)
	key := groth16.NewProvingKey(snark.EcCurve)
	if KeyChecked {
		_, err = key.UnsafeReadFrom(buf)
	} else {
		_, err = key.ReadFrom(buf)
	}

	if err != nil {
		return nil, err
	}

	csBin, err := os.ReadFile(csFileName)
	if err != nil {
		return nil, err
	}
	csBuf := bytes.NewBuffer(csBin)

	cs := groth16.NewCS(snark.EcCurve)
	_, err = cs.ReadFrom(csBuf)
	if err != nil {
		return nil, err
	}

	prover := snark.SnarkProver{
		ConstraintSystem: cs,
		ProveKey:         key,
	}

	return &prover, nil
}

func LoadUserBasicKey(protocol string) (*snark.SnarkProver, error) {
	return LoadUserKey("", protocol)
}

func LoadGroupManagerUpdateKeys() ([]*snark.SizedSnarkVerifier, error) {
	verifiers := []*snark.SizedSnarkVerifier{}
	files, err := os.ReadDir(dump.KeyPath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		fileName := f.Name()
		if !strings.HasSuffix(fileName, ".meta.json") {
			continue
		}

		var verifierRe = regexp.MustCompile(`^(.+)_verifier-(.+)\.meta\.json$`)

		m := verifierRe.FindStringSubmatch(fileName)
		if len(m) != 3 {
			return nil, fmt.Errorf("unexpected format: %s", fileName)
		}

		name := m[2]
		verifier, err := LoadGroupManagerUpdateKey(name)
		if err != nil {
			return nil, err
		}

		verifiers = append(verifiers, verifier)
	}

	return verifiers, nil
}

func LoadGroupManagerMetadata(name string) (*witness.RevocationListSize, error) {
	metaDataFileName := dump.FileName(name, "update", dump.MetaFileNameFormat)

	metadataFile, err := os.ReadFile(metaDataFileName)
	if err != nil {
		return nil, err
	}

	var metadata witness.RevocationListSize
	err = json.Unmarshal(metadataFile, &metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

func LoadGroupManagerUpdateKey(name string) (*snark.SizedSnarkVerifier, error) {
	key, err := LoadGroupManagerKey(name, "update")
	if err != nil {
		return nil, err
	}

	metadata, err := LoadGroupManagerMetadata(name)
	if err != nil {
		return nil, err
	}

	verifier := snark.SizedSnarkVerifier{
		VerifyKey: key,
		RLSize:    *metadata,
		Name:      name,
	}

	return &verifier, nil
}

func LoadBasicGroupManagerKey(protocol string) (*groth16.VerifyingKey, error) {
	return LoadGroupManagerKey("", protocol)
}

func LoadGroupManagerKey(name, protocol string) (*groth16.VerifyingKey, error) {
	keyFileName := dump.FileName(name, protocol, dump.VerifierKeyFileNameFormat)

	vkBin, err := os.ReadFile(keyFileName)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(vkBin)

	var key = groth16.NewVerifyingKey(snark.EcCurve)
	if KeyChecked {
		_, err = key.UnsafeReadFrom(buf)
	} else {
		_, err = key.ReadFrom(buf)
	}

	if err != nil {
		return nil, err
	}

	return &key, nil
}
