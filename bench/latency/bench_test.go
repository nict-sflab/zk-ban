package latency

import (
	"testing"

	"github.com/akakou/zk-ban/bench/latency/utils/usefulbench"
)

func BenchmarkBaselineRun(b *testing.B) {
	BenchmarkBaseline(&usefulbench.StandardBenchMarker{B: b})

}

func BenchmarkUpdateRun(b *testing.B) {
	BenchmarkUpdate(&usefulbench.StandardBenchMarker{B: b})
}
