package zkban

import (
	"github.com/consensys/gnark-crypto/hash"
)

type Certificate struct {
	Signature []byte
}

func (gsk *GroupSecretKey) IssuseCertificate(pk *UserPublicKey) (*Certificate, error) {
	hasher := hash.MIMC_BN254.New()
	signature, err := gsk.Sign(pk.UserPublicKey, hasher)
	if err != nil {
		return nil, err
	}

	cert := Certificate{Signature: signature}
	return &cert, nil
}
