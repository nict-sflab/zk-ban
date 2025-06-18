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
			nyms = append(nyms, big.NewInt(1))
		}

		rl = append(rl, RevokedNymsPerSession{
			Nyms:       nyms,
			SessionTag: big.NewInt(1),
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
			Nyms:       nyms,
			SessionTag: big.NewInt(0),
		})
	}

	return rl
}

func HashRevocationList(rl RevocationList) *big.Int {
	var aa []*big.Int
	for _, r := range rl {
		aa = append(aa, r.SessionTag)
		aa = append(aa, r.Nyms...)
	}

	h, err := snark.CommitHash(aa...)
	if err != nil {
		panic(err)
	}

	return h
}
