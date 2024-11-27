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

func ComputeSignCommit(m, bsn, period *big.Int, usk *UserSecretKey) (*SignCommit, error) {
	commit1, err := hash(m, usk.Number)
	if err != nil {
		return nil, err
	}

	commit2, err := hash(period, bsn, usk.Number)
	if err != nil {
		return nil, err
	}

	return &SignCommit{Commit1: commit1, Commit2: commit2}, nil
}
