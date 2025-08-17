package main

import (
	"encoding/json"
	"fmt"
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

	baseNym := 108000
	baseSess := 60

	rlU, _ := test.EmptyUniformRevocationList(baseSess, baseNym)
	benchBasicUpdate(rlU, "uniform", &params, result)

	rlP, _ := test.EmptyProportionalRevocationList(baseSess, baseNym)
	benchBasicUpdate(rlP, "proportional", &params, result)

	rlG, _ := test.EmptyGaussianRevocationList(baseSess, baseNym)
	benchBasicUpdate(rlG, "gaussian", &params, result)

	j, err := json.Marshal(result)
	test.PanicIfErr(err)
	fmt.Printf("%s", j)
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
