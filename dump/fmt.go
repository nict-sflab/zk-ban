package dump

import "fmt"

var UpdateProverKeyFileNameFormat = "update_prover-%s.bin"
var UpdateCircuitFileNameFormat = "update_circuit-%s.bin"
var UpdateVerifierKeyFileNameFormat = "update_verifier-%s.bin"
var UpdateVKMetaFileNameFormat = "update_verifier-%s.meta.json"

var KeyPath = "./"

func FileName(name string, format string) string {
	return fmt.Sprintf(format, name)
}
