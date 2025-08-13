package test

import (
	"fmt"
	"testing"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/precomputes"
	"github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/signature/eddsa"
	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	x := 10
	y := 10

	rlSize := witness.MakeUniformRLSize(x, y)
	rl := witness.EmptyRevocationList(rlSize)
	_, _, updateCircuit := prepareCircuit(rl, true)

	index := updateCircuit.ConstraintSystem.GetNbPublicVariables() - x*y - x - 1 - 2

	wit, err := frontend.NewWitness(
		&precomputes.UpdateCircuit{
			circuit.UpdateCircuit{
				RevocationList: circuit.NewRevocationListAssigned(rl),
				UserSecretKey:  1,
				Credential: eddsa.Signature{
					R: twistededwards.Point{
						X: 1,
						Y: 1,
					},
					S: 1,
				},
				NextInfo: circuit.PublicKeyAuthInfo{
					UserPublicKey: 1,
					Period:        1,
				},
				CurrentInfo: circuit.PublicKeyAuthInfo{
					UserPublicKey: 1,
					Period:        1,
				},
				GroupPublicKey: eddsa.PublicKey{
					A: twistededwards.Point{
						X: 0,
						Y: 0,
					},
				},
			},
		}, ecc.BLS12_381.ScalarField())

	assert.NoError(t, err)

	pubWit, err := wit.Public()
	assert.NoError(t, err)
	fmt.Printf("%v", pubWit.Vector())

	assert.Equal(t, precomputes.UpdatePreparableIndex, index)

}
