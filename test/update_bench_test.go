package zkbantest

// import (
// 	"fmt"
// 	"testing"

// 	gnarkprecomputes "github.com/akakou/gnark-precomputes"
// 	zkban "github.com/akakou/zk-ban"
// 	zkbanc "github.com/akakou/zk-ban/circuit"
// 	"github.com/akakou/zk-ban/snark"
// 	zkbanw "github.com/akakou/zk-ban/witness"
// 	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
// 	"github.com/consensys/gnark/backend/groth16"
// 	"github.com/consensys/gnark/backend/witness"
// 	"github.com/consensys/gnark/frontend"
// 	// 	groth16_bls12381 "github.com/consensys/gnark/backend/groth16/bls12-381"
// 	// fr_bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
// )

// func BenchmarkUpdate(t *testing.B) {
// 	max := 10
// 	base := 6000
// 	session := 270

// 	for i := 1; i <= max; i++ {
// 		benchmarkUpdate(session, base*i, t)
// 	}
// }

// func benchmarkUpdate(a, b int, t *testing.B) {
// 	params := prepareParams()

// 	var proof groth16.Proof
// 	var pubWit witness.Witness
// 	var err error

// 	var rl zkbanw.RevocationList
// 	var tag = ""
// 	rl = zkbanw.EmptyConstantRevocationAddList(a, b)
// 	tag = fmt.Sprintf("%d,%d,%d,%v", a*b, a, b, "constant")

// 	_, _, updateCircuit := prepareCircuit(rl, true)
// 	tag += fmt.Sprintf("-%v", updateCircuit.ConstraintSystem.GetNbPublicVariables())

// 	t.Run("Prove ,"+tag, func(b *testing.B) {
// 		for range b.N {
// 			proof2, err := zkban.RequestUpdate(params.nextPeriod, params.signer(), rl, params.gpk, updateCircuit.Prover())
// 			panicIfErr(err)

// 			proof = proof2
// 			wit, err := frontend.NewWitness(updateAssign, snark.EcCurve.ScalarField())
// 			panicIfErr(err)

// 			pubWit, err = wit.Public()
// 			panicIfErr(err)
// 		}
// 	})

// 	vk, err := gnarkprecomputes.FromBLS12381GnarkKey(updateCircuit.VerifyKey, updateCircuit.Circuit.(*zkbanc.UpdateCircuit))
// 	panicIfErr(err)

// 	var prepare *bls12381.G1Jac
// 	t.Run("Verify-Precomputes,"+tag, func(b *testing.B) {
// 		for range b.N {
// 			prepare, err = vk.PreparePublicInputs(pubWit)
// 			panicIfErr(err)

// 		}
// 	})

// 	t.Run("Verify ,"+tag, func(b *testing.B) {
// 		for range b.N {
// 			params.gsk.IssueCredential(params.upk)

// 			err = vk.VerifyPrepared(proof, pubWit, prepare)
// 			panicIfErr(err)

// 		}
// 	})
// }
