package witness

import (
	"math/big"

	_ "github.com/consensys/gnark-crypto/ecc/bls12-377/fr/poseidon2"
	gnarkhash "github.com/consensys/gnark-crypto/hash"
)

const HashType = gnarkhash.POSEIDON2_BLS12_377

func Hash(data ...*big.Int) (*big.Int, error) {
	hasher := HashType.New()
	for _, d := range data {
		_, err := hasher.Write(d.Bytes())
		if err != nil {
			return nil, err
		}
	}

	sum := hasher.Sum(nil)

	i := big.NewInt(0)
	i.SetBytes(sum)

	return i, nil
}
