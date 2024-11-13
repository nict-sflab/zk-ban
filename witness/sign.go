package witness

import (
	"math/big"
)

type SignCommit struct {
	Commit1 []byte
	Commit2 []byte
	Commit3 []byte
}

type Signer struct {
	UserSecretKey *UserSecretKey
	UserPublicKey *UserPublicKey
	Certificate   *Certificate
}

func ComputeSignCommit(m, bsn, nonce *big.Int, usk *UserSecretKey) (*SignCommit, error) {
	commit1, err := mimcHash(bsn.Bytes(), usk.Number.Bytes())
	if err != nil {
		return nil, err
	}
	commit2, err := mimcHash(m.Bytes(), nonce.Bytes())
	if err != nil {
		return nil, err
	}

	commit3, err := mimcHash(commit2, usk.Number.Bytes())
	if err != nil {
		return nil, err
	}

	return &SignCommit{Commit1: commit1, Commit2: commit2, Commit3: commit3}, nil
}
