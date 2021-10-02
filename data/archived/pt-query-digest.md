pt-query-digestとRows examine/Rows sentについて---2018-08-15 19:09:52

[pt-query-digest](https://www.percona.com/doc/percona-toolkit/LATEST/pt-query-digest.html)というMySQLのスロークエリ調査ツールがある。
ここでは、pt-query-digestでわかるRows examine/Rows sentの見方について書いておく。

### 使い方

```
pt-query-digest /var/log/mysql/slow_query.log
```
のように使う。MySQL側の設定で、スロークエリログを吐く設定にしてある必要がある。
(以降は、ISUCON7のデータを使用している)

```
### 1.1s user time, 0 system time, 28.47M rss, 81.44M vsz
### Current date: Sat Aug 18 12:34:31 2018
### Hostname: ubuntu-xenial
### Files: ./slow_query.log
### Overall: 12.37k total, 17 unique, 199.52 QPS, 5.08x concurrency ________
### Time range: 2018-08-18T12:24:26 to 2018-08-18T12:25:28
### Attribute          total     min     max     avg     95%  stddev  median
### ============     ======= ======= ======= ======= ======= ======= =======
### Exec time           315s     1ms      3s    25ms    61ms    79ms    13ms
### Lock time            17s       0      1s     1ms   596us    19ms    14us
### Rows sent         12.46k       0     100    1.03    0.99    2.14    0.99
### Rows examine     122.40M       0  10.77k  10.13k  10.29k   1.12k   9.80k
### Query size       699.69k      17     216   57.92   56.92    3.93   56.92

### Profile
### Rank Query ID                      Response time  Calls R/Call V/M   Ite
### ==== ============================= ============== ===== ====== ===== ===
###    1 0xB36405B8C0D2F0C74866C96B... 220.3156 70.0% 12092 0.0182  0.05 SELECT message
###    2 0xDA556F9115773A1A99AA0165...  31.9184 10.1%    52 0.6138  0.36 ADMIN PREPARE
###    3 0xDF664753A9FC79D32CF41F09...  22.1128  7.0%    34 0.6504  0.36 SELECT haveread
###    4 0x6CB85ADF7C155CAA91283C30...  11.3231  3.6%    80 0.1415  0.06 SELECT message
###    5 0xEC97B582BC15F46A85799218...   7.9755  2.5%    55 0.1450  0.01 SELECT message
###    6 0x63464133E79053D64BD8EFB4...   5.2594  1.7%     6 0.8766  0.29 SELECT channel
###    7 0x07890000813C4CC7111FD2D3...   4.1713  1.3%    16 0.2607  0.01 ADMIN CLOSE STMT
### MISC 0xMISC                         11.7454  3.7%    35 0.3356   0.0 <10 ITEMS>
```

最初のセクションで注目すべきは、

```
### Rows sent         12.46k       0     100    1.03    0.99    2.14    0.99
### Rows examine     122.40M       0  10.77k  10.13k  10.29k   1.12k   9.80k
```

である。
`Rows sent` は、発行されたクエリが、実際にクライアントにレスポンスした行数を意味している。
`Rows examine` は、クライアントにレスポンスを返すまでになめた行数を意味している。

例えば、 `id` カラムがPKになったusersテーブル(100レコード)に対して、

```
SELECT * FROM users WHERE id = 99;
```

というクエリを投げたときのことを考える。
idカラムはPKなため、当然WHERE句に指定すればインデックスが効く。
つまり、いきなりidが99の行を取りに行き、その行だけをクライアントに返すため、
Rows examine: 1, Rows sent: 1 ã«なる。

また、別の例では。
例えば `name` カラムがユニークかつインデックスがない場合を考える。


```
SELECT * FROM users WHERE name = "John";
```

Johnのidが `99` だった場合、インデックスがないため1件目から走査を行う必要がある。
その後、99件目に目的の行にたどり着いた後、 「 `name`　はユニークである」という前提に基づき、
その時点で走査をやめ、99行目のレコードを返す。
つまり、 Rows examine: 99, Rows sent: 1 になる。

また、上記の例で `name` がユニークでなければどうなるか？
この場合、99件目でJohnを見つけたとしても、まだJohnが残っているかもしれないため、走査を辞めることができない。
つまり、 Rows examine: 100, Rows sent: 1 になる。

さて、これらの例全てにおいて、Rows sentは1である。しかし、DB制約やインデックスにより、Rows examineの数字は大きく異なる。
そして言うまでもなく、 Rows examineの数が少ないほうがクエリの実行時間は短くなる。
また言うまでもないが、 Rows examine >= Rows sentの関係は崩れることはない。
つまり、 Rows examine / Rows sentが1に近いほど効率の良いクエリと言え、
1から離れているほどチューニングの余地がある(走査する件数を減らすことで、クエリの実行時間を減らしやすい)と言える。

例えば、pt-query-digestの結果をもう少し下まで見ていくと、下記のようなセクションがある。

```
### Query 1: 195.03 QPS, 3.55x concurrency, ID 0xB36405B8C0D2F0C74866C96BA295C56E at byte 3240347
### Scores: V/M = 0.05
### Time range: 2018-08-18T12:24:26 to 2018-08-18T12:25:28
### Attribute    pct   total     min     max     avg     95%  stddev  median
### ============ === ======= ======= ======= ======= ======= ======= =======
### Count         97   12092
### Exec time     69    220s     1ms      1s    18ms    53ms    30ms    12ms
### Lock time     47      8s     5us   773ms   658us   445us     9ms    14us
### Rows sent     94  11.81k       1       1       1       1       0       1
### Rows examine  98 121.03M   9.77k  10.77k  10.25k  10.29k  344.61   9.80k
### Query size    97 682.73k      56      58   57.82   56.92    0.76   56.92
### String:
### Databases    isubata
### Hosts        localhost
### Users        isucon
### Query_time distribution
###   1us
###  10us
### 100us
###   1ms  #################################################
###  10ms  ################################################################
### 100ms  #
###    1s  #
###  10s+
### Tables
###    SHOW TABLE STATUS FROM `isubata` LIKE 'message'\G
###    SHOW CREATE TABLE `isubata`.`message`\G
### EXPLAIN /*!50100 PARTITIONS*/
SELECT COUNT(*) as cnt FROM message WHERE channel_id = 283\G
```

これは、

```
SELECT COUNT(*) as cnt FROM message WHERE channel_id = 283
```

のようなクエリがどれだけ時間がかかったかを示している。
Rows examineとRows sentを見てみると、examineしまくっている割にsentは非常に少ない。(countだから)
そのため、countをなんとかしてなくす(専用のカラムを作っちゃうとか)ことで、パフォーマンスを向上できる可能性が高い。
