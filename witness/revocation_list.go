package witness

import (
	"math/big"
)

type RevocationList []RevokedNymsPerSession
type RevokedNymsPerSession struct {
	SessionTag *big.Int
	Nyms       []*big.Int
}

func EmptyConstantRevocationAddList(SessionSize, NymSizePerSession int) RevocationList {
	rl := RevocationList{}

	for i := 0; i < SessionSize; i++ {
		nyms := []*big.Int{}

		for j := 0; j < NymSizePerSession; j++ {
			n := big.NewInt(0)
			nyms = append(nyms, n)
		}

		n := big.NewInt(0)
		rl = append(rl, RevokedNymsPerSession{
			Nyms:       nyms,
			SessionTag: n,
		})
	}

	return rl
}

func EmptyRevocationList(NymSizePerSessions []int) RevocationList {
	rl := RevocationList{}

	for _, nymPerSession := range NymSizePerSessions {
		nyms := []*big.Int{}

		for i := 0; i < nymPerSession; i++ {
			n := big.NewInt(0)
			nyms = append(nyms, n)
		}

		n := big.NewInt(0)
		rl = append(rl, RevokedNymsPerSession{
			Nyms:       nyms,
			SessionTag: n,
		})
	}

	return rl
}
