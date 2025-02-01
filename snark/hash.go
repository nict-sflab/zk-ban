package snark

import (
	"github.com/akakou/zk-ban/primitives/hash/circuit"
	"github.com/akakou/zk-ban/primitives/hash/commit"
	"github.com/consensys/gnark-crypto/hash"
)

var CircuitHash = circuit.HashMaker(circuit.NewMIMC)
var CommitHash = commit.MimcHash(hash.MIMC_BLS12_381)

var Hasher = hash.MIMC_BLS12_381.New()
