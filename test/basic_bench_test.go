package test_test

import (
	"testing"

	"github.com/akakou/zk-ban/test"
)

func BenchmarkAll(b *testing.B) {
	test.BenchmarkAll(b)
}
