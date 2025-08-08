package circuit

// func TestSignCircuit(t *testing.T) {
// 	assert := test.NewAssert(t)

// 	usk := witness.UserSecretKey{primitives.NewBigInt(1)}
// 	m := big.NewInt(3)
// 	bsn := big.NewInt(4)
// 	period := int64(2024)

// 	gsk, gpk, err := witness.RandomGroupKeyPair()
// 	assert.NoError(err)

// 	upk, err := usk.PublicKey(period)
// 	assert.NoError(err)

// 	cert, err := gsk.IssueCredential(upk)
// 	assert.NoError(err)

// 	signer := witness.Signer{
// 		UserSecretKey: &usk,
// 		Credential:    cert,
// 		Period:        period,
// 	}

// 	commit, err := signer.CommitSign(&primitives.BigInt{*m}, &primitives.BigInt{*bsn})
// 	assert.NoError(err)

// 	authCircuit := SignCircuit{}

// 	witness, err := NewSignWitness(&primitives.BigInt{*m}, &primitives.BigInt{*bsn}, commit, &signer, gpk)
// 	assert.NoError(err)

// 	assert.ProverSucceeded(&authCircuit, witness, test.WithCurves(snark.EcCurve))
// }
