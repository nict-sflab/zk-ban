package test

import (
	"fmt"
	"os"

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
		VerifyKey: &updateSnark.VerifyKey,
	}
}

func PrepareKeyIfNotExist(rlSize witness.RevocationListSize, name string) {
	dump.KeyPath = TestKeyPath

	metadataFileName := dump.FileName(name, "update", dump.MetaFileNameFormat)

	fmt.Printf("search key at %s\n", metadataFileName)
	_, err := os.Stat(metadataFileName)
	if err != nil {
		fmt.Println("compile")
		dump.DumpUpdateKeys(name, rlSize)
	} else {
		fmt.Println("compile skip")
	}

}

func PrepareUpdateKeyCached(rlSize witness.RevocationListSize, name string) (*snark.SnarkProver, *snark.SizedSnarkVerifier) {
	dump.KeyPath = TestKeyPath
	PrepareKeyIfNotExist(rlSize, name)

	pk, err := load.LoadUserKey(name, "update")
	PanicIfErr(err)

	vk, err := load.LoadGroupManagerUpdateKey(name)
	PanicIfErr(err)

	return pk, vk
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

var InitBigInt = witness.MimcInitBigInt

func EmptyUniformRevocationList(periodSize, nymNum int) witness.RevocationList {
	witness.InitBigInt = InitBigInt
	rlSize := witness.MakeUniformRLSizeFromTotal(periodSize, nymNum)
	return witness.EmptyRevocationList(rlSize)
}

func EmptyProportionalRevocationList(periodSize, nymNum int) witness.RevocationList {
	witness.InitBigInt = InitBigInt
	rlSize := witness.MakeProportionalRLSizeFromTotal(periodSize, nymNum)

	return witness.EmptyRevocationList(rlSize)
}

const GaussianStandarDeviationDiv = 4

func EmptyGaussianRevocationListSize(periodSize, nymNum int) witness.RevocationListSize {
	witness.InitBigInt = InitBigInt
	sd := float64(periodSize) / GaussianStandarDeviationDiv
	rlSize := witness.MakeGaussianRLSizeFromTotal(periodSize, nymNum, sd)
	return rlSize
}

func EmptyGaussianRevocationList(periodSize, nymNum int) witness.RevocationList {
	witness.InitBigInt = InitBigInt
	rlSize := EmptyGaussianRevocationListSize(periodSize, nymNum)
	return witness.EmptyRevocationList(rlSize)
}
