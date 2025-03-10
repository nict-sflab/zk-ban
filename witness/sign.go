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
	UserSecretKey *UserSecretKey
	UserPublicKey *UserPublicKey
	Credential    *Credential
	Period        *big.Int
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

func (signer *Signer) NextWithoutCred(nextPeriod *big.Int) (*Signer, error) {
	upk, err := signer.UserSecretKey.PublicKey(nextPeriod)
	if err != nil {
		return nil, err
	}

	return &Signer{
		UserSecretKey: signer.UserSecretKey,
		UserPublicKey: upk,
		Period:        nextPeriod,
		Credential:    nil,
	}, nil
}

func (signer *Signer) SessionTag(counter *big.Int) *big.Int {
	sessionTag := SessionTag(counter, signer.Period)
	return sessionTag
}

func SessionTag(counter *big.Int, period *big.Int) *big.Int {
	sessionTagBytes := append(counter.Bytes(), period.Bytes()...)
	sessionTag := big.NewInt(0)
	sessionTag.SetBytes(sessionTagBytes)

	return sessionTag
}
