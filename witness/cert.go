package witness

import "github.com/akakou/zk-ban/snark"

type Certificate struct {
	Signature []byte
}

func (gsk *GroupSecretKey) IssueCertificate(upk *UserPublicKey) (*Certificate, error) {
	hasher := snark.HashAlg.New()
	signature, err := gsk.Sign(upk.Number.Bytes(), hasher)
	if err != nil {
		return nil, err
	}

	cert := Certificate{Signature: signature}
	return &cert, nil
}
