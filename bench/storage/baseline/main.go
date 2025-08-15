package main

import (
	"encoding/json"
	"fmt"
	"time"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/bench/storage"
	"github.com/akakou/zk-ban/test"
	"github.com/akakou/zk-ban/witness"
)

func main() {
	params := test.PrepareParams()

	result := make(storage.Result, 0)

	gsk, gpk, err := witness.RandomGroupKeyPair()
	test.PanicIfErr(err)

	now := int64(time.Now().UnixNano()) / int64(time.Hour*24)
	fmt.Printf("%v\n", now)

	result["baseline"] = make(map[string]int)

	fmt.Printf("gsk: %v\n", len(gsk.Bytes()))
	fmt.Printf("gpk: %v\n", len(gpk.Bytes()))

	storage.StoreSize("gsk", "baseline", gsk.Bytes(), result)
	storage.StoreSize("gpk", "baseline", gpk.Bytes(), result)

	joinCircuit, signCircuit, _ := test.PrepareCircuit(witness.EmptyRevocationList([]int{}), false)
	storage.StoreCircuitObjectSize("join", "baseline", joinCircuit, result)
	storage.StoreCircuitObjectSize("sign", "baseline", signCircuit, result)

	joinReq, usk, err := zkban.RequestJoin(params.Period, joinCircuit.Prover())
	test.PanicIfErr(err)

	storage.StoreSize("usk: ", "baseline", usk.Number.Bytes(), result)
	storage.StoreSize("upk: ", "baseline", joinReq.UserPublicKey.Number.Bytes(), result)
	storage.StoreWritableSize("join proof", "baseline", joinReq.Proof, result)

	cred, err := params.GSK.IssueCredential(params.UPK)
	test.PanicIfErr(err)
	storage.StoreSize("cred", "baseline", cred.Signature, result)

	tag := witness.SessionTag(2, now)
	fmt.Printf("tag: %v\n", len(tag.Bytes()))

	signature, err := zkban.Sign(params.M, params.CNT, params.Signer(), params.GPK, signCircuit.Prover())
	test.PanicIfErr(err)

	storage.StoreSize("sigma", "baseline", signature.Commit.Sigma.Bytes(), result)
	storage.StoreSize("nym", "baseline", signature.Commit.Nym.Bytes(), result)
	storage.StoreWritableSize("sign proof", "baseline", &signature.Proof, result)

	baseNym := 108000
	baseSess := 60

	rlU := test.EmptyUniformRevocationList(baseSess, baseNym)
	benchBasicUpdate(rlU, "uniform", &params, result)

	rlP := test.EmptyProportionalRevocationList(baseSess, baseNym)
	benchBasicUpdate(rlP, "proportional", &params, result)

	rlG := test.EmptyGaussianRevocationList(baseSess, baseNym)
	benchBasicUpdate(rlG, "gaussian", &params, result)

	j, err := json.Marshal(result)
	test.PanicIfErr(err)
	fmt.Printf("%s", j)
}

func benchBasicUpdate(rl witness.RevocationList, name string, params *test.TestParams, result storage.Result) {
	result[name] = make(map[string]int)
	_, _, updateCircuit := test.PrepareCircuit(rl, false)
	storage.StoreCircuitObjectSize(name+"-update", "baseline", updateCircuit, result)

	update, err := zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, updateCircuit.Prover())
	test.PanicIfErr(err)

	storage.StoreSize("upk", name, update.PublicKey.Number.Bytes(), result)
	storage.StoreSize(name+"ticket", name, update.UpdateTicket.Number.Bytes(), result)
	storage.StoreWritableSize("-update-proof", name, &update.Proof, result)
}
