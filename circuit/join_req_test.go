package circuit

// func TestJoinReq(t *testing.T) {
// 	assert := test.NewAssert(t)
// 	var joinCercuit JoinRequestCircuit

// 	usk := witness.UserSecretKey{primitives.NewBigInt(1)}
// 	period := int64(2024)

// 	upk, err := usk.PublicKey(period)
// 	assert.NoError(err)

// 	assign := NewJoinRequestWitness(
// 		period,
// 		upk,
// 		&usk,
// 	)

// 	assert.ProverSucceeded(&joinCercuit, assign, test.WithCurves(snark.EcCurve))
// }
