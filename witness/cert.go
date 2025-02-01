package witness

import "github.com/akakou/zk-ban/snark"

type Certificate struct {
	Signature []byte
}

func (gsk *GroupSecretKey) IssueCertificate(upk *UserPublicKey) (*Certificate, error) {
	signature, err := gsk.Sign(upk.Number.Bytes(), snark.Hasher)
	if err != nil {
		return nil, err
	}

	return &Certificate{
		signature,
	}, nil
}
