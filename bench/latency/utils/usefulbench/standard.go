package usefulbench

import "testing"

type StandardBenchMarker struct {
	B *testing.B
}

func (b *StandardBenchMarker) Loop() bool {
	return b.B.Loop()
}

func (b *StandardBenchMarker) Run(tag string, target func(b Benchmarker)) bool {
	newTarget := func(bb *testing.B) {
		marker := &StandardBenchMarker{
			B: bb,
		}

		target(marker)
	}

	return b.B.Run(tag, newTarget)
}
