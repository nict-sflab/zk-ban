package main

import (
	"fmt"
	"math/big"
	"time"

	zkban "github.com/akakou/zk-ban"
	zkbanc "github.com/akakou/zk-ban/circuit"
	snark "github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	signParams, err := snark.InitSNARK(&zkbanc.SignCircuit{})
	panicIfErr(err)

	syncParams, err := snark.InitSNARK(&zkbanc.SyncCircuit{})
	panicIfErr(err)

	usk := zkbanw.UserSecretKey{Number: big.NewInt(1)}
	bsn := big.NewInt(3)
	m := big.NewInt(4)
	period := big.NewInt(2024)

	gsk, gpk, err := zkbanw.RandomGroupKeyPair()
	panicIfErr(err)

	upk, err := usk.PublicKey(period)
	panicIfErr(err)

	cert, err := gsk.IssueCertificate(upk)
	panicIfErr(err)

	t := time.Now()

	signer := &zkbanw.Signer{
		UserSecretKey: &usk,
		UserPublicKey: upk,
		Certificate:   cert,
		Period:        period,
	}

	signProof, witness, err := zkban.Sign(m, bsn, signer, gpk, signParams.Prover())
	panicIfErr(err)

	signWitPub, err := witness.Public()
	panicIfErr(err)
	fmt.Printf("Sign Prove: %vms\n", time.Since(t))

	t = time.Now()
	err = groth16.Verify(signProof, signParams.VerifyKey, signWitPub)
	fmt.Printf("Sign Verify: %vms\n", time.Since(t))
	panicIfErr(err)

	commit2 := [zkbanc.RevocationListSize]*big.Int{}

	for i := 0; i < zkbanc.RevocationListSize; i++ {
		commit2[i] = big.NewInt(1000000)
	}

	syncAssign := zkbanc.NewSyncCircuitWitness(commit2, bsn, &zkbanw.Signer{
		UserSecretKey: &usk,
		UserPublicKey: upk,
		Certificate:   cert,
		Period:        period,
	}, gpk)

	syncWit, err := frontend.NewWitness(syncAssign, snark.EcCurve.ScalarField())
	panicIfErr(err)

	syncWitPub, err := syncWit.Public()
	panicIfErr(err)

	t = time.Now()
	syncProof, err := groth16.Prove(syncParams.ConstraintSystem, syncParams.ProveKey, syncWit)
	panicIfErr(err)
	fmt.Printf("Sync Prove: %v\n", time.Since(t))

	t = time.Now()
	err = groth16.Verify(syncProof, syncParams.VerifyKey, syncWitPub)
	fmt.Printf("Sync Verify: %v\n", time.Since(t))
	panicIfErr(err)
	print("ok")
}
