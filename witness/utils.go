package witness

import (
	"math/big"

	"github.com/iden3/go-iden3-crypto/poseidon"
)

type hashT = func(data ...*big.Int) ([]byte, error)

var hash = poseidonHash

func poseidonHash(data ...*big.Int) (*big.Int, error) {
	hash, err := poseidon.Hash(data)

	if err != nil {
		return nil, err
	}

	return hash, nil

}
