package zkban

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/commit"
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

type TestParams struct {
	gpk            *commit.GroupPublicKey
	gsk            *commit.GroupSecretKey
	joinSnark      *snark.SnarkParams
	signSnark      *snark.SnarkParams
	updateSnark    *snark.SnarkParams
	upk            *commit.UserPublicKey
	usk            *commit.UserSecretKey
	m              *big.Int
	period         *big.Int
	bsn            *big.Int
	cert           *commit.Credential
	nextPeriod     *big.Int
	sessionName    []frontend.Variable
	revocationList [][]frontend.Variable
}

func (params *TestParams) signer() *commit.Signer {
	signer := commit.Signer{
		UserSecretKey: params.usk,
		UserPublicKey: params.upk,
		Credential:    params.cert,
		Period:        params.period,
	}

	return &signer
}

func prepare() TestParams {
	gsk, gpk, err := commit.RandomGroupKeyPair()
	panicIfErr(err)

	joinSnark, err := snark.InitSNARK(&circuit.JoinRequestCircuit{})
	panicIfErr(err)

	signSnark, err := snark.InitSNARK(&circuit.SignCircuit{})
	panicIfErr(err)

	tags, nyms := circuit.NewDefaultRevocationListWitness(
		[circuit.SessionSize][circuit.RevokedNymsPerSession]*big.Int{},
		[circuit.SessionSize]*big.Int{},
	)

	fmt.Printf("nyms: %v\n", len(nyms))
	fmt.Printf("nyms[0]: %v\n", len(nyms[0]))
	fmt.Printf("tags: %v\n", len(tags))

	updateSnark, err := snark.InitSNARK(&circuit.UpdateCircuit{
		RevocationList: circuit.RevocationList{
			RevokedNyms: nyms,
			SessionTags: tags,
		},
	})
	panicIfErr(err)

	period := big.NewInt(2024)
	nextPeriod := big.NewInt(2025)

	var usk = &commit.UserSecretKey{
		Number: big.NewInt(100),
	}

	upk, err := usk.PublicKey(period)
	panicIfErr(err)

	m := big.NewInt(100)
	bsn := big.NewInt(100)

	cert, err := gsk.IssueCredential(upk)
	panicIfErr(err)

	sn := [circuit.SessionSize]*big.Int{}

	for i := 0; i < circuit.SessionSize; i++ {
		sn[i] = big.NewInt(300)
	}

	rl := [circuit.SessionSize][circuit.RevokedNymsPerSession]*big.Int{}

	for i := 0; i < circuit.SessionSize; i++ {
		for j := 0; j < circuit.RevokedNymsPerSession; j++ {
			rl[i][j] = big.NewInt(300)
		}
	}

	tags_, nyms_ := circuit.NewDefaultRevocationListWitness(rl, sn)

	return TestParams{
		gpk:            gpk,
		gsk:            gsk,
		usk:            usk,
		upk:            upk,
		joinSnark:      joinSnark,
		signSnark:      signSnark,
		updateSnark:    updateSnark,
		m:              m,
		bsn:            bsn,
		period:         period,
		cert:           cert,
		nextPeriod:     nextPeriod,
		sessionName:    tags_,
		revocationList: nyms_,
	}
}

func TestAll(t *testing.T) {
	assert := test.NewAssert(t)
	params := prepare()

	t.Run("join req", func(t *testing.T) {
		var proof groth16.Proof
		proof, pubWit, _, _, err := JoinRequest(params.period, params.joinSnark.Prover())
		assert.NoError(err)

		err = groth16.Verify(proof, params.joinSnark.VerifyKey, pubWit)
		assert.NoError(err)
	})

	t.Run("sign", func(t *testing.T) {
		proof, pubWit, err := Sign(params.m, params.bsn, params.signer(), params.gpk, params.signSnark.Prover())
		assert.NoError(err)

		err = groth16.Verify(proof, params.signSnark.VerifyKey, pubWit)
		assert.NoError(err)
	})

	t.Run("update", func(t *testing.T) {
		_, proof, pubWit, err := Update(params.nextPeriod, params.signer(), params.revocationList, params.sessionName, params.gpk, params.updateSnark.Prover())
		assert.NoError(err)

		err = groth16.Verify(proof, params.updateSnark.VerifyKey, pubWit)
		assert.NoError(err)
	})
}

func BenchmarkAll(t *testing.B) {
	params := prepare()

	var proof groth16.Proof
	var pubWit witness.Witness
	var err error

	t.Run("join req", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			proof, pubWit, _, _, err = JoinRequest(params.period, params.joinSnark.Prover())
			panicIfErr(err)
		}
	})

	t.Run("verify join req", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err = groth16.Verify(proof, params.joinSnark.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})

	t.Run("sign", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			proof, pubWit, err = Sign(params.m, params.bsn, params.signer(), params.gpk, params.signSnark.Prover())
			panicIfErr(err)
		}
	})

	t.Run("verify", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err = groth16.Verify(proof, params.signSnark.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})

	t.Run("update-req", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _, _, err := Update(params.nextPeriod, params.signer(), params.revocationList, params.sessionName, params.gpk, params.updateSnark.Prover())
			panicIfErr(err)
		}
	})

	t.Run("update-verify", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			params.gsk.IssueCredential(params.upk)

			err = groth16.Verify(proof, params.updateSnark.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})
}
