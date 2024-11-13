package witness

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
)

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
