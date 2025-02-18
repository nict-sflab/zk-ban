package witness

import (
	"math/big"

	"github.com/akakou/zk-ban/snark"
)

type SignCommit struct {
	Sigma *big.Int
	Nym   *big.Int
}

type Signer struct {
	UserSecretKey  *UserSecretKey
	UserPublicKey  *UserPublicKey
	Credential     *Credential
	Period         *big.Int
	GroupPublicKey *GroupPublicKey
}

func (signer *Signer) CommitSign(m, sessionTag *big.Int) (*SignCommit, error) {
	commit1, err := snark.CommitHash(m, signer.UserSecretKey.Number)
	if err != nil {
		return nil, err
	}

	commit2, err := snark.CommitHash(sessionTag, signer.UserSecretKey.Number)
	if err != nil {
		return nil, err
	}

	return &SignCommit{Sigma: commit1, Nym: commit2}, nil
}

func (signer *Signer) NextWithoutCred(nextPeriod *big.Int, nextGpk *GroupPublicKey) (*Signer, error) {
	upk, err := signer.UserSecretKey.PublicKey(nextPeriod)
	if err != nil {
		return nil, err
	}

	return &Signer{
		UserSecretKey:  signer.UserSecretKey,
		UserPublicKey:  upk,
		Period:         nextPeriod,
		GroupPublicKey: nextGpk,
		Credential:     nil,
	}, nil
}

func (signer *Signer) SessionTag(counter *big.Int) *big.Int {
	sessionTagBytes := append(counter.Bytes(), signer.Period.Bytes()...)
	sessionTag := big.NewInt(0)
	sessionTag.SetBytes(sessionTagBytes)

	return sessionTag
}
