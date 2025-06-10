package witness

import (
	hash "github.com/akakou/zk-ban/primitives/hash/witness"
)

type Credential struct {
	Signature []byte
}

func (gsk *GroupSecretKey) IssueCredential(upk *UserPublicKey) (*Credential, error) {
	signature, err := gsk.Sign(upk.Number.Bytes(), hash.HashType.New())
	if err != nil {
		return nil, err
	}

	return &Credential{
		signature,
	}, nil
}
