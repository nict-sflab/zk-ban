package main

import (
	"fmt"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/bench/latency/utils/usefulbench"
	"github.com/akakou/zk-ban/bench/storage"
	"github.com/akakou/zk-ban/test"
	"github.com/akakou/zk-ban/witness"
)

func BenchmarkUpdate(b usefulbench.Benchmarker) {
	max := 10

	baseSessionNum := 60
	baseNymNum := 108_000

	alpha := 1
	beta := 50

	result := make(storage.Result, 0)
	result["env"]["baseSessionNum"] = int(baseSessionNum)
	result["env"]["baseNymNum"] = int(baseNymNum)
	result["env"]["alpha"] = int(alpha)
	result["env"]["beta"] = int(beta)

	// increase nym
	for i := 1; i <= max; i++ {
		nymNum := baseNymNum * i * alpha
		rl := test.EmptyUniformRevocationList(baseSessionNum, nymNum)
		benchUpdate(nymNum, "nym-increase-uniform", rl, result)
	}

	for i := 1; i <= max; i++ {
		nymNum := baseNymNum * i * alpha
		rl := test.EmptyProportionalRevocationList(baseSessionNum, nymNum)
		benchUpdate(nymNum, "nym-increase-proportional", rl, result)
	}

	for i := 1; i <= max; i++ {
		nymNum := baseNymNum * i * alpha
		rl := test.EmptyGaussianRevocationList(baseSessionNum, nymNum)
		benchUpdate(nymNum, "nym-increase-gaussian", rl, result)
	}

	// increase sessionNumber
	for i := 1; i <= max; i++ {
		sessionNum := baseSessionNum * i * beta
		rl := test.EmptyUniformRevocationList(sessionNum, baseNymNum)
		benchUpdate(sessionNum, "sess-increase-uniform", rl, result)
	}

	for i := 1; i <= max; i++ {
		sessionNum := baseSessionNum * i * beta
		rl := test.EmptyProportionalRevocationList(sessionNum, baseNymNum)
		benchUpdate(sessionNum, "sess-increase-proportional", rl, result)
	}

	for i := 1; i <= max; i++ {
		sessionNum := baseSessionNum * i * beta
		rl := test.EmptyGaussianRevocationList(sessionNum, baseNymNum)
		benchUpdate(sessionNum, "sess-increase-gaussian", rl, result)
	}
}

func benchUpdate(v int, name string, rl witness.RevocationList, result storage.Result) {
	parent := fmt.Sprintf("%v: %d", name, v)
	params := test.PrepareParams()
	_, _, updateCircuit := test.PrepareCircuit(rl, false)
	storage.StoreCircuitObjectSize("circuit", parent, updateCircuit, result)

	update, err := zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, updateCircuit.Prover())
	test.PanicIfErr(err)

	storage.StoreSize("upk", parent, update.PublicKey.Number.Bytes(), result)
	storage.StoreSize("ticket", parent, update.UpdateTicket.Number.Bytes(), result)
	storage.StoreWritableSize("update-proof", parent, &update.Proof, result)
}
