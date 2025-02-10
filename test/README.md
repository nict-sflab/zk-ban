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
BenchmarkAll/join_req-16                       5          14789137 ns/op
BenchmarkAll/verify_join_req-16                5           2365974 ns/op
BenchmarkAll/sign-16                           5          74586961 ns/op
BenchmarkAll/verify-16                         5           2412382 ns/op
BenchmarkAll/update-req_(constant)-16                  5        1090865694 ns/op
BenchmarkAll/update-verify_(constant)-16               5           5914462 ns/op
BenchmarkAll/update-req_(linear)-16                    5        1017619754 ns/op
BenchmarkAll/update-verify_(linear)-16                 5           2374631 ns/op
BenchmarkAll/update-req_(one_session)-16               5         246248717 ns/op
BenchmarkAll/update-verify_(one_session)-16            5           6695769 ns/op
PASS
ok      github.com/akakou/zk-ban/test   63.898s
```

```
akakou@ra-webs:~/zk-ban/test$ go test --bench BenchmarkUpdate . -timeout 0 -benchtime 50x
goos: linux
goarch: amd64
pkg: github.com/akakou/zk-ban/test
cpu: Intel(R) Xeon(R) Gold 5118 CPU @ 2.30GHz
BenchmarkUpdate/Prove_,16200,270,60,constant-16                       50        1047657020 ns/op
BenchmarkUpdate/Verify_,16200,270,60,constant-16                      50           4316184 ns/op
BenchmarkUpdate/Prove_,32400,270,120,constant-16                      50        1080975384 ns/op
BenchmarkUpdate/Verify_,32400,270,120,constant-16                     50           6620026 ns/op
BenchmarkUpdate/Prove_,48600,270,180,constant-16                      50        1103391831 ns/op
BenchmarkUpdate/Verify_,48600,270,180,constant-16                     50           9461505 ns/op
BenchmarkUpdate/Prove_,64800,270,240,constant-16                      50        1145360969 ns/op
BenchmarkUpdate/Verify_,64800,270,240,constant-16                     50           9505896 ns/op
BenchmarkUpdate/Prove_,81000,270,300,constant-16                      50        1543388432 ns/op
BenchmarkUpdate/Verify_,81000,270,300,constant-16                     50           9617584 ns/op
BenchmarkUpdate/Prove_,97200,270,360,constant-16                      50        1675690335 ns/op
BenchmarkUpdate/Verify_,97200,270,360,constant-16                     50          10057529 ns/op
BenchmarkUpdate/Prove_,113400,270,420,constant-16                     50        1701163085 ns/op
BenchmarkUpdate/Verify_,113400,270,420,constant-16                    50          13876565 ns/op
BenchmarkUpdate/Prove_,129600,270,480,constant-16                     50        1740866238 ns/op
BenchmarkUpdate/Verify_,129600,270,480,constant-16                    50          15805044 ns/op
BenchmarkUpdate/Prove_,145800,270,540,constant-16                     50        1743694660 ns/op
BenchmarkUpdate/Verify_,145800,270,540,constant-16                    50          15325350 ns/op
BenchmarkUpdate/Prove_,162000,270,600,constant-16                     50        1785079660 ns/op
BenchmarkUpdate/Verify_,162000,270,600,constant-16                    50          15791444 ns/op
BenchmarkUpdate/Prove_,178200,270,660,constant-16                     50        1806082856 ns/op
BenchmarkUpdate/Verify_,178200,270,660,constant-16                    50          14256451 ns/op
BenchmarkUpdate/Prove_,194400,270,720,constant-16                     50        1832530352 ns/op
BenchmarkUpdate/Verify_,194400,270,720,constant-16                    50          15874414 ns/op
BenchmarkUpdate/Prove_,210600,270,780,constant-16                     50        1853679611 ns/op
BenchmarkUpdate/Verify_,210600,270,780,constant-16                    50          15381615 ns/op
BenchmarkUpdate/Prove_,226800,270,840,constant-16                     50        1874807044 ns/op
BenchmarkUpdate/Verify_,226800,270,840,constant-16                    50          31472376 ns/op
BenchmarkUpdate/Prove_,243000,270,900,constant-16                     50        1892264470 ns/op
BenchmarkUpdate/Verify_,243000,270,900,constant-16                    50          28382368 ns/op
BenchmarkUpdate/Prove_,259200,270,960,constant-16                     50        1949115115 ns/op
BenchmarkUpdate/Verify_,259200,270,960,constant-16                    50          26008267 ns/op
BenchmarkUpdate/Prove_,275400,270,1020,constant-16                    50        1889442317 ns/op
BenchmarkUpdate/Verify_,275400,270,1020,constant-16                   50          32463487 ns/op
BenchmarkUpdate/Prove_,291600,270,1080,constant-16                    50        1948684411 ns/op
BenchmarkUpdate/Verify_,291600,270,1080,constant-16                   50          27203652 ns/op
BenchmarkUpdate/Prove_,307800,270,1140,constant-16                    50        1995432458 ns/op
BenchmarkUpdate/Verify_,307800,270,1140,constant-16                   50          29955637 ns/op
BenchmarkUpdate/Prove_,324000,270,1200,constant-16                    50        2038765343 ns/op
BenchmarkUpdate/Verify_,324000,270,1200,constant-16                   50          25637912 ns/op
BenchmarkUpdate/Prove_,16200,270,-,linear-16                          50        1016534622 ns/op
BenchmarkUpdate/Verify_,16200,270,-,linear-16                         50           2479488 ns/op
BenchmarkUpdate/Prove_,32400,270,-,linear-16                          50        1024312381 ns/op
BenchmarkUpdate/Verify_,32400,270,-,linear-16                         50           2477492 ns/op
BenchmarkUpdate/Prove_,48600,270,-,linear-16                          50        1098450334 ns/op
BenchmarkUpdate/Verify_,48600,270,-,linear-16                         50           6590260 ns/op
BenchmarkUpdate/Prove_,64800,270,-,linear-16                          50        1094549447 ns/op
BenchmarkUpdate/Verify_,64800,270,-,linear-16                         50           6646138 ns/op
BenchmarkUpdate/Prove_,81000,270,-,linear-16                          50        1171285932 ns/op
BenchmarkUpdate/Verify_,81000,270,-,linear-16                         50           9996365 ns/op
BenchmarkUpdate/Prove_,97200,270,-,linear-16                          50        1180117777 ns/op
BenchmarkUpdate/Verify_,97200,270,-,linear-16                         50           9654023 ns/op
BenchmarkUpdate/Prove_,113400,270,-,linear-16                         50        1695253838 ns/op
BenchmarkUpdate/Verify_,113400,270,-,linear-16                        50          15226627 ns/op
BenchmarkUpdate/Prove_,129600,270,-,linear-16                         50        1697327640 ns/op
BenchmarkUpdate/Verify_,129600,270,-,linear-16                        50          13885090 ns/op
BenchmarkUpdate/Prove_,145800,270,-,linear-16                         50        1752562026 ns/op
BenchmarkUpdate/Verify_,145800,270,-,linear-16                        50          15947220 ns/op
BenchmarkUpdate/Prove_,162000,270,-,linear-16                         50        1760051094 ns/op
BenchmarkUpdate/Verify_,162000,270,-,linear-16                        50          15159442 ns/op
BenchmarkUpdate/Prove_,178200,270,-,linear-16                         50        1756396688 ns/op
BenchmarkUpdate/Verify_,178200,270,-,linear-16                        50          15083859 ns/op
BenchmarkUpdate/Prove_,194400,270,-,linear-16                         50        1810250190 ns/op
BenchmarkUpdate/Verify_,194400,270,-,linear-16                        50          16038905 ns/op
BenchmarkUpdate/Prove_,210600,270,-,linear-16                         50        1800022000 ns/op
BenchmarkUpdate/Verify_,210600,270,-,linear-16                        50          14287685 ns/op
BenchmarkUpdate/Prove_,226800,270,-,linear-16                         50        1839717044 ns/op
BenchmarkUpdate/Verify_,226800,270,-,linear-16                        50          33056146 ns/op
BenchmarkUpdate/Prove_,243000,270,-,linear-16                         50        1873267370 ns/op
BenchmarkUpdate/Verify_,243000,270,-,linear-16                        50          22996161 ns/op
BenchmarkUpdate/Prove_,259200,270,-,linear-16                         50        1918305309 ns/op
BenchmarkUpdate/Verify_,259200,270,-,linear-16                        50          26109750 ns/op
BenchmarkUpdate/Prove_,275400,270,-,linear-16                         50        1937807153 ns/op
BenchmarkUpdate/Verify_,275400,270,-,linear-16                        50          33825263 ns/op
BenchmarkUpdate/Prove_,291600,270,-,linear-16                         50        1920081060 ns/op
BenchmarkUpdate/Verify_,291600,270,-,linear-16                        50          31126362 ns/op
BenchmarkUpdate/Prove_,307800,270,-,linear-16                         50        1935189323 ns/op
BenchmarkUpdate/Verify_,307800,270,-,linear-16                        50          28401462 ns/op
BenchmarkUpdate/Prove_,324000,270,-,linear-16                         50        1935334627 ns/op
BenchmarkUpdate/Verify_,324000,270,-,linear-16                        50          27311007 ns/op
BenchmarkUpdate/Prove_,16200,1,16200,constant-16                      50         152081593 ns/op
BenchmarkUpdate/Verify_,16200,1,16200,constant-16                     50           4804870 ns/op
BenchmarkUpdate/Prove_,32400,1,32400,constant-16                      50         241378465 ns/op
BenchmarkUpdate/Verify_,32400,1,32400,constant-16                     50           6983310 ns/op
BenchmarkUpdate/Prove_,48600,1,48600,constant-16                      50         290232279 ns/op
BenchmarkUpdate/Verify_,48600,1,48600,constant-16                     50           9105354 ns/op
BenchmarkUpdate/Prove_,64800,1,64800,constant-16                      50         406274804 ns/op
BenchmarkUpdate/Verify_,64800,1,64800,constant-16                     50           8502390 ns/op
BenchmarkUpdate/Prove_,81000,1,81000,constant-16                      50         451303968 ns/op
BenchmarkUpdate/Verify_,81000,1,81000,constant-16                     50          10076969 ns/op
BenchmarkUpdate/Prove_,97200,1,97200,constant-16                      50         501601717 ns/op
BenchmarkUpdate/Verify_,97200,1,97200,constant-16                     50          10375044 ns/op
BenchmarkUpdate/Prove_,113400,1,113400,constant-16                    50         534663138 ns/op
BenchmarkUpdate/Verify_,113400,1,113400,constant-16                   50          15611401 ns/op
BenchmarkUpdate/Prove_,129600,1,129600,constant-16                    50         757952154 ns/op
BenchmarkUpdate/Verify_,129600,1,129600,constant-16                   50          13308338 ns/op
BenchmarkUpdate/Prove_,145800,1,145800,constant-16                    50         773049977 ns/op
BenchmarkUpdate/Verify_,145800,1,145800,constant-16                   50          14183583 ns/op
BenchmarkUpdate/Prove_,162000,1,162000,constant-16                    50         802452217 ns/op
BenchmarkUpdate/Verify_,162000,1,162000,constant-16                   50          15676298 ns/op
BenchmarkUpdate/Prove_,178200,1,178200,constant-16                    50         821671659 ns/op
BenchmarkUpdate/Verify_,178200,1,178200,constant-16                   50          14156404 ns/op
BenchmarkUpdate/Prove_,194400,1,194400,constant-16                    50         862268607 ns/op
BenchmarkUpdate/Verify_,194400,1,194400,constant-16                   50          15539766 ns/op
BenchmarkUpdate/Prove_,210600,1,210600,constant-16                    50         935012537 ns/op
BenchmarkUpdate/Verify_,210600,1,210600,constant-16                   50          16006004 ns/op
BenchmarkUpdate/Prove_,226800,1,226800,constant-16                    50         969651856 ns/op
BenchmarkUpdate/Verify_,226800,1,226800,constant-16                   50          31515846 ns/op
BenchmarkUpdate/Prove_,243000,1,243000,constant-16                    50        1007879300 ns/op
BenchmarkUpdate/Verify_,243000,1,243000,constant-16                   50          32957018 ns/op
BenchmarkUpdate/Prove_,259200,1,259200,constant-16                    50        1369280574 ns/op
BenchmarkUpdate/Verify_,259200,1,259200,constant-16                   50          24532676 ns/op
BenchmarkUpdate/Prove_,275400,1,275400,constant-16                    50        1374481676 ns/op
BenchmarkUpdate/Verify_,275400,1,275400,constant-16                   50          27658370 ns/op
BenchmarkUpdate/Prove_,291600,1,291600,constant-16                    50        1418351914 ns/op
BenchmarkUpdate/Verify_,291600,1,291600,constant-16                   50          27500642 ns/op
BenchmarkUpdate/Prove_,307800,1,307800,constant-16                    50        1441295195 ns/op
BenchmarkUpdate/Verify_,307800,1,307800,constant-16                   50          24976160 ns/op
BenchmarkUpdate/Prove_,324000,1,324000,constant-16                    50        1456618380 ns/op
BenchmarkUpdate/Verify_,324000,1,324000,constant-16                   50          28651468 ns/op
BenchmarkUpdate/Prove_,34980,60,583,constant-16                       50         465182369 ns/op
BenchmarkUpdate/Verify_,34980,60,583,constant-16                      50           6590076 ns/op
BenchmarkUpdate/Prove_,34920,120,291,constant-16                      50         618453923 ns/op
BenchmarkUpdate/Verify_,34920,120,291,constant-16                     50           6560799 ns/op
BenchmarkUpdate/Prove_,34920,180,194,constant-16                      50         936576485 ns/op
BenchmarkUpdate/Verify_,34920,180,194,constant-16                     50           6822290 ns/op
BenchmarkUpdate/Prove_,34800,240,145,constant-16                      50        1029328763 ns/op
BenchmarkUpdate/Verify_,34800,240,145,constant-16                     50           6362695 ns/op
BenchmarkUpdate/Prove_,34800,300,116,constant-16                      50        1130757632 ns/op
BenchmarkUpdate/Verify_,34800,300,116,constant-16                     50           6312522 ns/op
BenchmarkUpdate/Prove_,34920,360,97,constant-16                       50        1640131282 ns/op
BenchmarkUpdate/Verify_,34920,360,97,constant-16                      50           6281430 ns/op
BenchmarkUpdate/Prove_,34860,420,83,constant-16                       50        1842284362 ns/op
BenchmarkUpdate/Verify_,34860,420,83,constant-16                      50           6074392 ns/op
BenchmarkUpdate/Prove_,34560,480,72,constant-16                       50        1944000312 ns/op
BenchmarkUpdate/Verify_,34560,480,72,constant-16                      50           6585977 ns/op
BenchmarkUpdate/Prove_,34560,540,64,constant-16                       50        2070246144 ns/op
BenchmarkUpdate/Verify_,34560,540,64,constant-16                      50           6075325 ns/op
BenchmarkUpdate/Prove_,34800,600,58,constant-16                       50        2196325883 ns/op
BenchmarkUpdate/Verify_,34800,600,58,constant-16                      50           6257804 ns/op
BenchmarkUpdate/Prove_,34980,660,53,constant-16                       50        2287994252 ns/op
BenchmarkUpdate/Verify_,34980,660,53,constant-16                      50           6157900 ns/op
BenchmarkUpdate/Prove_,34560,720,48,constant-16                       50        2412093394 ns/op
BenchmarkUpdate/Verify_,34560,720,48,constant-16                      50           6038679 ns/op
BenchmarkUpdate/Prove_,34320,780,44,constant-16                       50        3143105568 ns/op
BenchmarkUpdate/Verify_,34320,780,44,constant-16                      50           6279063 ns/op
BenchmarkUpdate/Prove_,34440,840,41,constant-16                       50        3248478609 ns/op
BenchmarkUpdate/Verify_,34440,840,41,constant-16                      50           6408989 ns/op
BenchmarkUpdate/Prove_,34200,900,38,constant-16                       50        3432454926 ns/op
BenchmarkUpdate/Verify_,34200,900,38,constant-16                      50           6676432 ns/op
BenchmarkUpdate/Prove_,34560,960,36,constant-16                       50        3481362456 ns/op
BenchmarkUpdate/Verify_,34560,960,36,constant-16                      50           5986114 ns/op
BenchmarkUpdate/Prove_,34680,1020,34,constant-16                      50        3590001130 ns/op
BenchmarkUpdate/Verify_,34680,1020,34,constant-16                     50           6706296 ns/op
BenchmarkUpdate/Prove_,34560,1080,32,constant-16                      50        3695718191 ns/op
BenchmarkUpdate/Verify_,34560,1080,32,constant-16                     50           6561076 ns/op
BenchmarkUpdate/Prove_,34200,1140,30,constant-16                      50        3831867702 ns/op
BenchmarkUpdate/Verify_,34200,1140,30,constant-16                     50           7196443 ns/op
BenchmarkUpdate/Prove_,34800,1200,29,constant-16                      50        3963358767 ns/op
BenchmarkUpdate/Verify_,34800,1200,29,constant-16                     50           7093990 ns/op
BenchmarkUpdate/Prove_,35000,60,-,linear-16                           50         457207978 ns/op
BenchmarkUpdate/Verify_,35000,60,-,linear-16                          50           5942348 ns/op
BenchmarkUpdate/Prove_,35000,120,-,linear-16                          50         587332369 ns/op
BenchmarkUpdate/Verify_,35000,120,-,linear-16                         50           6084631 ns/op
BenchmarkUpdate/Prove_,35000,180,-,linear-16                          50         919712804 ns/op
BenchmarkUpdate/Verify_,35000,180,-,linear-16                         50           6523115 ns/op
BenchmarkUpdate/Prove_,35000,240,-,linear-16                          50        1020704843 ns/op
BenchmarkUpdate/Verify_,35000,240,-,linear-16                         50           6194327 ns/op
BenchmarkUpdate/Prove_,35000,300,-,linear-16                          50        1059275628 ns/op
BenchmarkUpdate/Verify_,35000,300,-,linear-16                         50           2487292 ns/op
BenchmarkUpdate/Prove_,35000,360,-,linear-16                          50        1177561282 ns/op
BenchmarkUpdate/Verify_,35000,360,-,linear-16                         50           2755439 ns/op
BenchmarkUpdate/Prove_,35000,420,-,linear-16                          50        1646395632 ns/op
BenchmarkUpdate/Verify_,35000,420,-,linear-16                         50           2545822 ns/op
BenchmarkUpdate/Prove_,35000,480,-,linear-16                          50        1845567608 ns/op
BenchmarkUpdate/Verify_,35000,480,-,linear-16                         50           2384242 ns/op
BenchmarkUpdate/Prove_,35000,540,-,linear-16                          50        1971380970 ns/op
BenchmarkUpdate/Verify_,35000,540,-,linear-16                         50           2476280 ns/op
BenchmarkUpdate/Prove_,35000,600,-,linear-16                          50        2120020169 ns/op
BenchmarkUpdate/Verify_,35000,600,-,linear-16                         50           2516736 ns/op
BenchmarkUpdate/Prove_,35000,660,-,linear-16                          50        2242809253 ns/op
BenchmarkUpdate/Verify_,35000,660,-,linear-16                         50           2494113 ns/op
BenchmarkUpdate/Prove_,35000,720,-,linear-16                          50        2314948627 ns/op
BenchmarkUpdate/Verify_,35000,720,-,linear-16                         50           2555873 ns/op
BenchmarkUpdate/Prove_,35000,780,-,linear-16                          50        3049339823 ns/op
BenchmarkUpdate/Verify_,35000,780,-,linear-16                         50           2762871 ns/op
BenchmarkUpdate/Prove_,35000,840,-,linear-16                          50        3145656765 ns/op
BenchmarkUpdate/Verify_,35000,840,-,linear-16                         50           2592122 ns/op
BenchmarkUpdate/Prove_,35000,900,-,linear-16                          50        3258945040 ns/op
BenchmarkUpdate/Verify_,35000,900,-,linear-16                         50           2565649 ns/op
BenchmarkUpdate/Prove_,35000,960,-,linear-16                          50        3448944174 ns/op
BenchmarkUpdate/Verify_,35000,960,-,linear-16                         50           2651618 ns/op
BenchmarkUpdate/Prove_,35000,1020,-,linear-16                         50        3548305247 ns/op
BenchmarkUpdate/Verify_,35000,1020,-,linear-16                        50           2649017 ns/op
BenchmarkUpdate/Prove_,35000,1080,-,linear-16                         50        3605621547 ns/op
BenchmarkUpdate/Verify_,35000,1080,-,linear-16                        50           2911843 ns/op
BenchmarkUpdate/Prove_,35000,1140,-,linear-16                         50        3741054433 ns/op
BenchmarkUpdate/Verify_,35000,1140,-,linear-16                        50           2842184 ns/op
BenchmarkUpdate/Prove_,35000,1200,-,linear-16                         50        3881443272 ns/op
BenchmarkUpdate/Verify_,35000,1200,-,linear-16                        50           2664456 ns/op
PASS
ok      github.com/akakou/zk-ban/test   11536.588s
```