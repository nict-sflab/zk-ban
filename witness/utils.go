package witness

import (
	"math/big"

	"github.com/akakou/zk-ban/snark"
	"github.com/iden3/go-iden3-crypto/poseidon"
)

type hashT = func(data ...*big.Int) ([]byte, error)

var hash = poseidonHash

func mimcHash(data ...*big.Int) (*big.Int, error) {
	hasher := snark.NewMIMC(snark.MimcWithByteOrder(snark.MimcFrBigEndian))

	for _, d := range data {
		_, err := hasher.Write(d.Bytes())
		if err != nil {
			return nil, err
		}
	}

	h := hasher.Sum(nil)

	i := big.NewInt(0)
	i.SetBytes(h)

	return i, nil
}

func poseidonHash(data ...*big.Int) (*big.Int, error) {
	hash, err := poseidon.Hash(data)

	if err != nil {
		return nil, err
	}

	return hash, nil

}
