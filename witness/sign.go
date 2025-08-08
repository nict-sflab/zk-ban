package witness

import (
	"math/big"

	"github.com/akakou/zk-ban/primitives"
	"github.com/akakou/zk-ban/snark"
)

type SignCommit struct {
	Sigma *primitives.BigInt
	Nym   *primitives.BigInt
}

type Signer struct {
	UserSecretKey *UserSecretKey
	Credential    *Credential
	Period        int64
}

func (signer *Signer) CommitSign(m, sessionTag *primitives.BigInt) (*SignCommit, error) {
	commit1, err := snark.CommitHash(&m.Int, &signer.UserSecretKey.Int)
	if err != nil {
		return nil, err
	}

	commit2, err := snark.CommitHash(&sessionTag.Int, &signer.UserSecretKey.Int)
	if err != nil {
		return nil, err
	}

	return &SignCommit{Sigma: &primitives.BigInt{*commit1}, Nym: &primitives.BigInt{*commit2}}, nil
}

// func (signer *Signer) NextWithoutCred(nextPeriod *big.Int) (*Signer, error) {
// 	signer = &Signer{
// 		UserSecretKey: signer.UserSecretKey,
// 		Period:        nextPeriod,
// 		Credential:    nil,
// 	}

// 	return signer, nil
// }

func (signer *Signer) SessionTag(counter int64) *primitives.BigInt {
	sessionTag := SessionTag(counter, signer.Period)
	return sessionTag
}

func SessionTag(counter int64, period int64) *primitives.BigInt {
	sessionTagBytes := append(big.NewInt(counter).Bytes(), big.NewInt(period).Bytes()...)
	sessionTag := big.NewInt(0)
	sessionTag.SetBytes(sessionTagBytes)

	return &primitives.BigInt{*sessionTag}
}

func (signer *Signer) PublicKey() (*UserPublicKey, error) {
	return signer.UserSecretKey.PublicKey(signer.Period)
}

func (signer *Signer) OneTimeTicket() (*OneTimeTicket, error) {
	return signer.UserSecretKey.OneTimeTicket(signer.Period)
}
