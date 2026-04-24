package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/akakou/zk-ban/dump"
	"github.com/akakou/zk-ban/test"
	"github.com/akakou/zk-ban/witness"
)

func BenchUpdate(rlSize witness.RevocationListSize, root, itemName, fileName string, result Result) {
	test.PrepareUpdateKeyCached(rlSize, fileName)

	cs, err := os.Stat(dump.FileName(fileName, "update", dump.CircuitFileNameFormat))
	test.PanicIfErr(err)
	pk, err := os.Stat(dump.FileName(fileName, "update", dump.ProverKeyFileNameFormat))
	test.PanicIfErr(err)
	vk, err := os.Stat(dump.FileName(fileName, "update", dump.VerifierKeyFileNameFormat))
	test.PanicIfErr(err)

	StoreSize(itemName, "pk", root, int(pk.Size()), result)
	StoreSize(itemName, "cs", root, int(cs.Size()), result)
	StoreSize(itemName, "vk", root, int(vk.Size()), result)

	j, err := json.Marshal(result)
	test.PanicIfErr(err)
	fmt.Printf("result %s\n", j)
}
