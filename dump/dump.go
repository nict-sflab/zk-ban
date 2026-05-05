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

type DumpKeys func(name, protocol string, params *snark.SnarkParams)

func dumpKeys(name, protocol string, params *snark.SnarkParams, secure bool) {
	csFileName := FileName(name, protocol, CircuitFileNameFormat, nil)
	verifierFileName := FileName(name, protocol, VerifierKeyFileNameFormat, nil)

	var export *string = nil
	if secure {
		tmp := "export"
		export = &tmp
	}

	proverFileName := FileName(name, protocol, ProverKeyFileNameFormat, export)

	proverKeyFile, err := os.Create(proverFileName)
	if err != nil {
		log.Fatalf("create pk file: %v", err)
	}
	defer proverKeyFile.Close()

	if secure {
		_, err = params.ProveKey.WriteTo(proverKeyFile)
		if err != nil {
			log.Fatalf("write dump pk: %v", err)
		}
	} else {
		err = params.ProveKey.WriteDump(proverKeyFile)
		if err != nil {
			log.Fatalf("write dump pk: %v", err)
		}
	}

	constraintSystemFile, err := os.Create(csFileName)
	if err != nil {
		log.Fatalf("create cs file: %v", err)
	}
	defer constraintSystemFile.Close()

	_, err = params.ConstraintSystem.WriteTo(constraintSystemFile)
	if err != nil {
		log.Fatalf("write cs: %v", err)
	}

	verifierKeyFile, err := os.Create(verifierFileName)
	if err != nil {
		log.Fatalf("create vk file: %v", err)
	}
	defer verifierKeyFile.Close()

	_, err = params.VerifyKey.WriteRawTo(verifierKeyFile)
	if err != nil {
		log.Fatalf("write raw vk: %v", err)
	}
}

func DumpUnsafeKeys(name, protocol string, params *snark.SnarkParams) {
	dumpKeys(name, protocol, params, false)
}

func DumpSafeKeys(name, protocol string, params *snark.SnarkParams) {
	dumpKeys(name, protocol, params, true)
}

func DumpBasicKeys(dumpKey DumpKeys) {
	joinParams, err := snark.InitSNARK(&circuit.JoinRequestCircuit{})
	if err != nil {
		log.Fatalf("init join snark: %v", err)
	}

	signParams, err := snark.InitSNARK(&circuit.SignCircuit{})
	if err != nil {
		log.Fatalf("init sign snark: %v", err)
	}

	dumpKey("", "join", joinParams)
	dumpKey("", "sign", signParams)
}

func DumpMetadata(name string, rlSize witness.RevocationListSize) {
	verifierMetaFileName := FileName(name, "update", MetaFileNameFormat, nil)
	metadataFile, err := os.Create(verifierMetaFileName)
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

func DumpUpdateKeys(name string, rlSize witness.RevocationListSize, dumpKeys DumpKeys) {
	rl := witness.EmptyRevocationList(rlSize)
	params, err := snark.InitSNARK(&precomputes.UpdateCircuit{
		UpdateCircuit: circuit.UpdateCircuit{
			RevocationList: circuit.NewRevocationListAssigned(rl),
		},
	})

	if err != nil {
		log.Fatalf("write vk metadata: %v", err)
	}

	dumpKeys(name, "update", params)
	DumpMetadata(name, rlSize)
}
