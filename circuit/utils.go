package circuit

import (
	"math/big"

	"github.com/consensys/gnark/frontend"
)

func NewRevocationListWitness(revokedNyms [][]*big.Int, sessionTags []*big.Int) RevocationList {
	nyms := RevokedNyms{}
	tags := SessionTags{}

	for i, tag := range sessionTags {
		tags = append(tags, tag)
		nyms = append(nyms, []frontend.Variable{})

		for j := 0; j < RevokedNymsPerSession; j++ {
			nyms[i] = append(nyms[i], revokedNyms[i][j])
		}
	}

	return RevocationList{
		SessionTags: tags,
		RevokedNyms: nyms,
	}
}

func EmptyRevocationList(SessionSize, RevokedNymsPerSession int) RevocationList {
	nyms := RevokedNyms{}
	tags := SessionTags{}

	for i := 0; i < SessionSize; i++ {
		tags = append(tags, 0)
		nyms = append(nyms, []frontend.Variable{})

		for j := 0; j < RevokedNymsPerSession; j++ {
			nyms[i] = append(nyms[i], 0)
		}
	}

	return RevocationList{
		SessionTags: tags,
		RevokedNyms: nyms,
	}
}
