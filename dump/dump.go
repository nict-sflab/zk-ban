package dump

import (
	"fmt"
	"os"

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
	updateProver, updateVerify, err := MakeUpdateCircuit(rlSize)
	if err != nil {
		fmt.Printf("failed to dump update keys")
	}

	proverFileName := fmt.Sprintf(UpdateProverKeyFileNameFormat, name)
	verifierFileName := fmt.Sprintf(UpdateVerifierKeyFileNameFormat, name)

	err = os.WriteFile(path+proverFileName, updateProver, 0644)
	if err != nil {
		fmt.Printf("failed to dump update keys")
	}

	err = os.WriteFile(path+verifierFileName, updateVerify, 0644)
	if err != nil {
		fmt.Printf("failed to dump update keys")
	}
}
