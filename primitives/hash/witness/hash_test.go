package witness

import (
	"math/big"
	"testing"

	_ "github.com/consensys/gnark-crypto/ecc/bls12-381/fr/poseidon2"
	gnarkhash "github.com/consensys/gnark-crypto/hash"
	"github.com/stretchr/testify/require"
)

func TestHashAbsorbsZeroAsFieldElement(t *testing.T) {
	hashFn := Hash(gnarkhash.POSEIDON2_BLS12_381)
	got, err := hashFn(big.NewInt(0), big.NewInt(7))
	require.NoError(t, err)

	hasher := gnarkhash.POSEIDON2_BLS12_381.New()
	zero := make([]byte, hasher.BlockSize())
	seven := make([]byte, hasher.BlockSize())
	big.NewInt(7).FillBytes(seven)

	_, err = hasher.Write(zero)
	require.NoError(t, err)
	_, err = hasher.Write(seven)
	require.NoError(t, err)

	expected := new(big.Int).SetBytes(hasher.Sum(nil))
	require.Equal(t, expected, got)
}

func TestHashRejectsInvalidFieldInputs(t *testing.T) {
	hashFn := Hash(gnarkhash.POSEIDON2_BLS12_381)

	_, err := hashFn(nil)
	require.ErrorContains(t, err, "is nil")

	_, err = hashFn(big.NewInt(-1))
	require.ErrorContains(t, err, "is negative")

	tooLarge := new(big.Int).Lsh(big.NewInt(1), uint(gnarkhash.POSEIDON2_BLS12_381.New().BlockSize()*8))
	_, err = hashFn(tooLarge)
	require.ErrorContains(t, err, "does not fit")
}
