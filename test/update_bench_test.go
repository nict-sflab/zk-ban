package test_test

import (
	"testing"

	"github.com/akakou/zk-ban/test"
)

func BenchmarkUpdate(b *testing.B) {
	test.BenchmarkUpdate(b)
}
