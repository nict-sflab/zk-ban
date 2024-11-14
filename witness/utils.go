package witness

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/iden3/go-iden3-crypto/poseidon"
)

type hashT = func(data ...[]byte) ([]byte, error)

var hash = mimcHash

func mimcHash(data ...[]byte) ([]byte, error) {
	hasher := mimc.NewMiMC(mimc.WithByteOrder(fr.BigEndian))

	for _, d := range data {
		_, err := hasher.Write(d)
		if err != nil {
			return nil, err
		}
	}

	return hasher.Sum(nil), nil
}

func poseidonHash(data ...[]byte) ([]byte, error) {
	buf := []byte{}

	for _, b := range data {
		buf = append(buf, b...)
	}
	h, err := poseidon.HashBytes(buf)
	if err != nil {
		return nil, err
	}
	return h.Bytes(), nil
}
