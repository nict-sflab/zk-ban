package main

import (
	"math/big"

	zkban "github.com/akakou/zk-ban"
	zkbanc "github.com/akakou/zk-ban/circuit"
	snark "github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"

	"github.com/consensys/gnark/backend/groth16"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// signParams, err := snark.InitSNARK(&zkbanc.SignCircuit{})
	// panicIfErr(err)

	syncParams, err := snark.InitSNARK(&zkbanc.SyncCircuit{})
	panicIfErr(err)

	usk := zkbanw.UserSecretKey{Number: big.NewInt(1)}
	// m := big.NewInt(4)

	sessions := [zkbanc.RevocationListSize]*big.Int{}
	for i := 0; i < zkbanc.RevocationListSize; i++ {
		sessions[i] = big.NewInt(int64(i + 1))
	}

	last := big.NewInt(int64(10))
	next := big.NewInt(int64(12))

	gsk, gpk, err := zkbanw.RandomGroupKeyPair()
	panicIfErr(err)

	upk, err := usk.PublicKey(last)
	panicIfErr(err)

	cert, err := gsk.IssueCertificate(upk)
	panicIfErr(err)

	signer := &zkbanw.Signer{
		UserSecretKey: &usk,
		UserPublicKey: upk,
		Certificate:   cert,
		Period:        last,
	}

	commit2 := [zkbanc.RevocationListSize]*big.Int{}

	for i := 0; i < zkbanc.RevocationListSize; i++ {
		commit2[i] = big.NewInt(1000000)
	}

	syncProof, witness, err := zkban.Sync(commit2, next, signer, sessions, gpk, syncParams.Prover())
	panicIfErr(err)

	pubWitness, err := witness.Public()
	panicIfErr(err)

	err = groth16.Verify(syncProof, syncParams.VerifyKey, pubWitness)
	panicIfErr(err)
	print("ok")
}
