package test

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

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

type TestParams struct {
	GPK        *witness.GroupPublicKey
	GSK        *witness.GroupSecretKey
	UPK        *witness.UserPublicKey
	USK        *witness.UserSecretKey
	M          *primitives.BigInt
	Period     int64
	CNT        int64
	Cert       *witness.Credential
	NextPeriod int64
}

func (params *TestParams) Signer() *witness.Signer {
	signer := witness.Signer{
		UserSecretKey: params.USK,
		Credential:    params.Cert,
		Period:        params.Period,
	}

	return &signer
}

func PrepareCircuit(rl witness.RevocationList, omitJoinAndSign bool) (*snark.SnarkParams, *snark.SnarkParams, *snark.SnarkParams) {
	var err error
	var joinSnark *snark.SnarkParams = nil
	var signSnark *snark.SnarkParams = nil

	if !omitJoinAndSign {
		joinSnark, err = snark.InitSNARK(&circuit.JoinRequestCircuit{})
		PanicIfErr(err)

		signSnark, err = snark.InitSNARK(&circuit.SignCircuit{})
		PanicIfErr(err)
	}

	witnessRL := circuit.NewRevocationListAssigned(rl)

	updateSnark, err := snark.InitSNARK(&precomputes.UpdateCircuit{
		circuit.UpdateCircuit{
			RevocationList: witnessRL,
		},
	})

	PanicIfErr(err)

	return joinSnark, signSnark, updateSnark

}

func PrepareParams() TestParams {
	gsk, gpk, err := witness.RandomGroupKeyPair()
	PanicIfErr(err)

	period := int64(2024)
	nextPeriod := int64(2025)
	var usk = &witness.UserSecretKey{
		primitives.NewBigInt(102),
	}

	upk, err := usk.PublicKey(period)
	PanicIfErr(err)

	m := primitives.BigInt{*big.NewInt(100)}
	cnt := int64(101)

	cert, err := gsk.IssueCredential(upk)
	PanicIfErr(err)

	return TestParams{
		GPK:        gpk,
		GSK:        gsk,
		USK:        usk,
		UPK:        upk,
		M:          &m,
		CNT:        cnt,
		Period:     period,
		Cert:       cert,
		NextPeriod: nextPeriod,
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
