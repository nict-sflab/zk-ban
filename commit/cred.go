package commit

import "github.com/akakou/zk-ban/snark"

type Credential struct {
	Signature []byte
}

func (gsk *GroupSecretKey) IssueCredential(upk *UserPublicKey) (*Credential, error) {
	signature, err := gsk.Sign(upk.Number.Bytes(), snark.NewCommitHash())
	if err != nil {
		return nil, err
	}

	return &Credential{
		signature,
	}, nil
}
