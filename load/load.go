package load

import (
	"fmt"
	"io/fs"
	"syscall"
	"time"

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
	fmt.Printf("\nread file: %s\n", name)
	now := time.Now()
	startCPU := getCPUTimeNS()

	buf, err := fs.ReadFile(f, name)
	if err != nil {
		return nil, err
	}
	latency := time.Since(now)
	elapsedCPU := getCPUTimeNS() - startCPU
	fmt.Printf("%v (%v)\n", latency, time.Nanosecond*time.Duration(elapsedCPU))

	key, err := decoder(buf)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func getCPUTimeNS() int64 {
	var ru syscall.Rusage
	_ = syscall.Getrusage(syscall.RUSAGE_SELF, &ru)
	user := int64(ru.Utime.Sec)*1e9 + int64(ru.Utime.Usec)*1e3
	sys := int64(ru.Stime.Sec)*1e9 + int64(ru.Stime.Usec)*1e3
	return user + sys
}
