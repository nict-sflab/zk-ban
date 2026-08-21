package circuit

import (
	"fmt"

	"github.com/akakou/zk-ban/primitives"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

type UpdateCircuit struct {
	UserSecretKey  frontend.Variable `gnark:",secret"`
	Credential     eddsa.Signature   `gnark:",secret"`
	CurrentInfo    PublicKeyAuthInfo
	NextInfo       PublicKeyAuthInfo
	GroupPublicKey eddsa.PublicKey `gnark:",public"`
	RevocationList RevocationList
}

func (circuit *UpdateCircuit) Define(api frontend.API) error {
	err := authCredential(api, circuit.UserSecretKey, circuit.Credential, circuit.CurrentInfo.Period, circuit.GroupPublicKey)
	if err != nil {
		return err
	}

	err = authPubKey(api, circuit.CurrentInfo, circuit.UserSecretKey, ONE_TIME_TICKET)
	if err != nil {
		return err
	}

	err = authPubKey(api, circuit.NextInfo, circuit.UserSecretKey, PUBLIC_KEY)
	if err != nil {
		return err
	}

	err = circuit.RevocationList.CheckRevocation(circuit.UserSecretKey, circuit.GroupPublicKey, api)
	if err != nil {
		return err
	}

	return nil
}

func NewUpdateCircuitWitness(
	nextPeriod int64,
	nextPublicKey *zkbanw.UserPublicKey,
	ticket *zkbanw.OneTimeTicket,
	signer *zkbanw.Signer,
	revocationList zkbanw.RevocationList,
	gpk *zkbanw.GroupPublicKey,
) (witness.Witness, error) {
	assign, err := newUpdateCircuitAssignment(nextPeriod, nextPublicKey, ticket, signer, gpk)
	if err != nil {
		return nil, err
	}

	proofList, err := revocationList.NonMembershipProofs(signer, MaxSession)
	if err != nil {
		return nil, err
	}
	assign.RevocationList, err = NewRevocationProofAssigned(proofList)
	if err != nil {
		return nil, err
	}

	return frontend.NewWitness(assign, snark.EcCurve.ScalarField())
}

func NewPublicUpdateCircuitWitness(
	nextPeriod int64,
	nextPublicKey *zkbanw.UserPublicKey,
	updateTicket *zkbanw.OneTimeTicket,
	lastPeriod int64,
	revocationList zkbanw.RevocationList,
	gpk *zkbanw.GroupPublicKey,
) (witness.Witness, error) {
	assign, err := newUpdateCircuitAssignment(nextPeriod, nextPublicKey, updateTicket, &zkbanw.Signer{
		UserSecretKey: &zkbanw.UserSecretKey{
			Number: primitives.NewBigInt(0),
		},
		Credential: &zkbanw.Credential{
			Signature: make([]byte, 32),
		},
		Period: lastPeriod,
	}, gpk)
	if err != nil {
		return nil, err
	}

	// Bounds and signatures are secret and are discarded by Public(). Concrete
	// zero placeholders are sufficient; only the period vector is public.
	assign.RevocationList = NewRevocationListAssigned(revocationList)

	wit, err := frontend.NewWitness(assign, snark.EcCurve.ScalarField())
	if err != nil {
		return nil, err
	}

	return wit.Public()
}

func newUpdateCircuitAssignment(
	nextPeriod int64,
	nextPublicKey *zkbanw.UserPublicKey,
	ticket *zkbanw.OneTimeTicket,
	signer *zkbanw.Signer,
	gpk *zkbanw.GroupPublicKey,
) (*UpdateCircuit, error) {
	if nextPublicKey == nil || nextPublicKey.Number == nil {
		return nil, fmt.Errorf("missing next public key")
	}
	if ticket == nil || ticket.Number == nil {
		return nil, fmt.Errorf("missing update ticket")
	}
	if signer == nil || signer.UserSecretKey == nil || signer.UserSecretKey.Number == nil || signer.Credential == nil {
		return nil, fmt.Errorf("missing signer witness")
	}
	if gpk == nil || gpk.PublicKey == nil {
		return nil, fmt.Errorf("missing group public key")
	}

	assign := &UpdateCircuit{
		UserSecretKey: signer.UserSecretKey.Number.Int,
		CurrentInfo: PublicKeyAuthInfo{
			Period:        signer.Period,
			UserPublicKey: ticket.Number.Int,
		},
		NextInfo: PublicKeyAuthInfo{
			Period:        nextPeriod,
			UserPublicKey: nextPublicKey.Number.Int,
		},
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	assign.Credential.Assign(snark.TwistededwardsCurve, signer.Credential.Signature)

	return assign, nil
}
