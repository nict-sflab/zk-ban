package zkban_test

import (
	"math/big"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/commit"
	"github.com/akakou/zk-ban/snark"
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
	sessionName    [circuit.SessionSize]*big.Int
	revocationList [circuit.SessionSize][circuit.RevokedNymsPerSession]*big.Int
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

	updateSnark, err := snark.InitSNARK(&circuit.UpdateCircuit{})
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
		sessionName:    sn,
		revocationList: rl,
	}
}
