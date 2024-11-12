package zkban

import (
	"crypto/rand"
	"errors"

	"github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark-crypto/signature"
	"github.com/consensys/gnark-crypto/signature/eddsa"
)

var ErrCreateGMSecretKey = errors.New("failed to create GM's secret key: ")

type GroupPublicKey struct{ signature.PublicKey }
type GroupSecretKey struct{ signature.Signer }

func RandomGroupKeyPair() (*GroupSecretKey, *GroupPublicKey, error) {
	gsk, err := eddsa.New(twistededwards.BN254, rand.Reader)
	if err != nil {
		return nil, nil, errors.Join(err)
	}

	gpk := gsk.Public()

	return &GroupSecretKey{gsk}, &GroupPublicKey{gpk}, nil
}
