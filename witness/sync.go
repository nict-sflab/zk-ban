package witness

import "math/big"

func ComputeSyncCommit(nonce *big.Int, usk *UserSecretKey) (*big.Int, error) {
	commit, err := hash(nonce, usk.Number)
	if err != nil {
		return nil, err
	}

	return commit, nil
}
