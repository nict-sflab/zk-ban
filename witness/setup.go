package witness

import (
	"crypto/rand"
	"errors"

	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark-crypto/signature"
	"github.com/consensys/gnark-crypto/signature/eddsa"
)

var ErrCreateGMSecretKey = errors.New("failed to create GM's secret key: ")

type GroupPublicKey struct{ signature.PublicKey }
type GroupSecretKey struct{ signature.Signer }

func RandomGroupKeyPair() (*GroupSecretKey, *GroupPublicKey, error) {
	gsk, err := eddsa.New(snark.TwistededwardsCurve, rand.Reader)
	if err != nil {
		return nil, nil, errors.Join(err)
	}

	gpk := gsk.Public()

	return &GroupSecretKey{gsk}, &GroupPublicKey{gpk}, nil
}
