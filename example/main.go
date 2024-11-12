package main

import (
	"fmt"
	"math/big"
	"time"

	zkbanc "github.com/akakou/zk-ban/circuit"
	zkbanw "github.com/akakou/zk-ban/witness"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	authCircuit := zkbanc.ProofCircuit{}
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &authCircuit)
	panicIfErr(err)

	// groth16 zkSNARK: Setup
	pk, vk, err := groth16.Setup(ccs)
	panicIfErr(err)

	usk := zkbanw.UserSecretKey{UserSecretKey: big.NewInt(1)}
	r := big.NewInt(2)
	m := big.NewInt(3)
	period := big.NewInt(2024)

	gsk, gpk, err := zkbanw.RandomGroupKeyPair()
	panicIfErr(err)

	upk, err := usk.PublicKey(period)
	panicIfErr(err)

	cert, err := gsk.IssuseCertificate(upk)
	panicIfErr(err)

	t := time.Now()

	proof, err := zkbanw.ProveWitness(m, r, &usk)
	panicIfErr(err)

	assign := &zkbanc.ProofCircuit{
		Nonce:         r,
		UserSecretKey: usk.UserSecretKey,
		Message:       m,
		Hash:          proof.Hash,
		UserPublicKey: upk.UserPublicKey,
	}

	assign.GroupPublicKey.Assign(twistededwards.BN254, gpk.Bytes())
	assign.Certificate.Assign(twistededwards.BN254, cert.Signature)

	// witness definition
	witness, err := frontend.NewWitness(assign, ecc.BN254.ScalarField())
	panicIfErr(err)

	publicWitness, err := witness.Public()
	panicIfErr(err)

	// groth16: Prove & Verify
	zkproof, err := groth16.Prove(ccs, pk, witness)
	panicIfErr(err)

	fmt.Printf("Prove: %vms\n", time.Since(t))

	t = time.Now()
	err = groth16.Verify(zkproof, vk, publicWitness)
	fmt.Printf("Verify: %vms\n", time.Since(t))

	panicIfErr(err)

	print("ok")
}
