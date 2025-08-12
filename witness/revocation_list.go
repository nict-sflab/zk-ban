package witness

import (
	"log"
	"math/big"

	"github.com/akakou/zk-ban/primitives"
	"github.com/akakou/zk-ban/snark"
)

type RevocationList []RevokedNymsPerSession
type RevokedNymsPerSession struct {
	SessionTag *primitives.BigInt
	Nyms       []*primitives.BigInt
}

var InitBigInt = ZeroInitBigInt
var ZeroInitBigInt = func() *primitives.BigInt { return primitives.NewBigInt(0) }
var MimcInitBigInt = func() *primitives.BigInt {
	n, err := snark.CommitHash(big.NewInt(202506252))
	if err != nil {
		log.Fatal(err)
	}
	return &primitives.BigInt{*n}
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
		nyms := []*primitives.BigInt{}

		for i := 0; i < nymPerSession; i++ {
			n := InitBigInt()
			nyms = append(nyms, n)
		}

		n := primitives.NewBigInt(0)
		rl = append(rl, RevokedNymsPerSession{
			Nyms:       nyms,
			SessionTag: n,
		})
	}

	return rl
}
