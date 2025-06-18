package circuit

import (
	"github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/frontend"
)

type RevokedNymsPerSession struct {
	Nyms       []frontend.Variable
	SessionTag frontend.Variable
}
type RevocationList []RevokedNymsPerSession

func EmptyRevocationList(SessionSize, NymSizePerSession int) RevocationList {
	emptyRL := witness.EmptyConstantRevocationAddList(SessionSize, NymSizePerSession)
	return NewRevocationListWitness(emptyRL)
}

func NewRevocationListWitness(revocationList witness.RevocationList) RevocationList {
	rl := RevocationList{}

	for _, rps := range revocationList {
		nyms := []frontend.Variable{}
		for _, nym := range rps.Nyms {
			nyms = append(nyms, nym)
		}

		rl = append(rl, RevokedNymsPerSession{
			Nyms:       nyms,
			SessionTag: rps.SessionTag,
		})
	}

	return rl
}
