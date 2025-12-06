package main

import (
	"github.com/akakou/zk-ban/dump"
	"github.com/akakou/zk-ban/witness"
)

func main() {
	dump.KeyPath = "../../../zk-ban-system/keys/"

	dump.DumpBasicKeys()
	// dump.DumpUpdateKeys(260, 540)
	// dump.DumpUpdateKeys(260, 270)
	// dump.DumpUpdateKeys(130, 540)
	// dump.DumpUpdateKeys(130, 270)
	// dump.DumpdumpUpdateKeys(75, 270)
	// dump.DumpUpdateKeys(130, 135)
	// dump.DumpUpdateKeys(75, 135)
	// dump.DumpUpdateKeys(32, 135)
	rl := witness.MakeGaussianRLSize(12, 12, 4)
	dump.DumpUpdateKeys("gausse", rl)
}
