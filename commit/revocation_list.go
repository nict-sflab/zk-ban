package commit

import "math/big"

type RevocationList []RevokedNymsPerSession
type RevokedNymsPerSession struct {
	SessionTag *big.Int
	Nyms       []*big.Int
}

func EmptyRevocationList(SessionSize, NymSizePerSession int) RevocationList {
	rl := RevocationList{}

	for i := 0; i < SessionSize; i++ {
		nyms := []*big.Int{}
		for j := 0; j < NymSizePerSession; j++ {
			nyms = append(nyms, big.NewInt(0))
		}

		rl = append(rl, RevokedNymsPerSession{
			Nyms:       nyms,
			SessionTag: big.NewInt(0),
		})
	}

	return rl
}
