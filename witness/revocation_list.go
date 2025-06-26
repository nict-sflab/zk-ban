package witness

import (
	"math/big"

	"github.com/akakou/zk-ban/snark"
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
			n, err := snark.CommitHash(big.NewInt(202506253))
			if err != nil {
				panic(err)
			}
			n = big.NewInt(0)

			nyms = append(nyms, n)
		}

		n, err := snark.CommitHash(big.NewInt(202506253))
		if err != nil {
			panic(err)
		}
		n = big.NewInt(0)

		rl = append(rl, RevokedNymsPerSession{
			Nyms:       nyms,
			SessionTag: n,
		})
	}

	return rl
}
