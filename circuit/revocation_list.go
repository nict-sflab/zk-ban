package circuit

import (
	"github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/frontend"
)

type RevokedNymsPerSession struct {
	Nyms       []frontend.Variable `gnark:",public"`
	SessionTag frontend.Variable   `gnark:",public"`
}
type RevocationList []RevokedNymsPerSession

func NewRevocationListAssigned(revocationList witness.RevocationList) RevocationList {
	rl := RevocationList{}

	for _, rps := range revocationList {
		nyms := []frontend.Variable{}
		for _, nym := range rps.Nyms {
			nyms = append(nyms, nym.Int)
		}

		rl = append(rl, RevokedNymsPerSession{
			Nyms:       nyms,
			SessionTag: rps.SessionTag.Int,
		})
	}

	return rl
}
