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

type RevokedNyms [][]frontend.Variable
type SessionTags []frontend.Variable

type UpdateCircuit struct {
	UserSecretKey     frontend.Variable `gnark:"sk"`
	LastCredential    eddsa.Signature   `gnark:"cert"`
	LastPeriod        frontend.Variable `gnark:",public"`
	GroupPublicKey    eddsa.PublicKey   `gnark:",public"`
	NextUserPublicKey frontend.Variable `gnark:",public"`
	NextPeriod        frontend.Variable `gnark:",public"`
	SessionTags       SessionTags       `gnark:",public"`
	RevokedNyms       RevokedNyms       `gnark:",public"`
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
	err := authCert(api, circuit.LastPeriod, circuit.UserSecretKey, circuit.LastCredential, circuit.GroupPublicKey)
	if err != nil {
		return err
	}

	err = authPubKey(api, circuit.NextPeriod, circuit.UserSecretKey, circuit.NextUserPublicKey)
	if err != nil {
		return err
	}

	for i := 0; i < SessionSize; i++ {
		nym, err := snark.CircuitHash(api, circuit.SessionTags[i], circuit.UserSecretKey)
		if err != nil {
			return err
		}

		Scan(nym, i, circuit.RevokedNyms, api)
	}

	return nil
}

type DefaultRevokedNyms [SessionSize][RevokedNymsPerSession]frontend.Variable
type DefaultSessionTags [SessionSize]frontend.Variable

func NewUpdateCircuitWitness(next, last *commit.Signer, revokedNyms [][]frontend.Variable, sessionTags []frontend.Variable, gpk signature.PublicKey) *UpdateCircuit {
	assign := &UpdateCircuit{
		UserSecretKey:     last.UserSecretKey.Number,
		LastPeriod:        last.Period,
		NextPeriod:        next.Period,
		NextUserPublicKey: next.UserPublicKey.Number,
		RevokedNyms:       revokedNyms,
		SessionTags:       sessionTags,
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	assign.LastCredential.Assign(snark.TwistededwardsCurve, last.Credential.Signature)

	return assign
}

func NewDefaultRevocationListWitness(revokedNyms [SessionSize][RevokedNymsPerSession]*big.Int, sessionTags [SessionSize]*big.Int) (SessionTags, RevokedNyms) {
	nyms := RevokedNyms{}
	tags := SessionTags{}

	for i := 0; i < SessionSize; i++ {
		tags = append(tags, sessionTags[i])
		nyms = append(nyms, []frontend.Variable{})

		for j := 0; j < RevokedNymsPerSession; j++ {
			nyms[i] = append(nyms[i], revokedNyms[i][j])
		}
	}

	return tags, nyms
}
