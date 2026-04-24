# Benchmark

### How to work

```sh
# Latency

go run ./latency/baseline
go run ./latency/scale
go run ./latency/compare 

# if want to use go test
# go test --bench ^BenchmarkBaselineRun$ . -timeout 0 -benchtime 5x
# go test --bench ^BenchmarkScalabilityRun$$ . -timeout 0 -benchtime 5x

# Storage 

go run ./storage/baseline
go run ./storage/scale
go run ./storage/compare
```
