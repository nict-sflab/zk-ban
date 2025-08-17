package storage

import (
	"bytes"
	"io"

	"github.com/akakou/zk-ban/snark"
	"github.com/akakou/zk-ban/test"
)

type Result = map[string]map[string]map[string]int

func StoreSize(name, parent, root string, size int, result Result) {
	if _, ok := result[root]; !ok {
		result[root] = make(map[string]map[string]int)
	}
	if _, ok := result[root][parent]; !ok {
		result[root][parent] = make(map[string]int)
	}

	result[root][parent][name] = size
}
func StoreBufSize(name, parent, root string, buf []byte, result Result) {
	StoreSize(name, parent, root, len(buf), result)
}

func StoreWritableSize(name, parent, root string, writer io.WriterTo, result Result) {
	var buf bytes.Buffer
	_, err := writer.WriteTo(&buf)
	test.PanicIfErr(err)

	StoreBufSize(name, parent, root, buf.Bytes(), result)
}

func StoreCircuitObjectSize(parent, root string, p *snark.SnarkParams, result Result) {
	StoreWritableSize("cs", parent, root, p.ConstraintSystem, result)
	StoreWritableSize("pk", parent, root, p.ProveKey, result)
	StoreWritableSize("vk", parent, root, p.VerifyKey, result)
}
