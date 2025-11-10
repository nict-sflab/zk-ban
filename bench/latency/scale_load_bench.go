package latency

import (
	"bytes"
	"fmt"
	"os"
	"runtime"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/bench"
	"github.com/akakou/zk-ban/bench/latency/utils/usefulbench"
	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	"github.com/akakou/zk-ban/test"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	// 	groth16_bls12381 "github.com/consensys/gnark/backend/groth16/bls12-381"
	// fr_bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

const BASELINE_NAME = "load_baseline"

func BenchmarkLoad(b usefulbench.Benchmarker) {
	BenchmarkBaselineLoad(b)

	BenchmarkNymLoad(b)
	BenchSessLoad(b)
}

func BenchmarkBaselineLoad(b usefulbench.Benchmarker) {
	unifomrRL := test.EmptyUniformRevocationList(baseSessionNum, baseNymNum)
	proportionalRL := test.EmptyProportionalRevocationList(baseSessionNum, baseNymNum)
	gaussiunRL := test.EmptyGaussianRevocationList(baseSessionNum, baseNymNum)

	benchmarkLoad(0, unifomrRL, BASELINE_NAME, b)
	benchmarkLoad(1, proportionalRL, BASELINE_NAME, b)
	benchmarkLoad(2, gaussiunRL, BASELINE_NAME, b)
}

func BenchmarkNymLoad(b usefulbench.Benchmarker) {
	logIfUsefulBenchInScaleBench(b)

	// increase nym
	for i := 1; i <= max; i++ {
		runtime.GC()

		nymNum := baseNymNum * i * alpha
		rl := test.EmptyUniformRevocationList(baseSessionNum, nymNum)
		benchmarkLoad(nymNum, rl, bench.NYM_INCREASE_UNIFORM, b)
	}

	for i := 1; i <= max; i++ {
		runtime.GC()

		nymNum := baseNymNum * i * alpha
		rl := test.EmptyProportionalRevocationList(baseSessionNum, nymNum)
		benchmarkLoad(nymNum, rl, bench.NYM_INCREASE_PROPORTIONAL, b)
	}

	for i := 1; i <= max; i++ {
		runtime.GC()

		nymNum := baseNymNum * i * alpha
		rl := test.EmptyGaussianRevocationList(baseSessionNum, nymNum)
		benchmarkLoad(nymNum, rl, bench.NYM_INCREASE_GAUSSIAN, b)
	}
}

func BenchSessLoad(b usefulbench.Benchmarker) {
	// increase sessionNumber
	BenchSessUniformLoad(b)
	BenchSessProportionalLoad(b)
	BenchSessGaussLoad(b)

}

func BenchSessUniformLoad(b usefulbench.Benchmarker) {
	logIfUsefulBenchInScaleBench(b)

	for i := 1; i <= max; i++ {
		runtime.GC()

		sessionNum := baseSessionNum * i * beta
		rl := test.EmptyUniformRevocationList(sessionNum, baseNymNum)
		benchmarkLoad(sessionNum, rl, bench.SESS_INCREASE_UNIFORM, b)
	}

}

func BenchSessProportionalLoad(b usefulbench.Benchmarker) {
	logIfUsefulBenchInScaleBench(b)

	for i := 1; i <= max; i++ {
		runtime.GC()

		sessionNum := baseSessionNum * i * beta
		rl := test.EmptyProportionalRevocationList(sessionNum, baseNymNum)
		benchmarkLoad(sessionNum, rl, bench.SESS_INCREASE_PROPORTIONAL, b)
	}
}

func BenchSessGaussLoad(b usefulbench.Benchmarker) {
	logIfUsefulBenchInScaleBench(b)

	for i := 1; i <= max; i++ {
		runtime.GC()

		sessionNum := baseSessionNum * i * beta
		rl := test.EmptyGaussianRevocationList(sessionNum, baseNymNum)
		benchmarkLoad(sessionNum, rl, bench.SESS_INCREASE_GAUSSIAN, b)
	}
}

var benchmarkLoad = benchmarkLoadWithWrite

func benchmarkLoadWithWrite(param int, rl zkbanw.RevocationList, name string, b usefulbench.Benchmarker) {
	parent := fmt.Sprintf("%s--%v", name, param)

	samplePk, sampleVK := test.PrepareUpdateKeyCached(rl.Sizes(), parent)

	var buf bytes.Buffer
	buf.Reset()
	_, err := samplePk.ProveKey.ProvingKey.WriteRawTo(&buf)
	test.PanicIfErr(err)

	pkFile := fmt.Sprintf("%s/%s.prover.bin", test.TestKeyPath, parent)
	err = os.WriteFile(pkFile, buf.Bytes(), 0o600)
	test.PanicIfErr(err)

	runtime.GC()

	var (
		pk    groth16.ProvingKey
		pkBin []byte
	)
	b.Run("io-read-pk-"+parent, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			pkBin, err = os.ReadFile(pkFile)
			test.PanicIfErr(err)
		}
	})
	b.Run("decode-pk-"+parent, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			pk = groth16.NewProvingKey(ecc.BLS12_381)
			_, err = pk.UnsafeReadFrom(bytes.NewBuffer(pkBin))
			test.PanicIfErr(err)
		}
	})

	buf.Reset()
	_, err = sampleVK.VerifyKey.VerifyingKey.WriteRawTo(&buf)
	test.PanicIfErr(err)

	vkFile := fmt.Sprintf("%s/%s.verifier.bin", test.TestKeyPath, parent)
	err = os.WriteFile(vkFile, buf.Bytes(), 0o600)
	test.PanicIfErr(err)

	var (
		vk    groth16.VerifyingKey
		vkBin []byte
	)
	b.Run("io-read-vk-"+parent, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			vkBin, err = os.ReadFile(vkFile)
			test.PanicIfErr(err)
		}
	})
	b.Run("decode-vk-"+parent, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			vk = groth16.NewVerifyingKey(ecc.BLS12_381)
			_, err = vk.UnsafeReadFrom(bytes.NewBuffer(vkBin))
			test.PanicIfErr(err)
		}
	})

	buf.Reset()
	_, err = samplePk.ConstraintSystem.WriteTo(&buf)
	test.PanicIfErr(err)

	csFile := fmt.Sprintf("%s/%s.cs.bin", test.TestKeyPath, parent)
	err = os.WriteFile(csFile, buf.Bytes(), 0o600)
	test.PanicIfErr(err)

	var (
		csBin []byte
		cs    constraint.ConstraintSystem
	)
	b.Run("io-read-cs-"+parent, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			csBin, err = os.ReadFile(csFile)
			test.PanicIfErr(err)
		}
	})
	b.Run("decode-cs-"+parent, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			cs = groth16.NewCS(ecc.BLS12_381)
			_, err = cs.ReadFrom(bytes.NewBuffer(csBin))
			test.PanicIfErr(err)
		}
	})

	checkISValidKeys(
		&snark.SnarkParams{
			ConstraintSystem: cs,
			ProveKey:         pk,
			VerifyKey:        vk,
		}, &rl)
}

