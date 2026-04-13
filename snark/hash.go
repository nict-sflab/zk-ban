package snark

import (
	"github.com/akakou/zk-ban/primitives/hash/circuit"
	"github.com/akakou/zk-ban/primitives/hash/witness"
	"github.com/consensys/gnark-crypto/hash"
)

// var NewCircuitHash = circuit.NewPoseidon
// var NewCommitHash = poseidon.NewPoseidon
// var CommitHash = commit.PoseidonHash

var NewCircuitHash = circuit.NewPoseidon2
var NewCommitHash = NewCommitHashBase.New
var CommitHash = witness.Hash(NewCommitHashBase)

var NewCommitHashBase = hash.POSEIDON2_BLS12_381

var CircuitHash = circuit.HashMaker(NewCircuitHash)
