package commit

import (
	"math/big"
)

type RevocationList []RevokedNymsPerSession
type RevokedNymsPerSession struct {
	Period   *big.Int
	Basename *big.Int
	Nyms     []*big.Int
}

func EmptyConstantRevocationAddList(SessionSize, NymSizePerSession int) RevocationList {
	rl := RevocationList{}

	for i := 0; i < SessionSize; i++ {
		nyms := []*big.Int{}
		for j := 0; j < NymSizePerSession; j++ {
			nyms = append(nyms, big.NewInt(0))
		}

		rl = append(rl, RevokedNymsPerSession{
			Nyms:     nyms,
			Period:   big.NewInt(0),
			Basename: big.NewInt(0),
		})
	}

	return rl
}

func EmptyLinerRevocationAddList(SessionSize, RevocationListSize int) RevocationList {
	rl := RevocationList{}

	a := 2 * RevocationListSize / (SessionSize * SessionSize)

	for i := 0; i < SessionSize; i++ {
		nyms := []*big.Int{}
		for j := 0; j < a*i; j++ {
			nyms = append(nyms, big.NewInt(0))
		}

		rl = append(rl, RevokedNymsPerSession{
			Nyms:     nyms,
			Period:   big.NewInt(0),
			Basename: big.NewInt(0),
		})
	}

	return rl
}
