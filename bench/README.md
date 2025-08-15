# Test

## Benchmark

### How to work

```sh
go run ./bench/baseline
go run ./bench/update
# go test --bench ^BenchmarkAll$ . -timeout 0 -benchtime 5x
# go test --bench ^BenchmarkUpdate$ . -timeout 0 -benchtime 5x
```
