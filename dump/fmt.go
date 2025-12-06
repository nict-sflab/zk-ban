package dump

import "fmt"

var ProverKeyFileNameFormat = "%s_prover-%s.bin"
var CircuitFileNameFormat = "%s_circuit-%s.bin"
var VerifierKeyFileNameFormat = "%s_verifier-%s.bin"
var MetaFileNameFormat = "%s_verifier-%s.meta.json"

var KeyPath = "./"

func FileName(name, protocol string, format string) string {
	return fmt.Sprintf(format, protocol, name)
}
