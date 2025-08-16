# Test

## Benchmark

### How to work

```sh
go run ./latency/baseline
go run ./latency/update

go run ./storage/baseline
go run ./storage/update

# go test --bench ^BenchmarkAll$ . -timeout 0 -benchtime 5x
# go test --bench ^BenchmarkUpdate$ . -timeout 0 -benchtime 5x
```
