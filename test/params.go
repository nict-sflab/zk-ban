package zkbantest

import (
	"math/big"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/precomputes"
	"github.com/akakou/zk-ban/primitives"
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
	m          *primitives.BigInt
	period     int64
	cnt        int64
	cert       *witness.Credential
	nextPeriod int64
}

func (params *TestParams) signer() *witness.Signer {
	signer := witness.Signer{
		UserSecretKey: params.usk,
		Credential:    params.cert,
		Period:        params.period,
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

	witnessRL := circuit.NewRevocationListAssigned(rl)

	updateSnark, err := snark.InitSNARK(&precomputes.UpdateCircuit{
		circuit.UpdateCircuit{
			RevocationList: witnessRL,
		},
	})

	panicIfErr(err)

	return joinSnark, signSnark, updateSnark

}

func prepareParams() TestParams {
	gsk, gpk, err := witness.RandomGroupKeyPair()
	panicIfErr(err)

	period := int64(2024)
	nextPeriod := int64(2025)
	var usk = &witness.UserSecretKey{
		primitives.NewBigInt(102),
	}

	upk, err := usk.PublicKey(period)
	panicIfErr(err)

	m := primitives.BigInt{*big.NewInt(100)}
	cnt := int64(101)

	cert, err := gsk.IssueCredential(upk)
	panicIfErr(err)

	return TestParams{
		gpk:        gpk,
		gsk:        gsk,
		usk:        usk,
		upk:        upk,
		m:          &m,
		cnt:        cnt,
		period:     period,
		cert:       cert,
		nextPeriod: nextPeriod,
	}
}

func EmptyUniformRevocationList(sessionSize, nymNum int) witness.RevocationList {
	nymNumPerSession := nymNum / sessionSize

	witness.InitBigInt = witness.MimcInitBigInt
	rlSize := witness.MakeUniformRLSize(sessionSize, nymNumPerSession)

	return witness.EmptyRevocationList(rlSize)
}

func EmptyProportionalRevocationList(sessionSize, nymNum int) witness.RevocationList {
	witness.InitBigInt = witness.MimcInitBigInt
	max := nymNum * 2 / sessionSize
	rlSize := witness.MakeLinearRLSize(sessionSize, max, 1)
	rlSize = witness.AjustRLSize(rlSize, nymNum)

	return witness.EmptyRevocationList(rlSize)
}

const GaussianStandarDeviationDiv = 4

func EmptyGaussianRevocationList(sessionSize, nymNum int) witness.RevocationList {
	witness.InitBigInt = witness.MimcInitBigInt
	sd := float64(sessionSize) / GaussianStandarDeviationDiv
	rlSize := witness.MakeGaussianRLSize(sessionSize, nymNum, sd)
	rlSize = witness.AjustRLSize(rlSize, nymNum)

	return witness.EmptyRevocationList(rlSize)
}
