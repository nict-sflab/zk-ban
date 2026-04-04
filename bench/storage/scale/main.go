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
	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/dump"
	"github.com/akakou/zk-ban/test"
	"github.com/akakou/zk-ban/witness"
)

var basePeriodNum = 30
var baseNymNum = 30000
var baseMaxSess = 10

var alpha = 1
var beta = 10
var gamma = 10
var max = 10

func main() {
	gnarkserializable.Unsafe = true

	result := make(storage.Result, 0)
	result["env"] = make(map[string]map[string]int)
	result["env"]["default"] = make(map[string]int)
	result["env"]["default"]["basePeriodNum"] = int(basePeriodNum)
	result["env"]["default"]["baseNymNum"] = int(baseNymNum)
	result["env"]["default"]["alpha"] = int(alpha)
	result["env"]["default"]["beta"] = int(beta)

	// increase nym
	BenchOneNymScalability(result, bench.NYM_INCREASE_UNIFORM, witness.MakeUniformRLSizeFromTotal)
	BenchOneNymScalability(result, bench.NYM_INCREASE_PROPORTIONAL, witness.MakeProportionalRLSizeFromTotal)
	BenchOneNymScalability(result, bench.NYM_INCREASE_GAUSSIAN, test.EmptyGaussianRevocationListSize)

	BenchOnePeriodScalability(result, bench.PERIOD_INCREASE_UNIFORM, witness.MakeUniformRLSizeFromTotal)
	BenchOnePeriodScalability(result, bench.PERIOD_INCREASE_PROPORTIONAL, witness.MakeProportionalRLSizeFromTotal)
	BenchOnePeriodScalability(result, bench.PERIOD_INCREASE_GAUSSIAN, test.EmptyGaussianRevocationListSize)

	BenchOnePeriodScalability(result, bench.PERIOD_INCREASE_UNIFORM, witness.MakeUniformRLSizeFromTotal)
	BenchOnePeriodScalability(result, bench.PERIOD_INCREASE_PROPORTIONAL, witness.MakeProportionalRLSizeFromTotal)
	BenchOnePeriodScalability(result, bench.PERIOD_INCREASE_GAUSSIAN, test.EmptyGaussianRevocationListSize)

	j, err := json.Marshal(result)
	test.PanicIfErr(err)
	fmt.Printf("%s", j)

	os.WriteFile("update-storage.json", j, 0644)
}

func BenchOneNymScalability(result storage.Result, name string, MakeRLSize func(int, int) witness.RevocationListSize) {
	for i := 1; i <= max; i++ {
		runtime.GC()

		nymNum := baseNymNum * i * alpha
		rlSize := witness.MakeUniformRLSizeFromTotal(basePeriodNum, nymNum)
		benchUpdate(nymNum, name, rlSize, result)
	}
}

func BenchOnePeriodScalability(result storage.Result, name string, MakeRLSize func(int, int) witness.RevocationListSize) {
	for i := 1; i <= max; i++ {
		runtime.GC()

		periodNum := basePeriodNum * i * beta
		rlSize := MakeRLSize(periodNum, baseNymNum)
		benchUpdate(periodNum, name, rlSize, result)
	}
}

func BenchOneMaxSessScalability(result storage.Result, name string, MakeRLSize func(int, int) witness.RevocationListSize) {
	for i := 1; i <= max; i++ {
		runtime.GC()

		maxSess := baseMaxSess * i * gamma
		circuit.MaxSession = maxSess

		rlSize := MakeRLSize(basePeriodNum, baseNymNum)
		benchUpdate(maxSess, name, rlSize, result)
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
