package dump

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/precomputes"
	"github.com/akakou/zk-ban/snark"
	"github.com/akakou/zk-ban/witness"
)

func DumpBasicKeys(path string) {
	joinProver, joinVerify, err := JoinRequestCircuit()
	if err != nil {
		fmt.Printf("failed to dump join keys")
	}

	err = os.WriteFile(path+"/join_prover.key.json", joinProver, 0644)
	if err != nil {
		fmt.Printf("failed to dump join keys")
	}

	err = os.WriteFile(path+"/join_verifier.key.json", joinVerify, 0644)
	if err != nil {
		fmt.Printf("failed to dump join keys")
	}

	signProver, signVerifyKey, err := SignCircuit()
	if err != nil {
		fmt.Printf("failed to dump sign keys")
	}

	err = os.WriteFile(path+"/sign_prover.key.json", signProver, 0644)
	if err != nil {
		fmt.Printf("failed to dump sign keys")
	}

	err = os.WriteFile(path+"/sign_verifier.key.json", signVerifyKey, 0644)
	if err != nil {
		fmt.Printf("failed to dump sign keys")
	}

}

func DumpUpdateKeys(name string, rlSize witness.RevocationListSize, path string) {
	rl := witness.EmptyRevocationList(rlSize)

	params, err := snark.InitSNARK(&precomputes.UpdateCircuit{
		UpdateCircuit: circuit.UpdateCircuit{
			RevocationList: circuit.NewRevocationListAssigned(rl),
		},
	})

	if err != nil {
		fmt.Printf("failed to generate update keys")
	}

	proverFileName := FileName(name, UpdateProverKeyFileNameFormat)
	csFileName := FileName(name, UpdateCircuitFileNameFormat)
	verifierFileName := FileName(name, UpdateVerifierKeyFileNameFormat)
	verifierMetaFileName := FileName(name, UpdateVKMetaFileNameFormat)

	proverKeyFile, err := os.Create(proverFileName)
	if err != nil {
		log.Fatalf("create pk file: %v", err)
	}
	defer proverKeyFile.Close()

	if _, err := params.ProveKey.WriteRawTo(proverKeyFile); err != nil {
		log.Fatalf("write raw pk: %v", err)
	}

	constraintSystemFile, err := os.Create(csFileName)
	if err != nil {
		log.Fatalf("create cs file: %v", err)
	}
	defer constraintSystemFile.Close()

	if _, err := params.ConstraintSystem.WriteTo(constraintSystemFile); err != nil {
		log.Fatalf("write cs: %v", err)
	}

	verifierKeyFile, err := os.Create(verifierFileName)
	if err != nil {
		log.Fatalf("create vk file: %v", err)
	}

	defer verifierKeyFile.Close()
	if _, err := params.VerifyKey.WriteRawTo(verifierKeyFile); err != nil {
		log.Fatalf("write raw vk: %v", err)
	}

	fMeta, err := os.Create(verifierMetaFileName)
	if err != nil {
		log.Fatalf("create vk metadata file: %v", err)
	}
	defer constraintSystemFile.Close()

	v, err := json.Marshal(rlSize)
	if err != nil {
		fmt.Printf("failed to marshal rl size")
	}

	if _, err := fMeta.Write(v); err != nil {
		log.Fatalf("write vk metadata: %v", err)
	}
}
