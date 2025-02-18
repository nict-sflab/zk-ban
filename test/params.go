package zkbantest

import (
	"math/big"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	"github.com/akakou/zk-ban/witness"
)

var SessionSize = 270
var RevokedNymsPerSession = 130

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

type TestParams struct {
	gpk        *witness.GroupPublicKey
	gsk        *witness.GroupSecretKey
	upk        *witness.UserPublicKey
	usk        *witness.UserSecretKey
	m          *big.Int
	period     *big.Int
	cnt        *big.Int
	cert       *witness.Credential
	nextPeriod *big.Int
}

func (params *TestParams) signer() *witness.Signer {
	signer := witness.Signer{
		UserSecretKey:  params.usk,
		UserPublicKey:  params.upk,
		Credential:     params.cert,
		Period:         params.period,
		GroupPublicKey: params.gpk,
	}

	return &signer
}

func prepareCircuit(rl witness.RevocationList, omitJoinAndSign bool) (*snark.SnarkParams, *snark.SnarkParams, *snark.SnarkParams) {
	var err error
	var joinSnark *snark.SnarkParams = nil
	var signSnark *snark.SnarkParams = nil

	if !omitJoinAndSign {
		joinSnark, err = snark.InitSNARK(&circuit.JoinRequestCircuit{})
		panicIfErr(err)

		signSnark, err = snark.InitSNARK(&circuit.SignCircuit{})
		panicIfErr(err)
	}

	witnessRL := circuit.NewRevocationListWitness(rl)

	updateSnark, err := snark.InitSNARK(&circuit.UpdateCircuit{
		RevocationList: witnessRL,
	})

	panicIfErr(err)

	return joinSnark, signSnark, updateSnark

}

func prepareParams() TestParams {
	gsk, gpk, err := witness.RandomGroupKeyPair()
	panicIfErr(err)

	period := big.NewInt(2024)
	nextPeriod := big.NewInt(2025)

	var usk = &witness.UserSecretKey{
		Number: big.NewInt(102),
	}

	upk, err := usk.PublicKey(period)
	panicIfErr(err)

	m := big.NewInt(100)
	bsn := big.NewInt(101)

	cert, err := gsk.IssueCredential(upk)
	panicIfErr(err)

	return TestParams{
		gpk:        gpk,
		gsk:        gsk,
		usk:        usk,
		upk:        upk,
		m:          m,
		cnt:        bsn,
		period:     period,
		cert:       cert,
		nextPeriod: nextPeriod,
	}
}
