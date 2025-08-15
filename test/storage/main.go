package storage

import (
	"bytes"
	"io"

	"github.com/akakou/zk-ban/snark"
	"github.com/akakou/zk-ban/test"
)

type Result = map[string]map[string]int

func StoreSize(name, parent string, buf []byte, result Result) {
	result[parent][name] = len(buf)
}

func StoreWritableSize(name, parent string, writer io.WriterTo, result Result) {
	var buf bytes.Buffer
	_, err := writer.WriteTo(&buf)
	test.PanicIfErr(err)

	StoreSize(name, parent, buf.Bytes(), result)
}

func StoreCircuitObjectSize(name, parent string, p *snark.SnarkParams, result Result) {
	StoreWritableSize(name+"-cs", parent, p.ConstraintSystem, result)
	StoreWritableSize(name+"-pk", parent, p.ProveKey, result)
	StoreWritableSize(name+"-vk", parent, p.VerifyKey, result)
}
