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
akakou@ra-webs:~/zk-ban/test$ go test -run . -bench BenchmarkAll -timeout 0 -benchtime 20x
goos: linux
goarch: amd64
pkg: github.com/akakou/zk-ban/test
cpu: Intel(R) Xeon(R) Gold 5118 CPU @ 2.30GHz
BenchmarkAll/join_req-16                      20          13599889 ns/op
BenchmarkAll/verify_join_req-16               20           2154118 ns/op
BenchmarkAll/sign-16                          20          74022378 ns/op
BenchmarkAll/verify-16                        20           2317964 ns/op
BenchmarkAll/update-req_(constant)-16                 20        1086758618 ns/op
BenchmarkAll/update-verify_(constant)-16              20           5903594 ns/op
BenchmarkAll/update-req_(linear)-16                   20        1008659901 ns/op
BenchmarkAll/update-verify_(linear)-16                20           2537904 ns/op
BenchmarkAll/update-req_(one_session)-16              20         244845114 ns/op
BenchmarkAll/update-verify_(one_session)-16           20           6906016 ns/op
PASS
ok      github.com/akakou/zk-ban/test   102.312s
```

```
akakou@ra-webs:~/zk-ban/test$ go test -run . -bench BenchmarkAll -timeout 0 -benchtime 20x
goos: linux
goarch: amd64
pkg: github.com/akakou/zk-ban/test
cpu: Intel(R) Xeon(R) Gold 5118 CPU @ 2.30GHz
BenchmarkUpdate/Prove_,162000,270,600,constant-16                     20        1804436504 ns/op
BenchmarkUpdate/Verify_,162000,270,600,constant-16                    20          16252460 ns/op
BenchmarkUpdate/Prove_,324000,270,1200,constant-16                    20        2079019135 ns/op
BenchmarkUpdate/Verify_,324000,270,1200,constant-16                   20          26507205 ns/op
BenchmarkUpdate/Prove_,486000,270,1800,constant-16                    20        2995796610 ns/op
BenchmarkUpdate/Verify_,486000,270,1800,constant-16                   20          48376743 ns/op
BenchmarkUpdate/Prove_,648000,270,2400,constant-16                    20        3215581149 ns/op
BenchmarkUpdate/Verify_,648000,270,2400,constant-16                   20          47617891 ns/op
BenchmarkUpdate/Prove_,810000,270,3000,constant-16                    20        3634036651 ns/op
BenchmarkUpdate/Verify_,810000,270,3000,constant-16                   20          48951192 ns/op
BenchmarkUpdate/Prove_,972000,270,3600,constant-16                    20        5194411128 ns/op
BenchmarkUpdate/Verify_,972000,270,3600,constant-16                   20          65403836 ns/op
BenchmarkUpdate/Prove_,1134000,270,4200,constant-16                   20        5364274254 ns/op
BenchmarkUpdate/Verify_,1134000,270,4200,constant-16                  20          64689012 ns/op
BenchmarkUpdate/Prove_,1296000,270,4800,constant-16                   20        5626837832 ns/op
BenchmarkUpdate/Verify_,1296000,270,4800,constant-16                  20          69988972 ns/op
BenchmarkUpdate/Prove_,1458000,270,5400,constant-16                   20        5820003878 ns/op
BenchmarkUpdate/Verify_,1458000,270,5400,constant-16                  20          66220425 ns/op
BenchmarkUpdate/Prove_,1620000,270,6000,constant-16                   20        6156559514 ns/op
BenchmarkUpdate/Verify_,1620000,270,6000,constant-16                  20          73179468 ns/op
BenchmarkUpdate/Prove_,162000,270,-,linear-16                         20        1750305490 ns/op
BenchmarkUpdate/Verify_,162000,270,-,linear-16                        20          14679464 ns/op
BenchmarkUpdate/Prove_,324000,270,-,linear-16                         20        1959222017 ns/op
BenchmarkUpdate/Verify_,324000,270,-,linear-16                        20          31928343 ns/op
BenchmarkUpdate/Prove_,486000,270,-,linear-16                         20        2987515886 ns/op
BenchmarkUpdate/Verify_,486000,270,-,linear-16                        20          45008778 ns/op
BenchmarkUpdate/Prove_,648000,270,-,linear-16                         20        3212873183 ns/op
BenchmarkUpdate/Verify_,648000,270,-,linear-16                        20          48078125 ns/op
BenchmarkUpdate/Prove_,810000,270,-,linear-16                         20        3619962998 ns/op
BenchmarkUpdate/Verify_,810000,270,-,linear-16                        20          43656778 ns/op
BenchmarkUpdate/Prove_,972000,270,-,linear-16                         20        5126154925 ns/op
BenchmarkUpdate/Verify_,972000,270,-,linear-16                        20          60803060 ns/op
BenchmarkUpdate/Prove_,1134000,270,-,linear-16                        20        5374041854 ns/op
BenchmarkUpdate/Verify_,1134000,270,-,linear-16                       20          64751340 ns/op
BenchmarkUpdate/Prove_,1296000,270,-,linear-16                        20        5661947581 ns/op
BenchmarkUpdate/Verify_,1296000,270,-,linear-16                       20          65703987 ns/op
BenchmarkUpdate/Prove_,1458000,270,-,linear-16                        20        5826111404 ns/op
BenchmarkUpdate/Verify_,1458000,270,-,linear-16                       20          68831633 ns/op
BenchmarkUpdate/Prove_,1620000,270,-,linear-16                        20        6158599088 ns/op
BenchmarkUpdate/Verify_,1620000,270,-,linear-16                       20          74999998 ns/op
BenchmarkUpdate/Prove_,162000,1,162000,constant-16                    20         812354426 ns/op
BenchmarkUpdate/Verify_,162000,1,162000,constant-16                   20          18395232 ns/op
BenchmarkUpdate/Prove_,324000,1,324000,constant-16                    20        1475704915 ns/op
BenchmarkUpdate/Verify_,324000,1,324000,constant-16                   20          29608958 ns/op
BenchmarkUpdate/Prove_,486000,1,486000,constant-16                    20        1834599479 ns/op
BenchmarkUpdate/Verify_,486000,1,486000,constant-16                   20          49430388 ns/op
BenchmarkUpdate/Prove_,648000,1,648000,constant-16                    20        2745537307 ns/op
BenchmarkUpdate/Verify_,648000,1,648000,constant-16                   20          50058343 ns/op
BenchmarkUpdate/Prove_,810000,1,810000,constant-16                    20        2940806608 ns/op
BenchmarkUpdate/Verify_,810000,1,810000,constant-16                   20          56463903 ns/op
BenchmarkUpdate/Prove_,972000,1,972000,constant-16                    20        3355653965 ns/op
BenchmarkUpdate/Verify_,972000,1,972000,constant-16                   20          63242317 ns/op
BenchmarkUpdate/Prove_,1134000,1,1134000,constant-16                  20        4969819716 ns/op
BenchmarkUpdate/Verify_,1134000,1,1134000,constant-16                 20          68644554 ns/op
BenchmarkUpdate/Prove_,1296000,1,1296000,constant-16                  20        5169861007 ns/op
BenchmarkUpdate/Verify_,1296000,1,1296000,constant-16                 20          69366226 ns/op
BenchmarkUpdate/Prove_,1458000,1,1458000,constant-16                  20        5388604201 ns/op
BenchmarkUpdate/Verify_,1458000,1,1458000,constant-16                 20          68186861 ns/op
BenchmarkUpdate/Prove_,1620000,1,1620000,constant-16                  20        5605384720 ns/op
BenchmarkUpdate/Verify_,1620000,1,1620000,constant-16                 20          72563572 ns/op
BenchmarkUpdate/Prove_,34800,600,58,constant-16                       20        2195543140 ns/op
BenchmarkUpdate/Verify_,34800,600,58,constant-16                      20           6367327 ns/op
BenchmarkUpdate/Prove_,34800,1200,29,constant-16                      20        4029006610 ns/op
BenchmarkUpdate/Verify_,34800,1200,29,constant-16                     20           6899477 ns/op
BenchmarkUpdate/Prove_,34200,1800,19,constant-16                      20        6455192579 ns/op
BenchmarkUpdate/Verify_,34200,1800,19,constant-16                     20           7014496 ns/op
BenchmarkUpdate/Prove_,33600,2400,14,constant-16                      20        8031715100 ns/op
BenchmarkUpdate/Verify_,33600,2400,14,constant-16                     20           6992628 ns/op
BenchmarkUpdate/Prove_,33000,3000,11,constant-16                      20        9124540794 ns/op
BenchmarkUpdate/Verify_,33000,3000,11,constant-16                     20           6412191 ns/op
BenchmarkUpdate/Prove_,32400,3600,9,constant-16                       20        13220098909 ns/op
BenchmarkUpdate/Verify_,32400,3600,9,constant-16                      20           6924624 ns/op
BenchmarkUpdate/Prove_,33600,4200,8,constant-16                       20        14364858250 ns/op
BenchmarkUpdate/Verify_,33600,4200,8,constant-16                      20           7099924 ns/op
BenchmarkUpdate/Prove_,33600,4800,7,constant-16                       20        15439059178 ns/op
BenchmarkUpdate/Verify_,33600,4800,7,constant-16                      20           6139195 ns/op
BenchmarkUpdate/Prove_,32400,5400,6,constant-16                       20        16578384297 ns/op
BenchmarkUpdate/Verify_,32400,5400,6,constant-16                      20           6555212 ns/op
BenchmarkUpdate/Prove_,30000,6000,5,constant-16                       20        17696136292 ns/op
BenchmarkUpdate/Verify_,30000,6000,5,constant-16                      20           7167172 ns/op
BenchmarkUpdate/Prove_,35100,600,-,linear-16                          20        2159222402 ns/op
BenchmarkUpdate/Verify_,35100,600,-,linear-16                         20           2698702 ns/op
BenchmarkUpdate/Prove_,35100,1200,-,linear-16                         20        3913710572 ns/op
BenchmarkUpdate/Verify_,35100,1200,-,linear-16                        20           2741960 ns/op
BenchmarkUpdate/Prove_,35100,1800,-,linear-16                         20        6353086923 ns/op
BenchmarkUpdate/Verify_,35100,1800,-,linear-16                        20           3139160 ns/op
BenchmarkUpdate/Prove_,35100,2400,-,linear-16                         20        7966427471 ns/op
BenchmarkUpdate/Verify_,35100,2400,-,linear-16                        20           3229675 ns/op
BenchmarkUpdate/Prove_,35100,3000,-,linear-16                         20        9079860864 ns/op
BenchmarkUpdate/Verify_,35100,3000,-,linear-16                        20           3065041 ns/op
BenchmarkUpdate/Prove_,35100,3600,-,linear-16                         20        13158499654 ns/op
BenchmarkUpdate/Verify_,35100,3600,-,linear-16                        20           3243003 ns/op
BenchmarkUpdate/Prove_,35100,4200,-,linear-16                         20        14351523401 ns/op
BenchmarkUpdate/Verify_,35100,4200,-,linear-16                        20           3822846 ns/op
BenchmarkUpdate/Prove_,35100,4800,-,linear-16                         20        15476400392 ns/op
BenchmarkUpdate/Verify_,35100,4800,-,linear-16                        20           3518799 ns/op
BenchmarkUpdate/Prove_,35100,5400,-,linear-16                         20        16443004798 ns/op
BenchmarkUpdate/Verify_,35100,5400,-,linear-16                        20           3532798 ns/op
BenchmarkUpdate/Prove_,35100,6000,-,linear-16                         20        17698496078 ns/op
BenchmarkUpdate/Verify_,35100,6000,-,linear-16                        20           3857218 ns/op
PASS
ok      github.com/akakou/zk-ban/test   12474.503s
```