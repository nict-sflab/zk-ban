package witness

import (
	"math/big"

	"github.com/akakou/zk-ban/primitives"
	"github.com/akakou/zk-ban/snark"
)

type UserPublicKey struct {
	Number *primitives.BigInt
}

// func (b UserPublicKey) MarshalJSON() ([]byte, error) {
// 	return b.BigInt.MarshalJSON()
// }

// func (b *UserPublicKey) UnmarshalJSON(p []byte) error {
// 	return b.BigInt.UnmarshalJSON(p)

// }

type OneTimeTicket struct {
	Number *primitives.BigInt
}

func (usk *UserSecretKey) PseudoRandom(period int64, name int64) (*primitives.BigInt, error) {
	hash, err := snark.CommitHash(big.NewInt(name), big.NewInt(period), &usk.Number.Int)

	if err != nil {
		return nil, err
	}

	return &primitives.BigInt{*hash}, nil
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
