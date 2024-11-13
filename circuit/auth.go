package circuit

import (
	"github.com/consensys/gnark/std/signature/eddsa"

	tw "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
)

func auth(api frontend.API, period, usk, upk frontend.Variable, cert eddsa.Signature, gpk eddsa.PublicKey) error {
	err := pubKeyAuth(api, period, usk, upk)
	if err != nil {
		return err
	}

	err = certAuth(api, upk, cert, gpk)
	if err != nil {
		return err
	}

	return nil
}

func pubKeyAuth(api frontend.API, period, usk, upk frontend.Variable) error {
	upk_dash, err := mimcHash(api, period, usk)
	if err != nil {
		return err
	}

	api.AssertIsEqual(upk, upk_dash)

	return nil
}

func certAuth(api frontend.API, upk frontend.Variable, cert eddsa.Signature, gpk eddsa.PublicKey) error {
	mimc0, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	curve, err := twistededwards.NewEdCurve(api, tw.BN254)
	if err != nil {
		return err
	}

	err = eddsa.Verify(curve, cert, upk, gpk, &mimc0)
	if err != nil {
		return err
	}

	return nil
}
