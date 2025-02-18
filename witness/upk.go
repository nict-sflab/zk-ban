package witness

import (
	"math/big"

	"github.com/akakou/zk-ban/snark"
)

type UserPublicKey struct {
	Number *big.Int
}

func (usk *UserSecretKey) PublicKey(period *big.Int) (*UserPublicKey, error) {
	hash, err := snark.CommitHash(period, usk.Number)

	if err != nil {
		return nil, err
	}

	return &UserPublicKey{Number: hash}, nil
}
