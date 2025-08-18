package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	gnarkserializable "github.com/akakou/gnark-serializable"
	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/bench/storage"
	"github.com/akakou/zk-ban/test"
	"github.com/akakou/zk-ban/witness"
)

func main() {
	gnarkserializable.Unsafe = true
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

	tag := witness.SessionTag(2, now)
	fmt.Printf("tag: %v\n", len(tag.Bytes()))

	signature, err := zkban.Sign(params.M, params.CNT, params.Signer(), params.GPK, signCircuit.Prover())
	test.PanicIfErr(err)

	storage.StoreBufSize("sigma", "baseline", "baseline", signature.Commit.Sigma.Bytes(), result)
	storage.StoreBufSize("nym", "baseline", "baseline", signature.Commit.Nym.Bytes(), result)
	storage.StoreWritableSize("sign proof", "baseline", "baseline", &signature.Proof, result)

	// baseNym := 108000
	// baseSess := 60

	// rlU, _ := test.EmptyUniformRevocationList(baseSess, baseNym)
	// benchBasicUpdate(rlU, "uniform", &params, result)

	// rlP, _ := test.EmptyProportionalRevocationList(baseSess, baseNym)
	// benchBasicUpdate(rlP, "proportional", &params, result)

	// rlG, _ := test.EmptyGaussianRevocationList(baseSess, baseNym)
	// benchBasicUpdate(rlG, "gaussian", &params, result)

	benchBasicUpdateCircuit("baseline-circuit", result)

	j, err := json.Marshal(result)

	test.PanicIfErr(err)
	fmt.Printf("%s", j)

	os.WriteFile("baseline-storage.json", j, 0o644)
}

func benchBasicUpdate(rl witness.RevocationList, name string, params *test.TestParams, result storage.Result) {
	_, _, updateCircuit := test.PrepareCircuit(rl)
	n := name + "-update"
	storage.StoreCircuitObjectSize(n, "baseline", updateCircuit, result)

	update, err := zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, updateCircuit.Prover())
	test.PanicIfErr(err)

	storage.StoreBufSize("upk", n, name, update.PublicKey.Number.Bytes(), result)
	storage.StoreBufSize(name+"ticket", n, name, update.UpdateTicket.Number.Bytes(), result)
	storage.StoreWritableSize("-update-proof", n, name, &update.Proof, result)
}

func benchBasicUpdateCircuit(root string, result storage.Result) {
	kappas := []int{15, 30, 60, 120}
	lambdas := []int{27_000, 54_000, 108_000, 216_000}

	// rlMakers := []func(int, int) (witness.RevocationList, witness.RevocationListSize){
	// 	test.EmptyUniformRevocationList,
	// 	test.EmptyProportionalRevocationList,
	// 	test.EmptyGaussianRevocationList,
	// }

	witness.InitBigInt = witness.ZeroInitBigInt
	rlMakers := []func(int, int) witness.RevocationListSize{
		witness.MakeUniformRLSizeFromTotal,
		witness.MakeProportionalRLSizeFromTotal,
		test.EmptyGaussianRevocationListSize,
	}

	rlMakerTags := []string{"uniform", "proportionl", "gaussian"}
	count := 0

	for _, kappa := range kappas {
		for _, lambda := range lambdas {
			for i, rlMaker := range rlMakers {
				runtime.GC()

				count++
				tag := rlMakerTags[i]
				rlSize := rlMaker(kappa, lambda)
				name := fmt.Sprintf("%d-%d-%s", kappa, lambda, tag)
				pk, vk := test.PrepareUpdateKeyCached(rlSize, name)

				storage.StoreWritableSize(name, "pk", root, pk.ProveKey, result)
				storage.StoreWritableSize(name, "cs", root, &pk.ConstraintSystem, result)
				storage.StoreWritableSize(name, "vk", root, vk.VerifyKey, result)

				j, err := json.Marshal(result)
				test.PanicIfErr(err)
				fmt.Printf("result %s\n", j)
				fmt.Printf("%d/%d is done...", count, 4*4*3)

			}
		}
	}
}
