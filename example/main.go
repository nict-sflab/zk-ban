package main

import (
	"fmt"
	"math/big"
	"time"

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
	nonce := big.NewInt(2)
	bsn := big.NewInt(3)
	m := big.NewInt(4)
	period := big.NewInt(2024)
	dummy_usk := zkbanw.UserSecretKey{Number: big.NewInt(5)}

	gsk, gpk, err := zkbanw.RandomGroupKeyPair()
	panicIfErr(err)

	upk, err := usk.PublicKey(period)
	panicIfErr(err)

	cert, err := gsk.IssueCertificate(upk)
	panicIfErr(err)

	t := time.Now()

	commit, err := zkbanw.ComputeSignCommit(m, bsn, nonce, &usk)
	panicIfErr(err)

	dummy_commit, err := zkbanw.ComputeSignCommit(m, bsn, nonce, &dummy_usk)
	panicIfErr(err)

	signAssign := &zkbanc.SignCircuit{
		UserPublicKey: upk.Number,
		Nonce:         nonce,
		UserSecretKey: usk.Number,
		Basename:      bsn,
		Message:       m,
		Commit1:       commit.Commit1,
		Commit2:       commit.Commit2,
		Commit3:       commit.Commit3,
		Period:        period,
	}

	signAssign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	signAssign.Certificate.Assign(snark.TwistededwardsCurve, cert.Signature)

	// witness definition
	signWit, err := frontend.NewWitness(signAssign, snark.EcCurve.ScalarField())
	panicIfErr(err)

	signWitPub, err := signWit.Public()
	panicIfErr(err)

	// groth16: Prove & Verify
	zkproof, err := groth16.Prove(signParams.ConstraintSystem, signParams.ProveKey, signWit)
	panicIfErr(err)
	fmt.Printf("Sign Prove: %vms\n", time.Since(t))

	t = time.Now()
	err = groth16.Verify(zkproof, signParams.VerifyKey, signWitPub)
	fmt.Printf("Sign Verify: %vms\n", time.Since(t))
	panicIfErr(err)

	commit2 := [zkbanc.RevocationListSize]*big.Int{}
	commit3 := [zkbanc.RevocationListSize]*big.Int{}

	for i := 0; i < zkbanc.RevocationListSize; i++ {
		commit2[i] = dummy_commit.Commit2
		commit3[i] = dummy_commit.Commit3
	}

	t = time.Now()

	syncAssign := zkbanc.SyncCircuitWitness(commit2, commit3, &zkbanw.Signer{
		UserSecretKey: &usk,
		UserPublicKey: upk,
		Certificate:   cert,
		Period:        period,
	}, gpk)

	syncWit, err := frontend.NewWitness(syncAssign, snark.EcCurve.ScalarField())
	panicIfErr(err)

	syncWitPub, err := syncWit.Public()
	panicIfErr(err)

	syncProof, err := groth16.Prove(syncParams.ConstraintSystem, syncParams.ProveKey, syncWit)
	panicIfErr(err)
	fmt.Printf("Sync Prove: %vms\n", time.Since(t))

	t = time.Now()
	err = groth16.Verify(syncProof, syncParams.VerifyKey, syncWitPub)
	fmt.Printf("Sync Verify: %vms\n", time.Since(t))
	panicIfErr(err)
	print("ok")
}
