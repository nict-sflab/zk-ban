package zkban

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
)

type UserPublicKey struct {
	UserPublicKey []byte
}

func (usk *UserSecretKey) PublicKey(period *big.Int) (*UserPublicKey, error) {
	hasher := mimc.NewMiMC(mimc.WithByteOrder(fr.BigEndian))

	_, err := hasher.Write(period.Bytes())
	if err != nil {
		return nil, err
	}

	_, err = hasher.Write(usk.UserSecretKey.Bytes())
	if err != nil {
		return nil, err
	}

	pk := hasher.Sum(nil)

	return &UserPublicKey{UserPublicKey: pk}, nil
}
