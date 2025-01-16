package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
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

func randBig() *big.Int {
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(128), nil).Sub(max, big.NewInt(1))

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic("failed to make random")
	}

	return n
}

func main() {
	// signParams, err := snark.InitSNARK(&zkbanc.SignCircuit{})
	// panicIfErr(err)

	syncParams, err := snark.InitSNARK(&zkbanc.SyncCircuit{})
	panicIfErr(err)

	usk := zkbanw.UserSecretKey{Number: randBig()}
	// m := big.NewInt(4)

	sessions := [zkbanc.SessionSize]*big.Int{}
	for i := 0; i < zkbanc.SessionSize; i++ {
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

	commit2 := [zkbanc.SessionSize][zkbanc.RevocationPerSession]*big.Int{}

	for i := 0; i < zkbanc.SessionSize; i++ {
		for j := 0; j < zkbanc.RevocationPerSession; j++ {
			commit2[i][j] = randBig()
		}
	}

	fmt.Printf("session: %v, revoke: %v\n", zkbanc.SessionSize, zkbanc.RevocationPerSession)

	syncProof, witness, err := zkban.Sync(commit2, next, signer, sessions, gpk, syncParams.Prover())
	panicIfErr(err)

	pubWitness, err := witness.Public()
	panicIfErr(err)

	err = groth16.Verify(syncProof, syncParams.VerifyKey, pubWitness)
	panicIfErr(err)

	b := bytes.Buffer{}
	syncProof.WriteTo(&b)
	fmt.Printf("len: %v\n", b.Len())

	print("ok")
}
