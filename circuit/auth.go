package circuit

import (
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/signature/eddsa"

	"github.com/consensys/gnark/frontend"
)

func certAuth(api frontend.API, period, usk frontend.Variable, cert eddsa.Signature, gpk eddsa.PublicKey) error {
	upk, err := snark.CircuitHash(api, period, usk)
	if err != nil {
		return err
	}

	curve, err := twistededwards.NewEdCurve(api, snark.TwistededwardsCurve)
	if err != nil {
		return err
	}

	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	// poseidon := circuits.NewPoseidonHash(api)

	err = eddsa.Verify(curve, cert, upk, gpk, &mimc)
	if err != nil {
		return err
	}

	// circuits.EdDSAPoseidonVerifier()

	return nil
}

func pubKeyAuth(api frontend.API, period, usk, upk frontend.Variable) error {
	upk_dash, err := snark.CircuitHash(api, period, usk)
	if err != nil {
		return err
	}

	api.AssertIsEqual(upk, upk_dash)

	return nil
}
