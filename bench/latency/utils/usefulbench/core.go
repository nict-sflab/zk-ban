package usefulbench

type Benchmarker interface {
	Loop() bool
	Run(string, func(b Benchmarker)) bool
}
