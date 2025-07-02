package highlevel

import (
	"math/big"

	gnarkprecomputes "github.com/akakou/gnark-precomputes"
	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	"github.com/akakou/zk-ban/witness"
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
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

func PrepareVerification(
	// proofBytes []byte,
	// nextupk []byte,
	// nextPeriod int64,
	// beforePeriod int64,
	// gpk []byte,
	rl witness.RevocationList,
	verifyKey []byte,
) ([]byte, error) {
	// nextPerioBig := big.NewInt(nextPeriod)

	verifyKeyObj, err := snark.DecodeVerifierKey(verifyKey)
	if err != nil {
		return nil, err
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
			// Period: beforePeriod,
			Period: 0,
		},
		NextInfo: circuit.PublicKeyAuthInfo{
			// Period:        nextPerioBig,
			// UserPublicKey: big.NewInt(0).SetBytes(nextupk),
			Period:        0,
			UserPublicKey: 0,
		},
		RevocationList: circuit.NewRevocationListWitness(rl),
		GroupPublicKey: eddsa.PublicKey{
			A: twistededwards.Point{
				X: 0,
				Y: 0,
			},
		},
	}

	wit, err := frontend.NewWitness(&assign, snark.EcCurve.ScalarField())
	if err != nil {
		return nil, err
	}

	pubWit, err := wit.Public()
	if err != nil {
		return nil, err
	}

	vk, err := gnarkprecomputes.FromBLS12381GnarkKey(verifyKeyObj, &assign)
	if err != nil {
		return nil, err
	}

	prepare, err := vk.PreparePublicInputs(pubWit)
	if err != nil {
		return nil, err
	}

	var encoded bls12381.G1Affine
	encoded.FromJacobian(prepare)
	res := encoded.Bytes()

	return res[:], err
}

func VerifyUpdateRequest(
	proofBytes []byte,
	nextupk []byte,
	nextPeriod int64,
	beforePeriod int64,
	rl witness.RevocationList,
	gpk []byte,
	prepared []byte,
	verifyKey []byte,
) error {
	var g1Aff bls12381.G1Affine
	_, err := g1Aff.SetBytes(prepared)
	if err != nil {
		return err
	}

	var g1Jac bls12381.G1Jac
	g1Jac.FromAffine(&g1Aff)
	if err != nil {
		return err
	}

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

	vk, err := gnarkprecomputes.FromBLS12381GnarkKey(verifyKeyObj, &assign)
	if err != nil {
		return err
	}

	err = vk.VerifyPrepared(proof, pubWit, &g1Jac)

	return err
}
