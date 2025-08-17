package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	gnarkserializable "github.com/akakou/gnark-serializable"
	"github.com/akakou/zk-ban/bench"
	"github.com/akakou/zk-ban/bench/storage"
	"github.com/akakou/zk-ban/test"
	"github.com/akakou/zk-ban/witness"
)

func main() {
	gnarkserializable.Unsafe = true
	max := 10

	baseSessionNum := 60
	baseNymNum := 108_000

	alpha := 1
	beta := 25

	result := make(storage.Result, 0)
	result["env"] = make(map[string]map[string]int)
	result["env"]["default"] = make(map[string]int)
	result["env"]["default"]["baseSessionNum"] = int(baseSessionNum)
	result["env"]["default"]["baseNymNum"] = int(baseNymNum)
	result["env"]["default"]["alpha"] = int(alpha)
	result["env"]["default"]["beta"] = int(beta)

	/*

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
	*/

	// increase sessionNumber
	// for i := 1; i <= max; i++ {
	// 	sessionNum := baseSessionNum * i * beta
	// 	rl, rlSize := test.EmptyUniformRevocationList(sessionNum, baseNymNum)
	// 	benchUpdate(sessionNum, bench.SESS_INCREASE_UNIFORM, rl, rlSize, result)
	// }

	// for i := 1; i <= max; i++ {
	// 	sessionNum := baseSessionNum * i * beta
	// 	rl, rlSize := test.EmptyProportionalRevocationList(sessionNum, baseNymNum)
	// 	benchUpdate(sessionNum, bench.SESS_INCREASE_PROPORTIONAL, rl, rlSize, result)
	// }

	for i := 1; i <= max; i++ {
		sessionNum := baseSessionNum * i * beta
		rl, rlSize := test.EmptyGaussianRevocationList(sessionNum, baseNymNum)
		benchUpdate(sessionNum, bench.SESS_INCREASE_GAUSSIAN, rl, rlSize, result)
	}

	j, err := json.Marshal(result)
	test.PanicIfErr(err)
	fmt.Printf("%s", j)

	os.WriteFile("update-storage.json", j, 0644)
}

func benchUpdate(v int, root string, rl witness.RevocationList, rlSize witness.RevocationListSize, result storage.Result) {
	parent := fmt.Sprintf("%v:%d", root, v)
	vs := strconv.Itoa(v)

	prover, verifier := test.PrepareUpdateKeyCached(rlSize, parent)
	storage.StoreWritableSize(vs, "cs", root, &prover.ConstraintSystem, result)
	storage.StoreWritableSize(vs, "pk", root, &prover.ProveKey, result)
	storage.StoreWritableSize(vs, "vk", root, verifier.VerifyKey, result)

	// params := test.PrepareParams()
	// update, err := zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, prover)
	// test.PanicIfErr(err)

	// storage.StoreBufSize(vs, "upk", root, update.PublicKey.Number.Bytes(), result)
	// storage.StoreBufSize(vs, "ticket", root, update.UpdateTicket.Number.Bytes(), result)
	// storage.StoreWritableSize(vs, "update-proof", root, &update.Proof, result)

	res, err := json.Marshal(result)
	test.PanicIfErr(err)
	fmt.Printf("%s\n", res)
}
