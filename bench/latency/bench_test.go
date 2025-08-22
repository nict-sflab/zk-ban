package latency

import (
	"testing"

	"github.com/akakou/zk-ban/bench/latency/utils/usefulbench"
)

func BenchmarkBaselineRun(b *testing.B) {
	BenchmarkBaseline(&usefulbench.StandardBenchMarker{B: b})

}

func BenchmarkUpdateRun(b *testing.B) {
	BenchmarkScalability(&usefulbench.StandardBenchMarker{B: b})
}

func BenchmarkCompareRun(b *testing.B) {
	BenchmarkCompare(&usefulbench.StandardBenchMarker{B: b})
}
