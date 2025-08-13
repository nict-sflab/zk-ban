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
[nix-shell:/mnt/c/Users/Kosei Akama/Documents/zk-ban/test]$ go test --bench ^BenchmarkUpdate$ . -timeout 0 -benchtime 20x
[1,1,1,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]u: [900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900 900]
p: [1 16 26 36 46 55 65 75 85 95 105 115 125 135 145 155 165 175 185 195 205 215 225 235 245 255 265 275 285 295 305 315 325 335 345 355 365 375 385 395 405 415 425 435 445 455 465 475 485 495 505 515 525 535 545 555 565 575 585 595 605 615 625 635 645 655 665 675 685 695 705 715 725 735 745 755 765 775 785 795 805 815 825 835 845 855 865 875 885 895 905 915 925 935 945 955 965 975 985 995 1005 1015 1025 1035 1045 1055 1065 1075 1085 1095 1105 1115 1125 1135 1145 1155 1165 1175 1185 1195 1205 1215 1225 1235 1245 1255 1265 1275 1285 1295 1305 1315 1325 1335 1345 1355 1365 1375 1385 1395 1405 1415 1425 1435 1445 1455 1465 1475 1485 1495 1505 1515 1525 1535 1545 1555 1565 1575 1585 1595 1605 1615 1625 1635 1645 1655 1665 1675 1685 1695 1705 1715 1725 1735 1745 1755 1765 1775 1785 1795]
g: [1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 2 2 3 4 5 6 8 9 10 12 13 15 17 18 20 22 25 27 29 32 35 38 41 44 47 51 55 59 63 68 73 78 83 88 94 100 107 113 120 128 135 143 152 161 170 179 189 200 211 222 234 246 259 272 287 301 316 331 347 363 380 398 416 435 454 474 495 516 538 560 583 607 631 656 682 708 735 762 790 819 848 878 909 940 971 1004 1036 1070 1103 1138 1172 1208 1243 1279 1316 1353 1390 1427 1465 1503 1541 1579 1618 1656 1695 1734 1772 1811 1850 1888 1926 1965 2002 2040 2077 2114 2151 2187 2222 2257 2292 2325 2358 2391 2422 2453 2483 2512 2540 2568 2594 2619 2643 2666 2688 2709 2728 2746 2764 2779 2794 2807 2819 2829 2838 2846 2853 2858 2861 2863 2864]
goos: linux
goarch: amd64
pkg: github.com/akakou/zk-ban/test
cpu: Intel(R) Core(TM) i9-14900
BenchmarkUpdate/Prove:_uniformuniform:_180-162000-32                  20         964258406 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_180-162000-32                     20         121993594 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_180-162000-32                                 20            828684 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_180-324000-32                                  20        1640575432 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_180-324000-32                     20         254098075 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_180-324000-32                                 20           1030147 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_180-486000-32                                  20        2299278902 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_180-486000-32                     20         338546184 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_180-486000-32                                 20           1171029 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_180-648000-32                                  20        2831865169 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_180-648000-32                     20         445456379 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_180-648000-32                                 20           1101728 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_180-810000-32                                  20        3315537505 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_180-810000-32                     20         539277315 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_180-810000-32                                 20           1037499 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_180-972000-32                                  20        4080812978 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_180-972000-32                     20         643971023 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_180-972000-32                                 20           1074208 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_180-1134000-32                                 20        4521095028 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_180-1134000-32                    20         733747657 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_180-1134000-32                                20           1031014 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_180-1296000-32                                 20        4885921732 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_180-1296000-32                    20         846760158 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_180-1296000-32                                20           1007613 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_180-1458000-32                                 20        5337321829 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_180-1458000-32                    20         925916132 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_180-1458000-32                                20           1092821 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_180-1620000-32                                 20        5755804065 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_180-1620000-32                    20        1030896157 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_180-1620000-32                                20            972272 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_180-162000-32                        20        1163592928 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_180-162000-32                   20         137730693 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_180-162000-32                               20           1230299 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_180-324000-32                                20        1696180008 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_180-324000-32                   20         249192126 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_180-324000-32                               20           1122507 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_180-486000-32                                20        2328048434 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_180-486000-32                   20         342031408 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_180-486000-32                               20           1144631 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_180-648000-32                                20        2832359320 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_180-648000-32                   20         441725721 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_180-648000-32                               20           1070763 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_180-810000-32                                20        3302334059 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_180-810000-32                   20         536815640 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_180-810000-32                               20           1066553 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_180-972000-32                                20        4098140581 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_180-972000-32                   20         636231067 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_180-972000-32                               20           1061965 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_180-1134000-32                               20        4441915887 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_180-1134000-32                  20         728542700 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_180-1134000-32                              20           1129037 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_180-1296000-32                               20        4972196171 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_180-1296000-32                  20         831525934 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_180-1296000-32                              20           1018442 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_180-1458000-32                               20        5310911007 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_180-1458000-32                  20         926943200 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_180-1458000-32                              20           1070816 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_180-1620000-32                               20        5748024840 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_180-1620000-32                  20        1026584693 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_180-1620000-32                              20           1013895 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_180-162000-32                                        20        1161346059 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_180-162000-32                           20         141596086 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_180-162000-32                                       20           1121056 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_180-324000-32                                        20        1699206078 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_180-324000-32                           20         250707452 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_180-324000-32                                       20           1238072 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_180-486000-32                                        20        2300919634 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_180-486000-32                           20         337089578 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_180-486000-32                                       20           1210715 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_180-648000-32                                        20        2825316163 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_180-648000-32                           20         446188442 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_180-648000-32                                       20           1052276 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_180-810000-32                                        20        3321738997 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_180-810000-32                           20         554742966 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_180-810000-32                                       20           1125692 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_180-972000-32                                        20        4093038368 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_180-972000-32                           20         635686869 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_180-972000-32                                       20           1122347 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_180-1134000-32                                       20        4443085982 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_180-1134000-32                          20         730614512 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_180-1134000-32                                      20           1199649 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_180-1296000-32                                       20        4868630278 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_180-1296000-32                          20         832180276 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_180-1296000-32                                      20            975150 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_180-1458000-32                                       20        5389156621 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_180-1458000-32                          20         948936485 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_180-1458000-32                                      20           1014216 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_180-1620000-32                                       20        5795180312 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_180-1620000-32                          20        1032996716 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_180-1620000-32                                      20           1015107 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_180-162000#01-32                                       20        1181615515 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_180-162000#01-32                          20         140288146 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_180-162000#01-32                                      20           1076350 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_360-162000-32                                          20        1436075447 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_360-162000-32                             20         140473395 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_360-162000-32                                         20            942053 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_540-162000-32                                          20        1814312389 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_540-162000-32                             20         138631248 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_540-162000-32                                         20           1302436 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_720-162000-32                                          20        2037214375 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_720-162000-32                             20         139219398 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_720-162000-32                                         20           1186090 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_900-162000-32                                          20        2317312552 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_900-162000-32                             20         141854513 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_900-162000-32                                         20           1190288 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_1080-162000-32                                         20        2534410578 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_1080-162000-32                            20         141330325 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_1080-162000-32                                        20           1067966 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_1260-162000-32                                         20        2758743896 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_1260-162000-32                            20         141316776 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_1260-162000-32                                        20           1093576 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_1440-162000-32                                         20        3413692579 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_1440-162000-32                            20         140777257 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_1440-162000-32                                        20           1034050 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_1620-162000-32                                         20        3617540608 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_1620-162000-32                            20         141573643 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_1620-162000-32                                        20           1077503 ns/op
BenchmarkUpdate/Prove:_uniformuniform:_1800-162000-32                                         20        3836388842 ns/op
BenchmarkUpdate/Verify-Precomputes:_uniformuniform:_1800-162000-32                            20         143173517 ns/op
BenchmarkUpdate/Verify:_uniformuniform:_1800-162000-32                                        20           1035885 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_180-162000#01-32                             20        1154450521 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_180-162000#01-32                20         138165526 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_180-162000#01-32                            20           1236293 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_360-162000-32                                20        1423647173 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_360-162000-32                   20         138296714 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_360-162000-32                               20           1199303 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_540-162000-32                                20        1816818509 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_540-162000-32                   20         141363279 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_540-162000-32                               20           1276272 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_720-162000-32                                20        2065233299 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_720-162000-32                   20         142101170 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_720-162000-32                               20           1208673 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_900-162000-32                                20        2291768915 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_900-162000-32                   20         140780749 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_900-162000-32                               20           1163814 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_1080-162000-32                               20        2511463107 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_1080-162000-32                  20         140406278 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_1080-162000-32                              20           1051769 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_1260-162000-32                               20        2794350016 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_1260-162000-32                  20         140875153 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_1260-162000-32                              20           1130425 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_1440-162000-32                               20        3403489842 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_1440-162000-32                  20         140721709 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_1440-162000-32                              20           1031440 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_1620-162000-32                               20        3621032834 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_1620-162000-32                  20         141945039 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_1620-162000-32                              20            953580 ns/op
BenchmarkUpdate/Prove:_proportionalproportional:_1800-162000-32                               20        3816646794 ns/op
BenchmarkUpdate/Verify-Precomputes:_proportionalproportional:_1800-162000-32                  20         143985232 ns/op
BenchmarkUpdate/Verify:_proportionalproportional:_1800-162000-32                              20           1006540 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_180-162000#01-32                                     20        1149755520 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_180-162000#01-32                        20         139430402 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_180-162000#01-32                                    20           1198264 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_360-162000-32                                        20        1417885014 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_360-162000-32                           20         137504875 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_360-162000-32                                       20           1041222 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_540-162000-32                                        20        1788985239 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_540-162000-32                           20         139613067 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_540-162000-32                                       20           1101810 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_720-162000-32                                        20        2045858568 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_720-162000-32                           20         139132802 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_720-162000-32                                       20           1176408 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_900-162000-32                                        20        2274816521 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_900-162000-32                           20         141143305 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_900-162000-32                                       20           1021407 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_1080-162000-32                                       20        2553904291 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_1080-162000-32                          20         141720819 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_1080-162000-32                                      20           1069562 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_1260-162000-32                                       20        2793192095 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_1260-162000-32                          20         142697014 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_1260-162000-32                                      20           1075239 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_1440-162000-32                                       20        3371351267 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_1440-162000-32                          20         142671431 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_1440-162000-32                                      20           1036453 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_1620-162000-32                                       20        3608837147 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_1620-162000-32                          20         142024636 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_1620-162000-32                                      20           1111251 ns/op
BenchmarkUpdate/Prove:_gaussiangaussian:_1800-162000-32                                       20        3825819267 ns/op
BenchmarkUpdate/Verify-Precomputes:_gaussiangaussian:_1800-162000-32                          20         144624819 ns/op
BenchmarkUpdate/Verify:_gaussiangaussian:_1800-162000-32                                      20            986551 ns/op
PASS
ok      github.com/akakou/zk-ban/test   6601.260s
```