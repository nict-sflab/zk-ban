# Benchmark

### How to work
#### Latency

```sh
go run ./latency/baseline
go run ./latency/scale
go run -tags purego ./latency/compare 

# if want to use go test
# go test --bench ^BenchmarkAll$ . -timeout 0 -benchtime 5x
# go test --bench ^BenchmarkUpdate$ . -timeout 0 -benchtime 5x
```

#### Storage 

```
go run ./storage/baseline
go run ./storage/scale
```
