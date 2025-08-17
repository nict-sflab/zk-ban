# Test

## Benchmark

### How to work

```sh
go run ./latency/baseline
go run ./latency/scale
go run -tags purego ./latency/compare 

go run ./storage/baseline
go run ./storage/scale

# go test --bench ^BenchmarkAll$ . -timeout 0 -benchtime 5x
# go test --bench ^BenchmarkUpdate$ . -timeout 0 -benchtime 5x
```
