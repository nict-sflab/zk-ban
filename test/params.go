package test

import (
	"fmt"
	"os"

	gnarkserializable "github.com/akakou/gnark-serializable"
	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/dump"
	"github.com/akakou/zk-ban/load"
	"github.com/akakou/zk-ban/precomputes"
	"github.com/akakou/zk-ban/primitives"
	"github.com/akakou/zk-ban/snark"
	"github.com/akakou/zk-ban/witness"
)

var TestKeyPath = "./"

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

func PrepareCircuit(rl witness.RevocationList) (*snark.SnarkParams, *snark.SnarkParams, *snark.SnarkParams) {
	var err error
	var joinSnark *snark.SnarkParams = nil
	var signSnark *snark.SnarkParams = nil

	joinSnark, err = snark.InitSNARK(&circuit.JoinRequestCircuit{})
	PanicIfErr(err)

	signSnark, err = snark.InitSNARK(&circuit.SignCircuit{})
	PanicIfErr(err)

	witnessRL := circuit.NewRevocationListAssigned(rl)

	updateSnark, err := snark.InitSNARK(&precomputes.UpdateCircuit{
		circuit.UpdateCircuit{
			RevocationList: witnessRL,
		},
	})

	return joinSnark, signSnark, updateSnark
}

type fileReader struct{}

func (*fileReader) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func PrepareUpdateKey(rlSize witness.RevocationListSize, _ string) (*snark.SnarkProver, *snark.SizedSnarkVerifier) {
	rl := circuit.NewRevocationListAssigned(witness.EmptyRevocationList(rlSize))
	updateSnark, err := snark.InitSNARK(&precomputes.UpdateCircuit{
		circuit.UpdateCircuit{
			RevocationList: rl,
		},
	})

	PanicIfErr(err)

	return updateSnark.Prover(), &snark.SizedSnarkVerifier{
		RLSize:    rlSize,
		VerifyKey: &gnarkserializable.VerifyingKey{updateSnark.VerifyKey},
	}
}

func PrepareUpdateKeyCached(rlSize witness.RevocationListSize, name string) (*snark.SnarkProver, *snark.SizedSnarkVerifier) {
	proverFileName := fmt.Sprintf(dump.UpdateProverKeyFileNameFormat, name)
	verifierFileName := fmt.Sprintf(dump.UpdateVerifierKeyFileNameFormat, name)

	_, err := os.Stat(proverFileName)
	if err != nil {
		dump.DumpUpdateKeys(name, rlSize, TestKeyPath)
	} else {
		fmt.Println("compile skip")
	}

	pk, err := load.LoadUpdateKey(proverFileName, os.DirFS(TestKeyPath), load.DocodeProver)
	PanicIfErr(err)

	vk, err := load.LoadUpdateKey(verifierFileName, os.DirFS(TestKeyPath), load.DocodeSizedVerifyingKey)
	PanicIfErr(err)

	return *pk, *vk
}

func PrepareParams() TestParams {
	gsk, gpk, err := witness.RandomGroupKeyPair()
	PanicIfErr(err)

	period := int64(20240101)
	nextPeriod := int64(20250101)
	var usk = &witness.UserSecretKey{
		primitives.RandBigInt(),
	}

	upk, err := usk.PublicKey(period)
	PanicIfErr(err)

	m := witness.MimcInitBigInt()

	cnt := int64(2)

	cert, err := gsk.IssueCredential(upk)
	PanicIfErr(err)

	return TestParams{
		GPK:        gpk,
		GSK:        gsk,
		USK:        usk,
		UPK:        upk,
		M:          m,
		CNT:        cnt,
		Period:     period,
		Cert:       cert,
		NextPeriod: nextPeriod,
	}
}

func EmptyUniformRevocationList(sessionSize, nymNum int) (witness.RevocationList, witness.RevocationListSize) {
	nymNumPerSession := nymNum / sessionSize

	witness.InitBigInt = witness.MimcInitBigInt
	rlSize := witness.MakeUniformRLSize(sessionSize, nymNumPerSession)

	return witness.EmptyRevocationList(rlSize), rlSize
}

func EmptyProportionalRevocationList(sessionSize, nymNum int) (witness.RevocationList, witness.RevocationListSize) {
	witness.InitBigInt = witness.MimcInitBigInt
	max := nymNum * 2 / sessionSize
	rlSize := witness.MakeLinearRLSize(sessionSize, max, 1)
	rlSize = witness.AjustRLSize(rlSize, nymNum)

	return witness.EmptyRevocationList(rlSize), rlSize
}

const GaussianStandarDeviationDiv = 4

func EmptyGaussianRevocationList(sessionSize, nymNum int) (witness.RevocationList, witness.RevocationListSize) {
	witness.InitBigInt = witness.MimcInitBigInt
	sd := float64(sessionSize) / GaussianStandarDeviationDiv
	rlSize := witness.MakeGaussianRLSize(sessionSize, nymNum, sd)
	rlSize = witness.AjustRLSize(rlSize, nymNum)

	return witness.EmptyRevocationList(rlSize), rlSize
}
