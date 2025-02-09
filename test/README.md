# Test

## Benchmark

### How to work

```sh
go test --bench BenchmarkAll . -timeout 0 -benchtime 5x
go test --bench BenchmarkUpdate . -timeout 0 -benchtime 5x
```

### Our result

Environments:
- CPU: Intel Xeon Gold 5118 @ 16x 2.295GHz
- RAM: 1379MiB / 32150MiB
- OS: Ubuntu 22.04 jammy

Result:
```
akakou@ra-webs:~/zk-ban/test$ go test --bench BenchmarkAll . -timeout 0 -benchtime 5x
goos: linux
goarch: amd64
pkg: github.com/akakou/zk-ban/test
cpu: Intel(R) Xeon(R) Gold 5118 CPU @ 2.30GHz
BenchmarkAll/join_req-16                       5          14418801 ns/op
BenchmarkAll/verify_join_req-16                5           2088670 ns/op
BenchmarkAll/sign-16                           5          75351809 ns/op
BenchmarkAll/verify-16                         5           2103861 ns/op
BenchmarkAll/update-req-16                     5        1077088720 ns/op
BenchmarkAll/update-verify-16                  5           5847203 ns/op
PASS
ok      github.com/akakou/zk-ban/test   39.990s

akakou@ra-webs:~/zk-ban/test$ go test --bench BenchmarkUpdate . -timeout 0 -benchtime 5x
goos: linux
goarch: amd64
pkg: github.com/akakou/zk-ban/test
cpu: Intel(R) Xeon(R) Gold 5118 CPU @ 2.30GHz
BenchmarkUpdate/Prove_,8100,30,270,constant-16                 5        1057832089 ns/op
BenchmarkUpdate/Verify_,8100,30,270,constant-16                5           4153520 ns/op
BenchmarkUpdate/Prove_,16200,60,270,constant-16                5        1061309241 ns/op
BenchmarkUpdate/Verify_,16200,60,270,constant-16               5           4268631 ns/op
BenchmarkUpdate/Prove_,24300,90,270,constant-16                5        1039168731 ns/op
BenchmarkUpdate/Verify_,24300,90,270,constant-16               5           6570278 ns/op
BenchmarkUpdate/Prove_,32400,120,270,constant-16               5        1110044233 ns/op
BenchmarkUpdate/Verify_,32400,120,270,constant-16              5           6093876 ns/op
BenchmarkUpdate/Prove_,40500,150,270,constant-16               5        1101068903 ns/op
BenchmarkUpdate/Verify_,40500,150,270,constant-16              5           5632847 ns/op
BenchmarkUpdate/Prove_,48600,180,270,constant-16               5        1103441025 ns/op
BenchmarkUpdate/Verify_,48600,180,270,constant-16              5           8921289 ns/op
BenchmarkUpdate/Prove_,56700,210,270,constant-16               5        1107650396 ns/op
BenchmarkUpdate/Verify_,56700,210,270,constant-16              5          10595995 ns/op
BenchmarkUpdate/Prove_,64800,240,270,constant-16               5        1126638620 ns/op
BenchmarkUpdate/Verify_,64800,240,270,constant-16              5          10524568 ns/op
BenchmarkUpdate/Prove_,72900,270,270,constant-16               5        1161237235 ns/op
BenchmarkUpdate/Verify_,72900,270,270,constant-16              5           9290615 ns/op
BenchmarkUpdate/Prove_,81000,300,270,constant-16               5        1545642754 ns/op
BenchmarkUpdate/Verify_,81000,300,270,constant-16              5           9948194 ns/op
BenchmarkUpdate/Prove_,270,8100,-,linear-16                    5        1028858991 ns/op
BenchmarkUpdate/Verify_,270,8100,-,linear-16                   5           2525095 ns/op
BenchmarkUpdate/Prove_,270,16200,-,linear-16                   5        1031774981 ns/op
BenchmarkUpdate/Verify_,270,16200,-,linear-16                  5           2667736 ns/op
BenchmarkUpdate/Prove_,270,24300,-,linear-16                   5        1019505529 ns/op
BenchmarkUpdate/Verify_,270,24300,-,linear-16                  5           2373650 ns/op
BenchmarkUpdate/Prove_,270,32400,-,linear-16                   5        1006646364 ns/op
BenchmarkUpdate/Verify_,270,32400,-,linear-16                  5           2582486 ns/op
BenchmarkUpdate/Prove_,270,40500,-,linear-16                   5        1029322542 ns/op
BenchmarkUpdate/Verify_,270,40500,-,linear-16                  5           2619356 ns/op
BenchmarkUpdate/Prove_,270,48600,-,linear-16                   5        1014944719 ns/op
BenchmarkUpdate/Verify_,270,48600,-,linear-16                  5           2424135 ns/op
BenchmarkUpdate/Prove_,270,56700,-,linear-16                   5        1004979633 ns/op
BenchmarkUpdate/Verify_,270,56700,-,linear-16                  5           2489899 ns/op
BenchmarkUpdate/Prove_,270,64800,-,linear-16                   5        1035282316 ns/op
BenchmarkUpdate/Verify_,270,64800,-,linear-16                  5           2628096 ns/op
BenchmarkUpdate/Prove_,270,72900,-,linear-16                   5        1026691667 ns/op
BenchmarkUpdate/Verify_,270,72900,-,linear-16                  5           2412101 ns/op
BenchmarkUpdate/Prove_,270,81000,-,linear-16                   5        1013893173 ns/op
BenchmarkUpdate/Verify_,270,81000,-,linear-16                  5           2354352 ns/op
BenchmarkUpdate/Prove_,8100,8100,1,constant-16                 5         127937934 ns/op
BenchmarkUpdate/Verify_,8100,8100,1,constant-16                5           3899478 ns/op
BenchmarkUpdate/Prove_,16200,16200,1,constant-16               5         154322584 ns/op
BenchmarkUpdate/Verify_,16200,16200,1,constant-16              5           5399007 ns/op
BenchmarkUpdate/Prove_,24300,24300,1,constant-16               5         224092388 ns/op
BenchmarkUpdate/Verify_,24300,24300,1,constant-16              5           5841594 ns/op
BenchmarkUpdate/Prove_,32400,32400,1,constant-16               5         238504461 ns/op
BenchmarkUpdate/Verify_,32400,32400,1,constant-16              5           8259045 ns/op
BenchmarkUpdate/Prove_,40500,40500,1,constant-16               5         266083367 ns/op
BenchmarkUpdate/Verify_,40500,40500,1,constant-16              5           6935951 ns/op
BenchmarkUpdate/Prove_,48600,48600,1,constant-16               5         292768163 ns/op
BenchmarkUpdate/Verify_,48600,48600,1,constant-16              5           8781235 ns/op
BenchmarkUpdate/Prove_,56700,56700,1,constant-16               5         400250604 ns/op
BenchmarkUpdate/Verify_,56700,56700,1,constant-16              5           9505635 ns/op
BenchmarkUpdate/Prove_,64800,64800,1,constant-16               5         408393004 ns/op
BenchmarkUpdate/Verify_,64800,64800,1,constant-16              5           8774633 ns/op
BenchmarkUpdate/Prove_,72900,72900,1,constant-16               5         428038107 ns/op
BenchmarkUpdate/Verify_,72900,72900,1,constant-16              5           8999437 ns/op
BenchmarkUpdate/Prove_,81000,81000,1,constant-16               5         446548797 ns/op
BenchmarkUpdate/Verify_,81000,81000,1,constant-16              5           9802422 ns/op
BenchmarkUpdate/Prove_,34980,583,60,constant-16                5         449946964 ns/op
BenchmarkUpdate/Verify_,34980,583,60,constant-16               5           6087098 ns/op
BenchmarkUpdate/Prove_,34920,291,120,constant-16               5         602074802 ns/op
BenchmarkUpdate/Verify_,34920,291,120,constant-16              5           6873341 ns/op
BenchmarkUpdate/Prove_,34920,194,180,constant-16               5         935132609 ns/op
BenchmarkUpdate/Verify_,34920,194,180,constant-16              5           6402739 ns/op
BenchmarkUpdate/Prove_,34800,145,240,constant-16               5        1045688853 ns/op
BenchmarkUpdate/Verify_,34800,145,240,constant-16              5           5809427 ns/op
BenchmarkUpdate/Prove_,34800,116,300,constant-16               5        1139758820 ns/op
BenchmarkUpdate/Verify_,34800,116,300,constant-16              5           6660059 ns/op
BenchmarkUpdate/Prove_,34920,97,360,constant-16                5        1606800809 ns/op
BenchmarkUpdate/Verify_,34920,97,360,constant-16               5           6827919 ns/op
BenchmarkUpdate/Prove_,34860,83,420,constant-16                5        1841766324 ns/op
BenchmarkUpdate/Verify_,34860,83,420,constant-16               5           6216134 ns/op
BenchmarkUpdate/Prove_,34560,72,480,constant-16                5        1942777841 ns/op
BenchmarkUpdate/Verify_,34560,72,480,constant-16               5           5934215 ns/op
BenchmarkUpdate/Prove_,34560,64,540,constant-16                5        2051235662 ns/op
BenchmarkUpdate/Verify_,34560,64,540,constant-16               5           6481385 ns/op
BenchmarkUpdate/Prove_,34800,58,600,constant-16                5        2187075100 ns/op
BenchmarkUpdate/Verify_,34800,58,600,constant-16               5           6573966 ns/op
BenchmarkUpdate/Prove_,60,35000,-,linear-16                    5         338769653 ns/op
BenchmarkUpdate/Verify_,60,35000,-,linear-16                   5           3979915 ns/op
BenchmarkUpdate/Prove_,120,35000,-,linear-16                   5         549889008 ns/op
BenchmarkUpdate/Verify_,120,35000,-,linear-16                  5           5012249 ns/op
BenchmarkUpdate/Prove_,180,35000,-,linear-16                   5         639330490 ns/op
BenchmarkUpdate/Verify_,180,35000,-,linear-16                  5           2880695 ns/op
BenchmarkUpdate/Prove_,240,35000,-,linear-16                   5         964177380 ns/op
BenchmarkUpdate/Verify_,240,35000,-,linear-16                  5           2332662 ns/op
BenchmarkUpdate/Prove_,300,35000,-,linear-16                   5        1062079294 ns/op
BenchmarkUpdate/Verify_,300,35000,-,linear-16                  5           2412848 ns/op
BenchmarkUpdate/Prove_,360,35000,-,linear-16                   5        1192028571 ns/op
BenchmarkUpdate/Verify_,360,35000,-,linear-16                  5           2461066 ns/op
BenchmarkUpdate/Prove_,420,35000,-,linear-16                   5        1679182000 ns/op
BenchmarkUpdate/Verify_,420,35000,-,linear-16                  5           2505912 ns/op
BenchmarkUpdate/Prove_,480,35000,-,linear-16                   5        1844927301 ns/op
BenchmarkUpdate/Verify_,480,35000,-,linear-16                  5           2935508 ns/op
BenchmarkUpdate/Prove_,540,35000,-,linear-16                   5        1968149159 ns/op
BenchmarkUpdate/Verify_,540,35000,-,linear-16                  5           2850361 ns/op
BenchmarkUpdate/Prove_,600,35000,-,linear-16                   5        2120412885 ns/op
BenchmarkUpdate/Verify_,600,35000,-,linear-16                  5           2497784 ns/op
PASS
ok      github.com/akakou/zk-ban/test   980.478s
```