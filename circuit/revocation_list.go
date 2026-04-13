package circuit

import (
	"github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/frontend"
)

type RevokedNymsPerPeriod struct {
	Nyms   []frontend.Variable `gnark:",public"`
	Period frontend.Variable   `gnark:",public"`
}
type RevocationList []RevokedNymsPerPeriod

func NewRevocationListAssigned(revocationList witness.RevocationList) RevocationList {
	rl := RevocationList{}

	for _, rps := range revocationList {
		nyms := []frontend.Variable{}
		for _, nym := range rps.Nyms {
			nyms = append(nyms, nym.Int)
		}

		rl = append(rl, RevokedNymsPerPeriod{
			Nyms:   nyms,
			Period: rps.Period.Int,
		})
	}

	return rl
}
