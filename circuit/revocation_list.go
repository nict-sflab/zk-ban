package circuit

import (
	"github.com/akakou/zk-ban/commit"
	"github.com/consensys/gnark/frontend"
)

type RevokedNymsPerSession struct {
	Nyms     []frontend.Variable `gnark:",public"`
	Period   frontend.Variable   `gnark:",public"`
	Basename frontend.Variable   `gnark:",public"`
}
type RevocationList []RevokedNymsPerSession

func EmptyRevocationList(SessionSize, NymSizePerSession int) RevocationList {
	emptyRL := commit.EmptyConstantRevocationAddList(SessionSize, NymSizePerSession)
	return NewRevocationListWitness(emptyRL)
}

func NewRevocationListWitness(revocationList commit.RevocationList) RevocationList {
	rl := RevocationList{}

	for _, rps := range revocationList {
		nyms := []frontend.Variable{}
		for _, nym := range rps.Nyms {
			nyms = append(nyms, nym)
		}

		rl = append(rl, RevokedNymsPerSession{
			Nyms:     nyms,
			Period:   rps.Period,
			Basename: rps.Basename,
		})
	}

	return rl
}
