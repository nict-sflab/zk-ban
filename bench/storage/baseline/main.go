package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/bench/storage"
	"github.com/akakou/zk-ban/test"
	"github.com/akakou/zk-ban/witness"
)

const KAPPA = 30
const LAMBDA = 30000

func main() {
	params := test.PrepareParams()

	result := make(storage.Result, 0)

	gsk, gpk, err := witness.RandomGroupKeyPair()
	test.PanicIfErr(err)

	now := int64(time.Now().UnixNano()) / int64(time.Hour*24)
	fmt.Printf("%v\n", now)

	fmt.Printf("gsk: %v\n", len(gsk.Bytes()))
	fmt.Printf("gpk: %v\n", len(gpk.Bytes()))

	storage.StoreBufSize("gsk", "baseline", "baseline", gsk.Bytes(), result)
	storage.StoreBufSize("gpk", "baseline", "baseline", gpk.Bytes(), result)

	joinCircuit, signCircuit, _ := test.PrepareCircuit(witness.EmptyRevocationList([]int{}))
	storage.StoreCircuitObjectSize("join", "baseline", joinCircuit, result)
	storage.StoreCircuitObjectSize("sign", "baseline", signCircuit, result)

	joinReq, usk, err := zkban.RequestJoin(params.Period, joinCircuit.Prover())
	test.PanicIfErr(err)

	storage.StoreBufSize("usk: ", "baseline", "baseline", usk.Number.Bytes(), result)
	storage.StoreBufSize("upk: ", "baseline", "baseline", joinReq.UserPublicKey.Number.Bytes(), result)
	storage.StoreWritableSize("join proof", "baseline", "baseline", joinReq.Proof, result)

	cred, err := params.GSK.IssueCredential(params.UPK)
	test.PanicIfErr(err)
	storage.StoreBufSize("cred", "baseline", "baseline", cred.Signature, result)

	signature, err := zkban.Sign(params.M, params.CNT, params.Signer(), params.GPK, signCircuit.Prover())
	test.PanicIfErr(err)

	storage.StoreBufSize("sigma", "baseline", "baseline", signature.Commit.Sigma.Bytes(), result)
	storage.StoreBufSize("nym", "baseline", "baseline", signature.Commit.Nym.Bytes(), result)
	storage.StoreWritableSize("sign proof", "baseline", "baseline", signature.Proof, result)

	benchBasicUpdateCircuit("baseline-circuit", result)

	j, err := json.Marshal(result)

	test.PanicIfErr(err)
	fmt.Printf("%s", j)

	os.WriteFile("baseline-storage.json", j, 0o644)
}

func benchBasicUpdateCircuit(root string, result storage.Result) {
	kappas := []int{KAPPA / 2, KAPPA, KAPPA * 2}
	lambdas := []int{LAMBDA / 2, LAMBDA, LAMBDA * 2}

	witness.InitBigInt = witness.ZeroInitBigInt
	rlMakers := []func(int, int) witness.RevocationListSize{
		witness.MakeUniformRLSizeFromTotal,
	}

	rlMakerTags := []string{
		"uniform",
	}
	count := 0

	for _, kappa := range kappas {
		// for _, lambda := range lambdas {
		for i, rlMaker := range rlMakers {
			runtime.GC()

			count++
			tag := rlMakerTags[i]
			rlSize := rlMaker(kappa, 60)
			name := fmt.Sprintf("%d-%d-%s", kappa, 1, tag)

			storage.BenchUpdate(rlSize, root, name, name, result)
			fmt.Printf("%d/%d is done...", count, len(kappas)*len(lambdas)*len(rlMakers))
			// }
		}
	}
}
