package highlevel

import (
	"math/big"

	core "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	"github.com/akakou/zk-ban/witness"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

type Signature struct {
	Proof []byte
	Nym   []byte
	Sigma []byte
}

type HighLevelSigner struct {
	Credential     []byte
	Secret         []byte
	UserPublicKey  []byte
	GroupPublicKey []byte
	Period         int64
}

func (signer *HighLevelSigner) ToSigner() (*zkbanw.Signer, error) {
	periodBig := big.NewInt(signer.Period)
	upkBig := big.NewInt(0).SetBytes(signer.UserPublicKey)
	secretBig := big.NewInt(0).SetBytes(signer.Secret)

	signerStruct := zkbanw.Signer{
		UserSecretKey: &zkbanw.UserSecretKey{Number: secretBig},
		UserPublicKey: &zkbanw.UserPublicKey{Number: upkBig},
		Credential:    &zkbanw.Credential{Signature: signer.Credential},
		Period:        periodBig,
	}

	return &signerStruct, nil
}

func (signer *HighLevelSigner) FromSigner(signerObj *zkbanw.Signer) {
	if signer.Credential != nil {
		signer.Credential = signerObj.Credential.Signature
	}

	signer.Secret = signerObj.UserSecretKey.Number.Bytes()
	signer.UserPublicKey = signerObj.UserPublicKey.Number.Bytes()
	signer.Period = signerObj.Period.Int64()
}

func Sign(m []byte, counter int64, signer HighLevelSigner, gpk []byte, prover *HighLevelSnarkProver) (*Signature, error) {
	snarkProver, err := prover.ToSnarkProver()
	if err != nil {
		return nil, err
	}

	mBig := big.NewInt(0).SetBytes(m)
	counterBig := big.NewInt(counter)

	signerObj, err := signer.ToSigner()
	if err != nil {
		return nil, err
	}

	gpkObj, err := witness.GroupPublicKeyFromBytes(gpk)
	if err != nil {
		return nil, err
	}

	proof, assign, err := core.Sign(mBig, counterBig, signerObj, gpkObj, snarkProver)
	if err != nil {
		return nil, err
	}

	proofBytes, err := snark.EncodeProof(proof)
	pi := big.NewInt(0).SetBytes(proofBytes)
	nym := assign.Nym.(*big.Int)
	signature := assign.Signature.(*big.Int)

	result := Signature{
		Proof: pi.Bytes(),
		Nym:   nym.Bytes(),
		Sigma: signature.Bytes(),
	}

	return &result, err
}

func Verify(signature *Signature, m []byte, counter, period int64, gpk, verifyKeyBytes []byte) error {
	proofObj, err := snark.DecodeProof(signature.Proof)
	if err != nil {
		return err
	}

	verifyKeyObj, err := snark.DecodeVerifierKey(verifyKeyBytes)
	if err != nil {
		return err
	}

	mBig := big.NewInt(0).SetBytes(m)
	counterBig := big.NewInt(counter)
	periodBig := big.NewInt(period)

	publicKey := eddsa.PublicKey{}
	publicKey.Assign(snark.TwistededwardsCurve, gpk)

	dummy := eddsa.Signature{}
	dummyBuf := [32]byte{}
	dummy.Assign(snark.TwistededwardsCurve, dummyBuf[:])

	assign := circuit.SignCircuit{
		UserSecretKey: 0,
		CredentialAuthInfo: circuit.CredentialAuthInfo{
			Credential: dummy,
			Period:     period,
		},
		GroupPublicKey: publicKey,
		SessionTag:     witness.SessionTag(counterBig, periodBig),
		Message:        mBig,
		Signature:      signature.Sigma,
		Nym:            signature.Nym,
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
	if err != nil {
		return err
	}

	return nil
}
