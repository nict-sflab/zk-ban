package witness

import (
	"math/big"

	"github.com/akakou/zk-ban/snark"
)

type UserPublicKey struct {
	*big.Int
}

type OneTimeTicket struct {
	*big.Int
}

func (usk *UserSecretKey) PseudoRandom(period int64, name int64) (*big.Int, error) {
	hash, err := snark.CommitHash(big.NewInt(name), big.NewInt(period), &usk.Int)

	if err != nil {
		return nil, err
	}

	return hash, nil
}

func (usk *UserSecretKey) PublicKey(period int64) (*UserPublicKey, error) {
	upk, err := usk.PseudoRandom(period, PUBLIC_KEY)
	if err != nil {
		return nil, err
	}

	return &UserPublicKey{upk}, nil
}

func (usk *UserSecretKey) OneTimeTicket(period int64) (*OneTimeTicket, error) {
	ticket, err := usk.PseudoRandom(period, ONE_TIME_TICKET)

	if err != nil {
		return nil, err
	}

	return &OneTimeTicket{ticket}, nil
}