func checkISValidKeys(updateCircuit *snark.SnarkParams, rl *zkbanw.RevocationList) {
	params := test.PrepareParams()

	updateCircuit.Circuit = &circuit.UpdateCircuit{RevocationList: circuit.NewRevocationListAssigned(*rl)}

	updateReq, err := zkban.RequestUpdate(params.NextPeriod, params.Signer(), *rl, params.GPK, updateCircuit.Prover())
	test.PanicIfErr(err)

	err = updateReq.Verify(params.NextPeriod, params.Period, *rl, params.GPK, updateCircuit.VerifyKey)
	test.PanicIfErr(err)

}

// func benchmarkLoad(param int, rl zkbanw.RevocationList, name string, b usefulbench.Benchmarker) {
// 	parent := fmt.Sprintf("%s--%v", name, param)
// 	samplePk, _ := test.PrepareUpdateKeyCached(rl.Sizes(), parent)

// 	var buf bytes.Buffer
// 	_, err := samplePk.ProveKey.ProvingKey.WriteRawTo(&buf)
// 	test.PanicIfErr(err)

// 	proverKeyFileName := fmt.Sprintf("%s.prover.bin", parent)

// 	err = os.WriteFile(proverKeyFileName, buf.Bytes(), 0x600)
// 	test.PanicIfErr(err)

// 	var pk groth16.ProvingKey
// 	var bin []byte
// 	b.Run("read-pk-"+parent, func(b usefulbench.Benchmarker) {
// 		for b.Loop() {
// 			bin, err = os.ReadFile(proverKeyFileName)
// 			test.PanicIfErr(err)
// 		}
// 	})

// 	b.Run("read-pk-"+parent, func(b usefulbench.Benchmarker) {
// 		for b.Loop() {
// 			pk = groth16.NewProvingKey(ecc.BLS12_381)
// 			_, err = pk.UnsafeReadFrom(bytes.NewBuffer(bin))
// 			test.PanicIfErr(err)
// 		}
// 	})
// }

// func benchmarkLoadWithDump(param int, rl zkbanw.RevocationList, name string, b usefulbench.Benchmarker) {
// 	parent := fmt.Sprintf("%s--%v", name, param)
// 	samplePk, sampleVK := test.PrepareUpdateKeyCached(rl.Sizes(), parent)

// 	var buf bytes.Buffer
// 	buf.Reset()
// 	_, err := samplePk.ProveKey.ProvingKey.WriteRawTo(&buf)
// 	test.PanicIfErr(err)
// 	pkRaw := fmt.Sprintf("%s.prover.raw.bin", parent)
// 	test.PanicIfErr(os.WriteFile(pkRaw, buf.Bytes(), 0o600))

// 	var pk groth16.ProvingKey
// 	var pkRawBin []byte
// 	b.Run("io-read-pk-raw-"+parent, func(b usefulbench.Benchmarker) {
// 		for b.Loop() {
// 			pkRawBin, err = os.ReadFile(pkRaw)
// 			test.PanicIfErr(err)
// 		}
// 	})
// 	b.Run("decode-pk-raw-"+parent, func(b usefulbench.Benchmarker) {
// 		for b.Loop() {
// 			pk = groth16.NewProvingKey(ecc.BLS12_381)
// 			_, err = pk.UnsafeReadFrom(bytes.NewReader(pkRawBin))
// 			test.PanicIfErr(err)
// 		}
// 	})

