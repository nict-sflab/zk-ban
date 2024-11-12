package zkban

import (
	"testing"

	"github.com/consensys/gnark/test"
)

func TestAuth(t *testing.T) {
	assert := test.NewAssert(t)
	authCircuit := AuthCircuit{}

	assert.ProverSucceeded(&authCircuit, &AuthCircuit{SK: 1})
}
