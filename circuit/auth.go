package circuit

import (
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/signature/eddsa"

	"github.com/consensys/gnark/frontend"
)

const (
	PUBLIC_KEY = iota + 1
	ONE_TIME_TICKET
)

type CredentialAuthInfo struct {
	Credential eddsa.Signature   `gnark:",secret"`
	Period     frontend.Variable `gnark:",public"`
}

type PublicKeyAuthInfo struct {
	UserPublicKey frontend.Variable `gnark:",public"`
	Period        frontend.Variable `gnark:",public"`
}

func authCredential(api frontend.API, usk frontend.Variable, cred eddsa.Signature, period frontend.Variable, gpk eddsa.PublicKey) error {
	upk, err := snark.CircuitHash(api, PUBLIC_KEY, period, usk)
	if err != nil {
		return err
	}

	curve, err := twistededwards.NewEdCurve(api, snark.TwistededwardsCurve)
	if err != nil {
		return err
	}

	hash, err := snark.NewCircuitHash(api)
	if err != nil {
		return err
	}

	err = eddsa.Verify(curve, cred, upk, gpk, hash)
	if err != nil {
		return err
	}

	return nil
}

func authPubKey(api frontend.API, info PublicKeyAuthInfo, usk frontend.Variable, name frontend.Variable) error {
	upk_dash, err := snark.CircuitHash(api, name, info.Period, usk)
	if err != nil {
		return err
	}

	api.AssertIsEqual(info.UserPublicKey, upk_dash)

	return nil
}
