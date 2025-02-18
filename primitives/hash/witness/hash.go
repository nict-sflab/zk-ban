package witness

import (
	"math/big"

	gnarkhash "github.com/consensys/gnark-crypto/hash"
	"github.com/iden3/go-iden3-crypto/poseidon"
)

func PoseidonHash(data ...*big.Int) (*big.Int, error) {
	hash, err := poseidon.Hash(data)

	if err != nil {
		return nil, err
	}

	return hash, nil

}

func MimcHash(h gnarkhash.Hash) func(data ...*big.Int) (*big.Int, error) {
	return func(data ...*big.Int) (*big.Int, error) {
		hasher := h.New()
		for _, d := range data {
			_, err := hasher.Write(d.Bytes())
			if err != nil {
				return nil, err
			}
		}

		h := hasher.Sum(nil)

		i := big.NewInt(0)
		i.SetBytes(h)

		return i, nil
	}
}