// 	buf.Reset()
// 	if dumper, ok := any(&samplePk.ProveKey.ProvingKey).(gnarkio.BinaryDumper); ok {
// 		test.PanicIfErr(dumper.WriteDump(&buf))
// 		pkDump := fmt.Sprintf("%s.prover.dump.bin", parent)
// 		test.PanicIfErr(os.WriteFile(pkDump, buf.Bytes(), 0o600))

// 		var pkDumpBin []byte
// 		b.Run("io-read-pk-dump-"+parent, func(b usefulbench.Benchmarker) {
// 			for b.Loop() {
// 				pkDumpBin, err = os.ReadFile(pkDump)
// 				test.PanicIfErr(err)
// 			}
// 		})
// 		b.Run("decode-pk-dump-"+parent, func(b usefulbench.Benchmarker) {
// 			for b.Loop() {
// 				pk = groth16.NewProvingKey(ecc.BLS12_381)
// 				d, ok := pk.(gnarkio.BinaryDumper)
// 				if !ok {
// 					test.PanicIfErr(fmt.Errorf("pk not BinaryDumper"))
// 				}
// 				test.PanicIfErr(d.ReadDump(bytes.NewReader(pkDumpBin)))
// 			}
// 		})
// 	}

// 	buf.Reset()
// 	_, err = sampleVK.VerifyKey.VerifyingKey.WriteRawTo(&buf)
// 	test.PanicIfErr(err)
// 	vkRaw := fmt.Sprintf("%s.verifier.raw.bin", parent)
// 	test.PanicIfErr(os.WriteFile(vkRaw, buf.Bytes(), 0o600))

// 	var vk groth16.VerifyingKey
// 	var vkRawBin []byte
// 	b.Run("io-read-vk-raw-"+parent, func(b usefulbench.Benchmarker) {
// 		for b.Loop() {
// 			vkRawBin, err = os.ReadFile(vkRaw)
// 			test.PanicIfErr(err)
// 		}
// 	})
// 	b.Run("decode-vk-raw-"+parent, func(b usefulbench.Benchmarker) {
// 		for b.Loop() {
// 			vk = groth16.NewVerifyingKey(ecc.BLS12_381)
// 			_, err = vk.UnsafeReadFrom(bytes.NewReader(vkRawBin))
// 			test.PanicIfErr(err)
// 		}
// 	})

// 	buf.Reset()
// 	if dumper, ok := any(&sampleVK.VerifyKey.VerifyingKey).(gnarkio.BinaryDumper); ok {
// 		test.PanicIfErr(dumper.WriteDump(&buf))
// 		vkDump := fmt.Sprintf("%s.verifier.dump.bin", parent)
// 		test.PanicIfErr(os.WriteFile(vkDump, buf.Bytes(), 0o600))

// 		var vkDumpBin []byte
// 		b.Run("io-read-vk-dump-"+parent, func(b usefulbench.Benchmarker) {
// 			for b.Loop() {
// 				vkDumpBin, err = os.ReadFile(vkDump)
// 				test.PanicIfErr(err)
// 			}
// 		})
// 		b.Run("decode-vk-dump-"+parent, func(b usefulbench.Benchmarker) {
// 			for b.Loop() {
// 				vk = groth16.NewVerifyingKey(ecc.BLS12_381)
// 				d, ok := vk.(gnarkio.BinaryDumper)
// 				if !ok {
// 					test.PanicIfErr(fmt.Errorf("vk not BinaryDumper"))
// 				}
// 				test.PanicIfErr(d.ReadDump(bytes.NewReader(vkDumpBin)))
// 			}
// 		})
// 	}

// }

// func benchmarkLoadWithDump2(param int, rl zkbanw.RevocationList, name string, b usefulbench.Benchmarker) {
// 	parent := fmt.Sprintf("%s--%v", name, param)
// 	samplePk, _ := test.PrepareUpdateKeyCached(rl.Sizes(), parent)

// 	var buf bytes.Buffer
// 	buf.Reset()
// 	_, err := samplePk.ProveKey.ProvingKey.WriteRawTo(&buf)
// 	test.PanicIfErr(err)
// 	pkRaw := fmt.Sprintf("%s.prover.raw.bin", parent)
// 	test.PanicIfErr(os.WriteFile(pkRaw, buf.Bytes(), 0o600))

// 	var pk groth16.ProvingKey
// 	f, err := os.Open(pkRaw)
// 	test.PanicIfErr(err)

// 	b.Run("decode-pk-raw-"+parent, func(b usefulbench.Benchmarker) {
// 		for b.Loop() {
// 			pk = groth16.NewProvingKey(ecc.BLS12_381)
// 			_, err = pk.UnsafeReadFrom(f)
// 			test.PanicIfErr(err)
// 		}
// 	})
// }
