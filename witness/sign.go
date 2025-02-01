package witness

import (
	"math/big"
)

type SignCommit struct {
	Commit1 *big.Int
	Commit2 *big.Int
}

type Signer struct {
	UserSecretKey *UserSecretKey
	UserPublicKey *UserPublicKey
	Certificate   *Certificate
	Period        *big.Int
}

func (signer *Signer) ComputeSignCommit(m, bsn *big.Int) (*SignCommit, error) {
	commit1, err := hash(m, signer.UserSecretKey.Number)
	if err != nil {
		return nil, err
	}

	commit2, err := hash(signer.Period, bsn, signer.UserSecretKey.Number)
	if err != nil {
		return nil, err
	}

	return &SignCommit{Commit1: commit1, Commit2: commit2}, nil
}

func (signer *Signer) NextSignerCandidate(nextPeriod *big.Int) (*Signer, error) {
	upk, err := signer.UserSecretKey.PublicKey(nextPeriod)
	if err != nil {
		return nil, err
	}

	return &Signer{
		UserSecretKey: signer.UserSecretKey,
		UserPublicKey: upk,
		Period:        nextPeriod,
		Certificate:   nil,
	}, nil
}
