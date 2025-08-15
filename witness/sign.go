package witness

import (
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
	commit1, err := snark.CommitHash(&m.Int, &signer.UserSecretKey.Number.Int)
	if err != nil {
		return nil, err
	}

	commit2, err := snark.CommitHash(&sessionTag.Int, &signer.UserSecretKey.Number.Int)
	if err != nil {
		return nil, err
	}

	return &SignCommit{Sigma: &primitives.BigInt{*commit1}, Nym: &primitives.BigInt{*commit2}}, nil
}

func (signer *Signer) SessionTag(counter int64) *primitives.BigInt {
	sessionTag := SessionTag(counter, signer.Period)
	return sessionTag
}

func SessionTag(counter int64, period int64) *primitives.BigInt {
	var counterBuf = [8]byte{}
	var periodBuf = [8]byte{}

	copy(counterBuf[:], primitives.NewBigInt(counter).Bytes())
	copy(periodBuf[:], primitives.NewBigInt(period).Bytes())

	concated := append(counterBuf[:], periodBuf[:]...)
	return primitives.BigIntFromBytes(concated)
}

func (signer *Signer) PublicKey() (*UserPublicKey, error) {
	return signer.UserSecretKey.PublicKey(signer.Period)
}

func (signer *Signer) OneTimeTicket() (*OneTimeTicket, error) {
	return signer.UserSecretKey.OneTimeTicket(signer.Period)
}
