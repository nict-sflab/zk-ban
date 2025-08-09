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
		return nil, nil, err
	}

	gpk := gsk.Public()
	return &GroupSecretKey{gsk}, &GroupPublicKey{gpk}, nil
}

func GroupPublicKeyFromBytes(gpk []byte) (*GroupPublicKey, error) {
	gsk, err := eddsa.New(snark.TwistededwardsCurve, rand.Reader)
	if err != nil {
		return nil, err
	}

	gpkObj := gsk.Public()
	gpkObj.SetBytes(gpk)

	return &GroupPublicKey{gpkObj}, nil
}

func GroupSecretKeyFromBytes(gsk []byte) (*GroupSecretKey, error) {
	signer, err := eddsa.New(snark.TwistededwardsCurve, rand.Reader)
	if err != nil {
		return nil, err
	}

	_, err = signer.SetBytes(gsk)
	if err != nil {
		return nil, err
	}

	return &GroupSecretKey{signer}, nil
}
