package witness

import "github.com/liyue201/gnark-circomlib/utils/poseidon"

type Certificate struct {
	Signature []byte
}

func (gsk *GroupSecretKey) IssueCertificate(upk *UserPublicKey) (*Certificate, error) {
	// generate signature
	hFunc := poseidon.NewPoseidon()

	signature, err := gsk.Sign(upk.Number.Bytes(), hFunc)
	if err != nil {
		return nil, err
	}

	return &Certificate{
		signature,
	}, nil
}
