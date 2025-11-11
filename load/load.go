package load

import (
	"io/fs"

	"github.com/akakou/zk-ban/snark"
)

func LoadUserUpdateKey(fs fs.ReadDirFS) ([]*snark.SnarkProver, error) {
	return LoadUpdateKeys(fs, DocodeProver)
}

func LoadGroupManagerUpdateKey(fs fs.ReadDirFS) ([]*snark.SizedSnarkVerifier, error) {
	return LoadUpdateKeys(fs, DocodeSizedVerifyingKey)
}

func LoadUpdateKeys[T any](fs fs.ReadDirFS, decoder func([]byte) (T, error)) ([]T, error) {
	keys := []T{}
	files, err := fs.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		name := f.Name()

		key, err := LoadUpdateKey(name, fs, decoder)
		if err != nil {
			return nil, err
		}

		keys = append(keys, *key)
	}

	return keys, nil

}

func LoadUpdateKey[T any](name string, f fs.FS, decoder func([]byte) (T, error)) (*T, error) {
	buf, err := fs.ReadFile(f, name)
	if err != nil {
		return nil, err
	}

	key, err := decoder(buf)
	if err != nil {
		return nil, err
	}

	return &key, nil
}
