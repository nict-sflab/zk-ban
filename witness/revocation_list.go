package witness

import (
	"math/big"
)

type RevocationList []RevokedNymsPerSession
type RevokedNymsPerSession struct {
	SessionTag *big.Int
	Nyms       []*big.Int
}

func EmptyConstantRevocationAddList(sessionSize, nymSizePerSession int) RevocationList {
	rlSize := make([]int, sessionSize)
	for i := range sessionSize {
		rlSize[i] = nymSizePerSession
	}

	return EmptyRevocationList(rlSize)
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
