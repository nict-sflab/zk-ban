package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/bench"
	"github.com/akakou/zk-ban/bench/storage"
	"github.com/akakou/zk-ban/test"
	"github.com/akakou/zk-ban/witness"
)

func main() {
	max := 10

	baseSessionNum := 60
	baseNymNum := 108_000

	alpha := 1
	beta := 50

	result := make(storage.Result, 0)
	result["env"] = make(map[string]int, 0)
	result["env"]["baseSessionNum"] = int(baseSessionNum)
	result["env"]["baseNymNum"] = int(baseNymNum)
	result["env"]["alpha"] = int(alpha)
	result["env"]["beta"] = int(beta)

	// increase nym
	for i := 1; i <= max; i++ {
		nymNum := baseNymNum * i * alpha
		rl, rlSize := test.EmptyUniformRevocationList(baseSessionNum, nymNum)
		benchUpdate(nymNum, bench.NYM_INCREASE_UNIFORM, rl, rlSize, result)
	}

	for i := 1; i <= max; i++ {
		nymNum := baseNymNum * i * alpha
		rl, rlSize := test.EmptyProportionalRevocationList(baseSessionNum, nymNum)
		benchUpdate(nymNum, bench.NYM_INCREASE_PROPORTIONAL, rl, rlSize, result)
	}

	for i := 1; i <= max; i++ {
		nymNum := baseNymNum * i * alpha
		rl, rlSize := test.EmptyGaussianRevocationList(baseSessionNum, nymNum)
		benchUpdate(nymNum, bench.NYM_INCREASE_GAUSSIAN, rl, rlSize, result)
	}

	// increase sessionNumber
	for i := 1; i <= max; i++ {
		sessionNum := baseSessionNum * i * beta
		rl, rlSize := test.EmptyUniformRevocationList(sessionNum, baseNymNum)
		benchUpdate(sessionNum, bench.SESS_INCREASE_UNIFORM, rl, rlSize, result)
	}

	for i := 1; i <= max; i++ {
		sessionNum := baseSessionNum * i * beta
		rl, rlSize := test.EmptyProportionalRevocationList(sessionNum, baseNymNum)
		benchUpdate(sessionNum, bench.SESS_INCREASE_PROPORTIONAL, rl, rlSize, result)
	}

	for i := 1; i <= max; i++ {
		sessionNum := baseSessionNum * i * beta
		rl, rlSize := test.EmptyGaussianRevocationList(sessionNum, baseNymNum)
		benchUpdate(sessionNum, bench.SESS_INCREASE_GAUSSIAN, rl, rlSize, result)
	}

	j, err := json.Marshal(result)
	test.PanicIfErr(err)
	fmt.Printf("%s", j)

	ioutil.WriteFile("update-storage.json", j, 0644)

}

func benchUpdate(v int, name string, rl witness.RevocationList, rlSize witness.RevocationListSize, result storage.Result) {
	parent := fmt.Sprintf("%v:%d", name, v)
	result[parent] = make(map[string]int, 0)

	params := test.PrepareParams()
	prover, verifier := test.PrepareUpdateKeyCached(rlSize, name)
	storage.StoreWritableSize("cs", parent, &prover.ConstraintSystem, result)
	storage.StoreWritableSize("pk", parent, &prover.ProveKey, result)
	storage.StoreWritableSize("vk", parent, verifier.VerifyKey, result)

	update, err := zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, prover)
	test.PanicIfErr(err)

	storage.StoreSize("upk", parent, update.PublicKey.Number.Bytes(), result)
	storage.StoreSize("ticket", parent, update.UpdateTicket.Number.Bytes(), result)
	storage.StoreWritableSize("update-proof", parent, &update.Proof, result)
}
