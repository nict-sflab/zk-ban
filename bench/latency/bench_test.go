package latency

import (
	"testing"

	"github.com/akakou/zk-ban/bench/latency/utils/usefulbench"
)

func BenchmarkBaselineRun(b *testing.B) {
	BenchmarkBaseline(&usefulbench.StandardBenchMarker{B: b})

}

func BenchmarkScalabilityRun(b *testing.B) {
	BenchmarkScalability(&usefulbench.StandardBenchMarker{B: b})
}

func BenchmarkCompareRun(b *testing.B) {
	BenchmarkCompare(&usefulbench.StandardBenchMarker{B: b})
}

func BenchmarkLoadRun(b *testing.B) {
	BenchmarkLoad(&usefulbench.StandardBenchMarker{B: b})
}
