package circuit

import (
	"math/big"

	"github.com/akakou/zk-ban/commit"
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark-crypto/signature"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

// const RevocationListSize = 223
const RevocationSpeed = 0.82

const RevokedNymsPerSession = 130
const SessionSize = 270

type RevokedNyms [SessionSize][RevokedNymsPerSession]frontend.Variable
type SessionTags [SessionSize]frontend.Variable

type RevocationList struct {
	SessionTags SessionTags `gnark:",public"`
	RevokedNyms RevokedNyms `gnark:",public"`
}

type UpdateCircuit struct {
	UserSecretKey  frontend.Variable `gnark:",secret"`
	CurrentInfo    CredentialAuthInfo
	NextInfo       PublicKeyAuthInfo
	RevocationList RevocationList
	GroupPublicKey eddsa.PublicKey `gnark:",public"`
}

var Scan = ScanConstantRevocations

func ScanConstantRevocations(nym frontend.Variable, sessionId int, revokedNyms RevokedNyms, api frontend.API) {
	for j := 0; j < RevokedNymsPerSession; j++ {
		api.AssertIsDifferent(nym, revokedNyms[sessionId][j])
	}
}
func ScanIncrementalRevocations(nym frontend.Variable, sessionId int, revokedNyms RevokedNyms, api frontend.API) {
	max := int(RevocationSpeed * float64(sessionId))
	for j := 0; j < max; j++ {
		api.AssertIsDifferent(nym, revokedNyms[sessionId][j])
	}
}

func (circuit *UpdateCircuit) Define(api frontend.API) error {
	err := authCredential(api, circuit.CurrentInfo, circuit.UserSecretKey, circuit.GroupPublicKey)
	if err != nil {
		return err
	}

	err = authPubKey(api, circuit.NextInfo, circuit.UserSecretKey)
	if err != nil {
		return err
	}

	for i := 0; i < SessionSize; i++ {
		nym, err := snark.CircuitHash(api, circuit.RevocationList.SessionTags[i], circuit.UserSecretKey)
		if err != nil {
			return err
		}

		Scan(nym, i, circuit.RevocationList.RevokedNyms, api)
	}

	return nil
}

func NewUpdateCircuitWitness(revokedNyms [SessionSize][RevokedNymsPerSession]*big.Int, sessionTags [SessionSize]*big.Int, next, last *commit.Signer, gpk signature.PublicKey) *UpdateCircuit {
	assign := &UpdateCircuit{
		UserSecretKey: last.UserSecretKey.Number,
		CurrentInfo: CredentialAuthInfo{
			Period: last.Period,
		},

		NextInfo: PublicKeyAuthInfo{
			Period:        next.Period,
			UserPublicKey: next.UserPublicKey.Number,
		},
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	assign.CurrentInfo.Credential.Assign(snark.TwistededwardsCurve, last.Credential.Signature)

	for i := 0; i < SessionSize; i++ {
		assign.RevocationList.SessionTags[i] = sessionTags[i]

		for j := 0; j < RevokedNymsPerSession; j++ {
			assign.RevocationList.RevokedNyms[i][j] = revokedNyms[i][j]
		}
	}

	return assign
}
