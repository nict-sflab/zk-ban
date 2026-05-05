package dump

import "fmt"

var ProverKeyFileNameFormat = "%s%s_prover-%s.bin"
var CircuitFileNameFormat = "%s%s_circuit-%s.bin"
var VerifierKeyFileNameFormat = "%s%s_verifier-%s.bin"
var MetaFileNameFormat = "%s%s_verifier-%s.meta.json"

var KeyPath = "./"

func FileName(name, protocol string, format string, option *string) string {
	if option != nil {
		*option = *option + "_"
	} else {
		tmp := ""
		option = &tmp
	}
	return KeyPath + fmt.Sprintf(format, *option, protocol, name)
}
