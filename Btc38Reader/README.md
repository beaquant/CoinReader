# Go btc38 datas reader

*Description: This library for readout history trade datas and orders from https://www.btc38.com

## Installation


```shell
go get -u github.com/jojopoper/CoinReader/Btc38Reader/...
```

## Usage


```go
>  BTC / XRP Open orders (Records length = 20)
>      ************ Buy ************                         ************ Sell ************
> Price         Amount          Total                   Price           Amount          Total
> 0.00001071    868.80058877    0.00930485              0.00001074      15933.88623733  0.17112994
> 0.00001068    59108.90600647  0.63128312              0.00001075      5541.50671585   0.05957120
> 0.00001067    174.38100622    0.00186065              0.00001076      7627.39017895   0.08207072
> 0.00001066    107.01128966    0.00114074              0.00001077      53559.58720111  0.57683675
> 0.00001065    2091.75028198   0.02227714              0.00001078      18262.69811700  0.19687189
> 0.00001063    5075.69806018   0.05395467              0.00001079      2862.71618242   0.03088871
> 0.00001062    7051.35827198   0.07488542              0.00001080      300995.11957211 3.25074729
> 0.00001061    21677.72473339  0.23000066              0.00001081      19926.68093753  0.21540742
> 0.00001060    3146.68305408   0.03335484              0.00001082      14420.36680169  0.15602837
> 0.00001059    7483.99336671   0.07925549              0.00001083      74.30910089     0.00080477
> 0.00001055    19556.92511849  0.20632556              0.00001084      282.91366851    0.00306678
> 0.00001054    14699.20089799  0.15492958              0.00001089      34186.63283476  0.37229243
> 0.00001053    500.00000000    0.00526500              0.00001090      11007.56203054  0.11998243
> 0.00001052    14727.14614684  0.15492958              0.00001091      1363.08383464   0.01487124
> 0.00001051    90038.88201713  0.94630865              0.00001092      627.74967791    0.00685503
> 0.00001050    300000.00000000 3.15000000              0.00001093      859.89511103    0.00939865
> 0.00001049    10200.00965399  0.10699810              0.00001095      15300.00000000  0.16753500
> 0.00001047    1442.91432171   0.01510731              0.00001096      310997.00000000 3.40852712
> 0.00001044    41045.40650793  0.42851404              0.00001097      2140.00000000   0.02347580
> 0.00001043    11998.65622510  0.12514598              0.00001099      363.96724294    0.00400000

>  BTC / STR Trade history datas (Records length = 20)
> DateTime              Type    Price           Amount          Total
> 2016-06-02 10:31:24   sell    0.00000279      189.82724616    0.00052961
> 2016-06-02 10:30:11   buy     0.00000281      4.05845040      0.00001140
> 2016-06-02 10:04:37   sell    0.00000280      362.67750000    0.00101549
> 2016-06-02 10:03:37   sell    0.00000280      3287.26110000   0.00920433
> 2016-06-02 09:59:02   sell    0.00000280      1318.77019967   0.00369255
> 2016-06-02 09:59:02   sell    0.00000280      41.51420032     0.00011623
> 2016-06-02 09:52:51   sell    0.00000280      1720.67940564   0.00481790
> 2016-06-02 09:52:51   sell    0.00000280      6473.91109325   0.01812695
> 2016-06-02 09:52:51   sell    0.00000281      67.53130110     0.00018976
> 2016-06-02 09:41:37   sell    0.00000281      994.78186616    0.00279533
> 2016-06-02 09:41:37   sell    0.00000281      12443.74783384  0.03496693
> 2016-06-02 09:34:51   sell    0.00000281      985.59724999    0.00276952
> 2016-06-02 09:33:34   sell    0.00000281      985.59724999    0.00276952
> 2016-06-02 09:33:02   sell    0.00000281      985.59724999    0.00276952
> 2016-06-02 09:32:23   sell    0.00000281      3780.69155000   0.01062374
> 2016-06-02 09:32:12   sell    0.00000281      2985.00000000   0.00838785
> 2016-06-02 09:32:02   sell    0.00000281      995.00000000    0.00279595
> 2016-06-02 09:31:48   sell    0.00000281      1990.00000000   0.00559190
> 2016-06-02 09:31:37   sell    0.00000281      4548.96090000   0.01278258
> 2016-06-02 09:31:26   sell    0.00000281      2985.00000000   0.00838785
```