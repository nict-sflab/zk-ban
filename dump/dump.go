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

func DumpKeys(name, protocol string, params *snark.SnarkParams) {
	proverFileName := FileName(name, protocol, ProverKeyFileNameFormat)
	csFileName := FileName(name, protocol, CircuitFileNameFormat)
	verifierFileName := FileName(name, protocol, VerifierKeyFileNameFormat)

	proverKeyFile, err := os.Create(KeyPath + "/" + proverFileName)
	if err != nil {
		log.Fatalf("create pk file: %v", err)
	}
	defer proverKeyFile.Close()

	err = params.ProveKey.WriteDump(proverKeyFile)
	if err != nil {
		log.Fatalf("write dump pk: %v", err)
	}

	constraintSystemFile, err := os.Create(KeyPath + "/" + csFileName)
	if err != nil {
		log.Fatalf("create cs file: %v", err)
	}
	defer constraintSystemFile.Close()

	_, err = params.ConstraintSystem.WriteTo(constraintSystemFile)
	if err != nil {
		log.Fatalf("write cs: %v", err)
	}

	verifierKeyFile, err := os.Create(KeyPath + "/" + verifierFileName)
	if err != nil {
		log.Fatalf("create vk file: %v", err)
	}
	defer verifierKeyFile.Close()

	_, err = params.VerifyKey.WriteRawTo(verifierKeyFile)
	if err != nil {
		log.Fatalf("write raw vk: %v", err)
	}

}

func DumpBasicKeys() {
	joinParams, err := snark.InitSNARK(&circuit.JoinRequestCircuit{})
	if err != nil {
		log.Fatalf("init join snark: %v", err)
	}

	signParams, err := snark.InitSNARK(&circuit.SignCircuit{})
	if err != nil {
		log.Fatalf("init sign snark: %v", err)
	}

	DumpKeys("", "join", joinParams)
	DumpKeys("", "sign", signParams)
}

func DumpMetadata(name string, rlSize witness.RevocationListSize) {
	verifierMetaFileName := FileName(name, "update", MetaFileNameFormat)
	metadataFile, err := os.Create(KeyPath + "/" + verifierMetaFileName)
	if err != nil {
		log.Fatalf("create vk metadata file: %v", err)
	}
	defer metadataFile.Close()

	v, err := json.Marshal(rlSize)
	if err != nil {
		fmt.Printf("failed to marshal rl size")
	}

	_, err = metadataFile.Write(v)
	if err != nil {
		log.Fatalf("write vk metadata: %v", err)
	}
}

func DumpUpdateKeys(name string, rlSize witness.RevocationListSize) {
	rl := witness.EmptyRevocationList(rlSize)
	params, err := snark.InitSNARK(&precomputes.UpdateCircuit{
		UpdateCircuit: circuit.UpdateCircuit{
			RevocationList: circuit.NewRevocationListAssigned(rl),
		},
	})

	if err != nil {
		log.Fatalf("write vk metadata: %v", err)
	}

	DumpKeys(name, "update", params)
	DumpMetadata(name, rlSize)
}
