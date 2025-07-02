package circuit

import (
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/signature/eddsa"

	"github.com/consensys/gnark/frontend"
)

type CredentialAuthInfo struct {
	Period     frontend.Variable `gnark:",public"`
	Credential eddsa.Signature   `gnark:",secret"`
}

type PublicKeyAuthInfo struct {
	UserPublicKey frontend.Variable `gnark:",public"`
	Period        frontend.Variable `gnark:",public"`
}

func authCredential(api frontend.API, info CredentialAuthInfo, usk frontend.Variable, gpk eddsa.PublicKey) error {
	upk, err := snark.CircuitHash(api, info.Period, usk)
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

	err = eddsa.Verify(curve, info.Credential, upk, gpk, hash)
	if err != nil {
		return err
	}

	return nil
}

func authPubKey(api frontend.API, info PublicKeyAuthInfo, usk frontend.Variable) error {
	upk_dash, err := snark.CircuitHash(api, info.Period, usk)
	if err != nil {
		return err
	}

	api.AssertIsEqual(info.UserPublicKey, upk_dash)

	return nil
}
