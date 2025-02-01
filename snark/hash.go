package snark

import (
	"github.com/akakou/zk-ban/primitives/hash/circuit"
	"github.com/akakou/zk-ban/primitives/hash/commit"
	"github.com/consensys/gnark-crypto/hash"
)

var NewCircuitHash = circuit.NewMIMC
var CircuitHash = circuit.HashMaker(NewCircuitHash)

var NewCommitHashBase = hash.MIMC_BLS12_381
var NewCommitHash = NewCommitHashBase.New()
var CommitHash = commit.MimcHash(NewCommitHashBase)
