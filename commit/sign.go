package commit

import (
	"math/big"

	"github.com/akakou/zk-ban/snark"
)

type SignCommit struct {
	Commit1 *big.Int
	Commit2 *big.Int
}

type Signer struct {
	UserSecretKey  *UserSecretKey
	UserPublicKey  *UserPublicKey
	Credential     *Credential
	Period         *big.Int
	GroupPublicKey *GroupPublicKey
}

func (signer *Signer) CommitSign(m, bsn *big.Int) (*SignCommit, error) {
	commit1, err := snark.CommitHash(m, signer.UserSecretKey.Number)
	if err != nil {
		return nil, err
	}

	commit2, err := snark.CommitHash(signer.Period, bsn, signer.UserSecretKey.Number)
	if err != nil {
		return nil, err
	}

	return &SignCommit{Commit1: commit1, Commit2: commit2}, nil
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
