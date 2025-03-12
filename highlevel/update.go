package highlevel

import (
	"math/big"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	"github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/signature/eddsa"
)

func UpdateRequest(
	nextPeriod int64,
	signer *HighLevelSigner,
	rl witness.RevocationList,
	gpk []byte,
	prover *HighLevelSnarkProver,
) (*HighLevelSigner, []byte, error) {
	nextPeriodBig := big.NewInt(nextPeriod)

	gpkObj, err := witness.GroupPublicKeyFromBytes(gpk)
	if err != nil {
		return nil, nil, err
	}

	signerObj, err := signer.ToSigner()
	if err != nil {
		return nil, nil, err
	}

	proverObj, err := prover.ToSnarkProver()
	if err != nil {
		return nil, nil, err
	}

	nextSignerObj, proof, _, err := zkban.UpdateRequest(nextPeriodBig, signerObj, rl, gpkObj, proverObj)

	if err != nil {
		return nil, nil, err
	}

	nextSigner := HighLevelSigner{}
	nextSigner.FromSigner(nextSignerObj)

	proofBytes, err := snark.EncodeProof(proof)

	return &nextSigner, proofBytes, err
}

func VerifyUpdateRequest(
	proofBytes []byte,
	nextupk []byte,
	nextPeriod int64,
	beforePeriod int64,
	rl witness.RevocationList,
	gpk []byte,
	verifyKey []byte,
) error {
	nextPerioBig := big.NewInt(nextPeriod)

	verifyKeyObj, err := snark.DecodeVerifierKey(verifyKey)
	if err != nil {
		return err
	}

	proof, err := snark.DecodeProof(proofBytes)
	if err != nil {
		return err
	}

	assign := circuit.UpdateCircuit{
		UserSecretKey: 0,
		CurrentInfo: circuit.CredentialAuthInfo{
			Credential: eddsa.Signature{
				R: twistededwards.Point{
					X: 0,
					Y: 0,
				},
				S: 0,
			},
			Period: beforePeriod,
		},
		NextInfo: circuit.PublicKeyAuthInfo{
			Period:        nextPerioBig,
			UserPublicKey: big.NewInt(0).SetBytes(nextupk),
		},
		RevocationList: circuit.NewRevocationListWitness(rl),
		GroupPublicKey: eddsa.PublicKey{},
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk)

	wit, err := frontend.NewWitness(&assign, snark.EcCurve.ScalarField())
	if err != nil {
		return err
	}

	pubWit, err := wit.Public()
	if err != nil {
		return err
	}

	err = groth16.Verify(proof, verifyKeyObj, pubWit)

	return err
}
