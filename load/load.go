package load

import (
	"io/fs"

	"github.com/akakou/zk-ban/snark"
)

func LoadUserUpdateKey(readDir fs.ReadDirFS, readFile fs.ReadFileFS) ([]*snark.SnarkProver, error) {
	return LoadUpdateKeys(readDir, readFile, DocodeProver)
}

func LoadGroupManagerUpdateKey(readDir fs.ReadDirFS, readFile fs.ReadFileFS) ([]*snark.SizedSnarkVerifier, error) {
	return LoadUpdateKeys(readDir, readFile, DocodeSizedVerifyingKey)
}

func LoadUpdateKeys[T any](readDir fs.ReadDirFS, readFile fs.ReadFileFS, decoder func([]byte) (T, error)) ([]T, error) {
	keys := []T{}
	files, err := readDir.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		name := f.Name()

		buf, err := readFile.ReadFile(name)
		if err != nil {
			return nil, err
		}

		key, err := decoder(buf)
		if err != nil {
			return nil, err
		}

		keys = append(keys, key)
	}

	return keys, nil

}
