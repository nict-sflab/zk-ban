package highlevel

import (
	"math/big"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

type JoinReq struct {
	UserSecretKey []byte
	UserPublicKey []byte
	Period        int64
}

func JoinRequest(period int64, circuitBytes, proveKey []byte) ([]byte, *JoinReq, error) {
	periodBig := big.NewInt(period)

	cs, err := snark.DecodeCircuit(circuitBytes)
	if err != nil {
		return nil, nil, err
	}

	proveKeyObj, err := snark.DecodeProverKey(proveKey)
	if err != nil {
		return nil, nil, err
	}

	prover := snark.SnarkProver{
		ConstraintSystem: cs,
		ProveKey:         proveKeyObj,
	}

	proof, assign, err := zkban.JoinRequest(periodBig, &prover)
	if err != nil {
		return nil, nil, err
	}

	proofBytes, err := snark.EncodeProof(proof)
	if err != nil {
		return nil, nil, err
	}

	req := JoinReq{
		UserSecretKey: assign.UserSecretKey.(*big.Int).Bytes(),
		UserPublicKey: assign.PublicKeyAuthInfo.UserPublicKey.(*big.Int).Bytes(),
		Period:        assign.PublicKeyAuthInfo.Period.(*big.Int).Int64(),
	}

	return proofBytes, &req, err
}

func VerifyJoinReq(proof, upk []byte, period int64, verifyKey []byte) error {
	verifyKeyObj, err := snark.DecodeVerifierKey(verifyKey)
	if err != nil {
		return err
	}

	proofObj, err := snark.DecodeProof(proof)
	if err != nil {
		return err
	}

	assign := circuit.JoinRequestCircuit{
		UserSecretKey: big.NewInt(0),
		PublicKeyAuthInfo: circuit.PublicKeyAuthInfo{
			UserPublicKey: big.NewInt(0).SetBytes(upk),
			Period:        period,
		},
	}

	wit, err := frontend.NewWitness(&assign, snark.EcCurve.ScalarField())
	if err != nil {
		return err
	}

	pubWit, err := wit.Public()
	if err != nil {
		return err
	}

	err = groth16.Verify(proofObj, verifyKeyObj, pubWit)

	return err
}
