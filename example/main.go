package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/commit"
	snark "github.com/akakou/zk-ban/snark"

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
	// signParams, err := snark.InitSNARK(&circuit.SignCircuit{})
	// panicIfErr(err)

	updateParams, err := snark.InitSNARK(&circuit.UpdateCircuit{})
	panicIfErr(err)

	usk := commit.UserSecretKey{Number: randBig()}
	// m := big.NewInt(4)

	sessions := [circuit.SessionSize]*big.Int{}
	for i := 0; i < circuit.SessionSize; i++ {
		sessions[i] = big.NewInt(int64(i + 1))
	}

	last := big.NewInt(int64(10))
	next := big.NewInt(int64(12))

	gsk, gpk, err := commit.RandomGroupKeyPair()
	panicIfErr(err)

	upk, err := usk.PublicKey(last)
	panicIfErr(err)

	cert, err := gsk.IssueCredential(upk)
	panicIfErr(err)

	signer := &commit.Signer{
		UserSecretKey: &usk,
		UserPublicKey: upk,
		Credential:    cert,
		Period:        last,
	}

	commit2 := [circuit.SessionSize][circuit.RevocationPerSession]*big.Int{}

	for i := 0; i < circuit.SessionSize; i++ {
		for j := 0; j < circuit.RevocationPerSession; j++ {
			commit2[i][j] = randBig()
		}
	}

	fmt.Printf("session: %v, revoke: %v\n", circuit.SessionSize, circuit.RevocationPerSession)

	updateProof, witness, err := zkban.Update(commit2, next, signer, sessions, gpk, updateParams.Prover())
	panicIfErr(err)

	pubWitness, err := witness.Public()
	panicIfErr(err)

	err = groth16.Verify(updateProof, updateParams.VerifyKey, pubWitness)
	panicIfErr(err)

	b := bytes.Buffer{}
	updateProof.WriteTo(&b)
	fmt.Printf("len: %v\n", b.Len())

	print("ok")
}
