package load

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/akakou/zk-ban/dump"
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/backend/groth16"
)

func LoadUserUpdateKey(name string) (*snark.SnarkProver, error) {
	ccsFileName := dump.FileName(name, dump.UpdateCircuitFileNameFormat)
	keyFileName := dump.FileName(name, dump.UpdateProverKeyFileNameFormat)

	keyFile, err := os.OpenFile(keyFileName, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	key := groth16.NewProvingKey(snark.EcCurve)
	_, err = key.ReadFrom(keyFile)
	if err != nil {
		return nil, err
	}

	ccsFile, err := os.OpenFile(ccsFileName, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	ccs := groth16.NewCS(snark.EcCurve)
	_, err = ccs.ReadFrom(ccsFile)
	if err != nil {
		return nil, err
	}

	prover := snark.SnarkProver{
		ConstraintSystem: ccs,
		ProveKey:         key,
	}

	return &prover, nil

}

func LoadGroupManagerUpdateKeys() ([]*snark.SizedSnarkVerifier, error) {
	verifiers := []*snark.SizedSnarkVerifier{}
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		fileName := f.Name()

		if !strings.HasSuffix(fileName, ".metadata.json") {
			return nil, fmt.Errorf("invalid file name")
		}

		name := strings.TrimPrefix(fileName, ".metadata.json")

		verifier, err := LoadGroupManagerUpdateKey(name)
		if err != nil {
			return nil, err
		}

		verifiers = append(verifiers, verifier)
	}

	return verifiers, nil
}

func LoadGroupManagerUpdateKey(name string) (*snark.SizedSnarkVerifier, error) {
	dirFS := os.DirFS(dump.KeyPath)

	metaDataFileName := dump.FileName(name, dump.UpdateVKMetaFileNameFormat)
	keyFileName := dump.FileName(name, dump.UpdateVerifierKeyFileNameFormat)

	metadataFile, err := fs.ReadFile(dirFS, metaDataFileName)
	if err != nil {
		return nil, err
	}

	var metadata []int
	err = json.Unmarshal(metadataFile, &metadata)
	if err != nil {
		fmt.Printf("failed to marshal rl size")
	}

	keyFile, err := dirFS.Open(keyFileName)
	if err != nil {
		return nil, err
	}

	var key = groth16.NewVerifyingKey(snark.EcCurve)
	_, err = key.ReadFrom(keyFile)
	if err != nil {
		return nil, err
	}

	verifier := snark.SizedSnarkVerifier{
		VerifyKey: &key,
		RLSize:    metadata,
	}

	return &verifier, nil
}
