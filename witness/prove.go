package zkban

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
)

type Proof struct {
	Hash []byte
}

func ComputeProveWitness(m, r *big.Int, usk *UserSecretKey) (*Proof, error) {
	hasher := mimc.NewMiMC(mimc.WithByteOrder(fr.BigEndian))
	_, err := hasher.Write(m.Bytes())
	if err != nil {
		return nil, err
	}

	_, err = hasher.Write(r.Bytes())
	if err != nil {
		return nil, err
	}

	_, err = hasher.Write(usk.UserSecretKey.Bytes())
	if err != nil {
		return nil, err
	}

	h := hasher.Sum(nil)

	return &Proof{Hash: h}, nil
}
