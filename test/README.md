# Test

## Benchmark

### How to work

```sh
go test --bench ^BenchmarkAll$ . -timeout 0 -benchtime 5x
go test --bench ^BenchmarkUpdate$ . -timeout 0 -benchtime 5x
```

### Our result

Environments:
- CPU: Intel Core i9-14900 @ 32x 1.997GHz
- RAM: 1416MiB / 61592MiB
- OS: NixOS 25.05 (on the Windows Subsystem for Linux)

Result:
```
go test --bench ^BenchmarkAll$ . -timeout 0 -benchtime 20x
[1,1,1,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]u: [900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900]
p: [1 16 26 36 46 55 65 75 85 95 105 115 125 135 145 155 165 175 185 195 205 215 225 235 245 255 265 275 285 295 305 315 325 335 345 355 365 375 385 395 405 415 425 435 445 455 465 475 485 495 505 515 525 535 545 555 565 575 585 595 605 615 625 635 645 655 665 675 685 695 705 715 725 735 745 755 765 775 785 795 805 815 825 835 845 855 865 875 885 895 905 915 925 935 945 955 965 975 985 995 1005 1015 1025 1035 1045 1055 1065 1075 1085 1095 1105 1115 1125 1135 1145 1155 1165 1175 1185 1195 1205 1215 1225 1235 1245 1255 1265 1275 1285 1295 1305 1315 1325 1335 1345 1355 1365 1375 1385 1395 1405 1415 1425 1435 1445 1455 1465 1475 1485 1495 1505 1515 1525 1535 1545 1555 1565 1575 1585 1595 1605 1615 1625 1635 1645 1655 1665 1675 1685 1695 1705 1715 1725 1735 1745 1755 1765 1775 1785 1795]
g: [1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 2 2 3 4 5 6 8 9 10 12 13 15 17 18 20 22 25 27 29 32 35 38 41 44 47 51 55 59 63 68 73 78 83 88 94 100 107 113 120 128 135 143 152 161 170 179 189 200 211 222 234 246 259 272 287 301 316 331 347 363 380 398 416 435 454 474 495 516 538 560 583 607 631 656 682 708 735 762 790 819 848 878 909 940 971 1004 1036 1070 1103 1138 1172 1208 1243 1279 1316 1353 1390 1427 1465 1503 1541 1579 1618 1656 1695 1734 1772 1811 1850 1888 1926 1965 2002 2040 2077 2114 2151 2187 2222 2257 2292 2325 2358 2391 2422 2453 2483 2512 2540 2568 2594 2619 2643 2666 2688 2709 2728 2746 2764 2779 2794 2807 2819 2829 2838 2846 2853 2858 2861 2863 2864]
goos: linux
goarch: amd64
pkg: github.com/akakou/zk-ban/test
cpu: Intel(R) Core(TM) i9-14900
BenchmarkAll/join_req-32                      20           6054634 ns/op
BenchmarkAll/verify_join_req-32               20            752986 ns/op
BenchmarkAll/issue_credential-32              20             67578 ns/op
BenchmarkAll/sign-32                          20          34449133 ns/op
BenchmarkAll/verify-32                        20            900581 ns/op
BenchmarkAll/update-req:_uniform-32                   20         705204861 ns/op
BenchmarkAll/update-verify-precomputes:_uniform-32                    20          84419111 ns/op
BenchmarkAll/update-verify_:_uniform-32                               20            847144 ns/op
BenchmarkAll/update-req:_proportional-32                              20         708906138 ns/op
BenchmarkAll/update-verify-precomputes:_proportional-32               20          83828385 ns/op
BenchmarkAll/update-verify_:_proportional-32                          20            945367 ns/op
BenchmarkAll/update-req:_gaussian-32                                  20         839578733 ns/op
BenchmarkAll/update-verify-precomputes:_gaussian-32                   20          93842770 ns/op
BenchmarkAll/update-verify_:_gaussian-32                              20           1098266 ns/op
PASS
ok      github.com/akakou/zk-ban/test   90.292s
```

```
akakou@ra-webs:~/zk-ban/test$ go test --bench BenchmarkUpdate . -timeout 0 -benchtime 20x
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