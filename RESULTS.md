Test Results:

```
$ echo "GET http://localhost:8080" | vegeta attack -rate 1000 -duration 10s| vegeta report
Requests      [total, rate, throughput]         10000, 1000.11, 928.81
Duration      [total, attack, wait]             10.766s, 9.999s, 767.509ms
Latencies     [min, mean, 50, 90, 95, 99, max]  5.274ms, 21.139ms, 10.29ms, 10.413ms, 71.925ms, 777.56ms, 1s
Bytes In      [total, mean]                     330000, 33.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:10000  
Error Set:
```

```
$ echo "GET http://localhost:8081" | vegeta attack -rate 1000 -duration 10s| vegeta report
Requests      [total, rate, throughput]         10000, 1000.10, 926.88
Duration      [total, attack, wait]             10.789s, 9.999s, 789.91ms
Latencies     [min, mean, 50, 90, 95, 99, max]  5.875ms, 19.256ms, 10.395ms, 12.325ms, 36.183ms, 162.737ms, 1.003s
Bytes In      [total, mean]                     330000, 33.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:10000  
Error Set:
```
