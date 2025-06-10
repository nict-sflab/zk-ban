package witness

import (
	"math/big"

	hash "github.com/akakou/zk-ban/primitives/hash/witness"
)

type UserPublicKey struct {
	Number *big.Int
}

func (usk *UserSecretKey) PublicKey(period *big.Int) (*UserPublicKey, error) {
	hash, err := hash.Hash(period, usk.Number)

	if err != nil {
		return nil, err
	}

	return &UserPublicKey{Number: hash}, nil
}
