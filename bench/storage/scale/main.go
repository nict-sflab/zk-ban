package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"

	gnarkserializable "github.com/akakou/gnark-serializable"
	"github.com/akakou/zk-ban/bench"
	"github.com/akakou/zk-ban/bench/storage"
	"github.com/akakou/zk-ban/dump"
	"github.com/akakou/zk-ban/test"
	"github.com/akakou/zk-ban/witness"
)

var baseSessionNum = 60
var baseNymNum = 108_000

var alpha = 1
var beta = 25
var max = 10

func main() {
	gnarkserializable.Unsafe = true

	result := make(storage.Result, 0)
	result["env"] = make(map[string]map[string]int)
	result["env"]["default"] = make(map[string]int)
	result["env"]["default"]["baseSessionNum"] = int(baseSessionNum)
	result["env"]["default"]["baseNymNum"] = int(baseNymNum)
	result["env"]["default"]["alpha"] = int(alpha)
	result["env"]["default"]["beta"] = int(beta)

	// increase nym
	BenchOneNymScalability(result, witness.MakeUniformRLSizeFromTotal)
	BenchOneNymScalability(result, witness.MakeProportionalRLSizeFromTotal)
	BenchOneNymScalability(result, test.EmptyGaussianRevocationListSize)

	BenchOneSessScalability(result, witness.MakeUniformRLSizeFromTotal)
	BenchOneSessScalability(result, witness.MakeProportionalRLSizeFromTotal)
	BenchOneSessScalability(result, test.EmptyGaussianRevocationListSize)

	j, err := json.Marshal(result)
	test.PanicIfErr(err)
	fmt.Printf("%s", j)

	os.WriteFile("update-storage.json", j, 0644)
}

func BenchOneNymScalability(result storage.Result, MakeRLSize func(int, int) witness.RevocationListSize) {
	for i := 1; i <= max; i++ {
		runtime.GC()

		nymNum := baseNymNum * i * alpha
		rlSize := witness.MakeUniformRLSizeFromTotal(baseSessionNum, nymNum)
		benchUpdate(nymNum, bench.NYM_INCREASE_UNIFORM, rlSize, result)
	}
}

func BenchOneSessScalability(result storage.Result, MakeRLSize func(int, int) witness.RevocationListSize) {
	for i := 1; i <= max; i++ {
		runtime.GC()

		sessionNum := baseSessionNum * i * beta
		rlSize := MakeRLSize(sessionNum, baseNymNum)
		benchUpdate(sessionNum, bench.SESS_INCREASE_UNIFORM, rlSize, result)
	}
}

func benchUpdate(v int, root string, rlSize witness.RevocationListSize, result storage.Result) {
	parent := fmt.Sprintf("%v--%d", root, v)
	vs := strconv.Itoa(v)

	test.PrepareUpdateKeyCached(rlSize, parent)

	cs, err := os.Stat(dump.FileName(parent, "update", dump.CircuitFileNameFormat))
	test.PanicIfErr(err)
	pk, err := os.Stat(dump.FileName(parent, "update", dump.ProverKeyFileNameFormat))
	test.PanicIfErr(err)
	vk, err := os.Stat(dump.FileName(parent, "update", dump.VerifierKeyFileNameFormat))
	test.PanicIfErr(err)

	storage.StoreSize(vs, "cs", root, int(cs.Size()), result)
	storage.StoreSize(vs, "pk", root, int(pk.Size()), result)
	storage.StoreSize(vs, "vk", root, int(vk.Size()), result)

	// params := test.PrepareParams()

	// update, err := zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, prover)
	// test.PanicIfErr(err)

	// storage.StoreBufSize(vs, "upk", parent, update.PublicKey.Number.Bytes(), result)
	// storage.StoreBufSize(vs, "ticket", parent, update.UpdateTicket.Number.Bytes(), result)
	// storage.StoreWritableSize(vs, "update-proof", parent, &update.Proof, result)

	res, err := json.Marshal(result)
	test.PanicIfErr(err)
	fmt.Printf("%s\n", res)
}
