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
BenchmarkUpdate/Prove_,16200,270,60,constant-16                        5        1053913112 ns/op
BenchmarkUpdate/Verify_,16200,270,60,constant-16                       5           4424008 ns/op
BenchmarkUpdate/Prove_,32400,270,120,constant-16                       5        1095179752 ns/op
BenchmarkUpdate/Verify_,32400,270,120,constant-16                      5           7725209 ns/op
BenchmarkUpdate/Prove_,48600,270,180,constant-16                       5        1139067851 ns/op
BenchmarkUpdate/Verify_,48600,270,180,constant-16                      5           9514676 ns/op
BenchmarkUpdate/Prove_,64800,270,240,constant-16                       5        1123956167 ns/op
BenchmarkUpdate/Verify_,64800,270,240,constant-16                      5           9063854 ns/op
BenchmarkUpdate/Prove_,81000,270,300,constant-16                       5        1535963438 ns/op
BenchmarkUpdate/Verify_,81000,270,300,constant-16                      5           9812829 ns/op
BenchmarkUpdate/Prove_,97200,270,360,constant-16                       5        1672649830 ns/op
BenchmarkUpdate/Verify_,97200,270,360,constant-16                      5           9439552 ns/op
BenchmarkUpdate/Prove_,113400,270,420,constant-16                      5        1709462225 ns/op
BenchmarkUpdate/Verify_,113400,270,420,constant-16                     5          14805405 ns/op
BenchmarkUpdate/Prove_,129600,270,480,constant-16                      5        1699893991 ns/op
BenchmarkUpdate/Verify_,129600,270,480,constant-16                     5          14399208 ns/op
BenchmarkUpdate/Prove_,145800,270,540,constant-16                      5        1750627316 ns/op
BenchmarkUpdate/Verify_,145800,270,540,constant-16                     5          14729637 ns/op
BenchmarkUpdate/Prove_,162000,270,600,constant-16                      5        1778040382 ns/op
BenchmarkUpdate/Verify_,162000,270,600,constant-16                     5          14758270 ns/op
BenchmarkUpdate/Prove_,16200,270,-,linear-16                           5        1015203273 ns/op
BenchmarkUpdate/Verify_,16200,270,-,linear-16                          5           2430192 ns/op
BenchmarkUpdate/Prove_,32400,270,-,linear-16                           5        1014011648 ns/op
BenchmarkUpdate/Verify_,32400,270,-,linear-16                          5           2324938 ns/op
BenchmarkUpdate/Prove_,48600,270,-,linear-16                           5        1088111470 ns/op
BenchmarkUpdate/Verify_,48600,270,-,linear-16                          5           6148534 ns/op
BenchmarkUpdate/Prove_,64800,270,-,linear-16                           5        1090665174 ns/op
BenchmarkUpdate/Verify_,64800,270,-,linear-16                          5           5670229 ns/op
BenchmarkUpdate/Prove_,81000,270,-,linear-16                           5        1153071868 ns/op
BenchmarkUpdate/Verify_,81000,270,-,linear-16                          5           9948261 ns/op
BenchmarkUpdate/Prove_,97200,270,-,linear-16                           5        1166819655 ns/op
BenchmarkUpdate/Verify_,97200,270,-,linear-16                          5          10715371 ns/op
BenchmarkUpdate/Prove_,113400,270,-,linear-16                          5        1696499699 ns/op
BenchmarkUpdate/Verify_,113400,270,-,linear-16                         5          14077767 ns/op
BenchmarkUpdate/Prove_,129600,270,-,linear-16                          5        1685452902 ns/op
BenchmarkUpdate/Verify_,129600,270,-,linear-16                         5          13021846 ns/op
BenchmarkUpdate/Prove_,145800,270,-,linear-16                          5        1742488962 ns/op
BenchmarkUpdate/Verify_,145800,270,-,linear-16                         5          15203998 ns/op
BenchmarkUpdate/Prove_,162000,270,-,linear-16                          5        1760661555 ns/op
BenchmarkUpdate/Verify_,162000,270,-,linear-16                         5          13653875 ns/op
BenchmarkUpdate/Prove_,16200,1,16200,constant-16                       5         148729102 ns/op
BenchmarkUpdate/Verify_,16200,1,16200,constant-16                      5           5429818 ns/op
BenchmarkUpdate/Prove_,32400,1,32400,constant-16                       5         256577273 ns/op
BenchmarkUpdate/Verify_,32400,1,32400,constant-16                      5           7184420 ns/op
BenchmarkUpdate/Prove_,48600,1,48600,constant-16                       5         294336801 ns/op
BenchmarkUpdate/Verify_,48600,1,48600,constant-16                      5           9217225 ns/op
BenchmarkUpdate/Prove_,64800,1,64800,constant-16                       5         402929840 ns/op
BenchmarkUpdate/Verify_,64800,1,64800,constant-16                      5           9981917 ns/op
BenchmarkUpdate/Prove_,81000,1,81000,constant-16                       5         442969859 ns/op
BenchmarkUpdate/Verify_,81000,1,81000,constant-16                      5           9405126 ns/op
BenchmarkUpdate/Prove_,97200,1,97200,constant-16                       5         521369047 ns/op
BenchmarkUpdate/Verify_,97200,1,97200,constant-16                      5          10413674 ns/op
BenchmarkUpdate/Prove_,113400,1,113400,constant-16                     5         519325212 ns/op
BenchmarkUpdate/Verify_,113400,1,113400,constant-16                    5          16054386 ns/op
BenchmarkUpdate/Prove_,129600,1,129600,constant-16                     5         752821171 ns/op
BenchmarkUpdate/Verify_,129600,1,129600,constant-16                    5          15091101 ns/op
BenchmarkUpdate/Prove_,145800,1,145800,constant-16                     5         770833601 ns/op
BenchmarkUpdate/Verify_,145800,1,145800,constant-16                    5          15529495 ns/op
BenchmarkUpdate/Prove_,162000,1,162000,constant-16                     5         807589782 ns/op
BenchmarkUpdate/Verify_,162000,1,162000,constant-16                    5          15826650 ns/op
BenchmarkUpdate/Prove_,34980,60,583,constant-16                        5         472367047 ns/op
BenchmarkUpdate/Verify_,34980,60,583,constant-16                       5           6623272 ns/op
BenchmarkUpdate/Prove_,34920,120,291,constant-16                       5         625271495 ns/op
BenchmarkUpdate/Verify_,34920,120,291,constant-16                      5           6616429 ns/op
BenchmarkUpdate/Prove_,34920,180,194,constant-16                       5         922089043 ns/op
BenchmarkUpdate/Verify_,34920,180,194,constant-16                      5           6520661 ns/op
BenchmarkUpdate/Prove_,34800,240,145,constant-16                       5        1022137959 ns/op
BenchmarkUpdate/Verify_,34800,240,145,constant-16                      5           6658356 ns/op
BenchmarkUpdate/Prove_,34800,300,116,constant-16                       5        1131740829 ns/op
BenchmarkUpdate/Verify_,34800,300,116,constant-16                      5           6170932 ns/op
BenchmarkUpdate/Prove_,34920,360,97,constant-16                        5        1599844093 ns/op
BenchmarkUpdate/Verify_,34920,360,97,constant-16                       5           7040791 ns/op
BenchmarkUpdate/Prove_,34860,420,83,constant-16                        5        1812891330 ns/op
BenchmarkUpdate/Verify_,34860,420,83,constant-16                       5           7004528 ns/op
BenchmarkUpdate/Prove_,34560,480,72,constant-16                        5        1946339055 ns/op
BenchmarkUpdate/Verify_,34560,480,72,constant-16                       5           7200387 ns/op
BenchmarkUpdate/Prove_,34560,540,64,constant-16                        5        2044743740 ns/op
BenchmarkUpdate/Verify_,34560,540,64,constant-16                       5           6469047 ns/op
BenchmarkUpdate/Prove_,34800,600,58,constant-16                        5        2212724350 ns/op
BenchmarkUpdate/Verify_,34800,600,58,constant-16                       5           6328057 ns/op
BenchmarkUpdate/Prove_,35000,60,-,linear-16                            5         463534978 ns/op
BenchmarkUpdate/Verify_,35000,60,-,linear-16                           5           5995842 ns/op
BenchmarkUpdate/Prove_,35000,120,-,linear-16                           5         600696352 ns/op
BenchmarkUpdate/Verify_,35000,120,-,linear-16                          5           5449238 ns/op
BenchmarkUpdate/Prove_,35000,180,-,linear-16                           5         927472874 ns/op
BenchmarkUpdate/Verify_,35000,180,-,linear-16                          5           5812736 ns/op
BenchmarkUpdate/Prove_,35000,240,-,linear-16                           5        1031482862 ns/op
BenchmarkUpdate/Verify_,35000,240,-,linear-16                          5           6029613 ns/op
BenchmarkUpdate/Prove_,35000,300,-,linear-16                           5        1059385365 ns/op
BenchmarkUpdate/Verify_,35000,300,-,linear-16                          5           2472453 ns/op
BenchmarkUpdate/Prove_,35000,360,-,linear-16                           5        1192501480 ns/op
BenchmarkUpdate/Verify_,35000,360,-,linear-16                          5           2493873 ns/op
BenchmarkUpdate/Prove_,35000,420,-,linear-16                           5        1650388527 ns/op
BenchmarkUpdate/Verify_,35000,420,-,linear-16                          5           2338120 ns/op
BenchmarkUpdate/Prove_,35000,480,-,linear-16                           5        1844940700 ns/op
BenchmarkUpdate/Verify_,35000,480,-,linear-16                          5           2441037 ns/op
BenchmarkUpdate/Prove_,35000,540,-,linear-16                           5        1993329031 ns/op
BenchmarkUpdate/Verify_,35000,540,-,linear-16                          5           2470075 ns/op
BenchmarkUpdate/Prove_,35000,600,-,linear-16                           5        2112257603 ns/op
BenchmarkUpdate/Verify_,35000,600,-,linear-16                          5           2349911 ns/op
PASS
ok      github.com/akakou/zk-ban/test   1192.821s
```