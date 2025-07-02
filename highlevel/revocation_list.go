package highlevel

import (
	"math/big"

	zkbanw "github.com/akakou/zk-ban/witness"
)

type HighLevelRevokedNymsPerSession struct {
	Nyms       [][]byte
	SessionTag []byte
}

type HighLevelRevocationList []HighLevelRevokedNymsPerSession

func (hRL *HighLevelRevocationList) ToRevocationList() (*zkbanw.RevocationList, error) {
	rl := zkbanw.RevocationList{}

	for _, h := range *hRL {
		nyms := []*big.Int{}

		for _, b := range h.Nyms {
			nym := big.NewInt(0)
			nym.FillBytes(b)

			nyms = append(nyms, nym)
		}

		tag := big.NewInt(0)
		tag.FillBytes(h.SessionTag)

		rl = append(rl, zkbanw.RevokedNymsPerSession{
			Nyms:       nyms,
			SessionTag: tag,
		})
	}

	return &rl, nil
}

func (hRL *HighLevelRevocationList) FromRevocationList(rlObj *zkbanw.RevocationList) {
	*hRL = HighLevelRevocationList{}

	for _, o := range *rlObj {
		nyms := [][]byte{}

		for _, nym := range o.Nyms {
			nb := nym.Bytes()
			nyms = append(nyms, nb)
		}

		tb := o.SessionTag.Bytes()

		*hRL = append(*hRL, HighLevelRevokedNymsPerSession{
			Nyms:       nyms,
			SessionTag: tb,
		})
	}
}
