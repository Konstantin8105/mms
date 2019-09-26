# mms
Memory managenet for slice


```
go test -v -bench=. -benchmem -count=10 > bench.prof
benchstat bench.prof 

name       time/op
/Direct-5   190µs ±21%
/Cache-5    203µs ±24%

name       alloc/op
/Direct-5  23.1kB ± 0%
/Cache-5   4.07kB ± 0%

name       allocs/op
/Direct-5     123 ± 0%
/Cache-5      123 ± 0%
```

```
go test -v -bench=. -benchmem

=== RUN   Test
--- PASS: Test (0.00s)
goos: linux
goarch: amd64
pkg: github.com/Konstantin8105/mms
Benchmark/Direct-5         	    5826	    178096 ns/op	   23071 B/op	     123 allocs/op
Benchmark/Cache-5          	    9283	    239751 ns/op	    4068 B/op	     123 allocs/op
PASS
ok  	github.com/Konstantin8105/mms	3.307s
```


```
go test -v -bench=. -benchmem -memprofile=mem.prof
go tool pprof mem.prof

(pprof) top10 -cum
Showing nodes accounting for 123.56MB, 99.60% of 124.06MB total
Dropped 5 nodes (cum <= 0.62MB)
      flat  flat%   sum%        cum   cum%
         0     0%     0%   121.06MB 97.58%  github.com/Konstantin8105/mms.Benchmark.func1.1
  105.06MB 84.68% 84.68%   105.06MB 84.68%  github.com/Konstantin8105/mms.(*Direct).Get
      16MB 12.90% 97.58%       16MB 12.90%  github.com/Konstantin8105/mms.(*Cache).Put
         0     0% 97.58%     2.50MB  2.02%  github.com/Konstantin8105/mms.Benchmark.func1
    2.50MB  2.02% 99.60%     2.50MB  2.02%  github.com/Konstantin8105/mms.getChan
         0     0% 99.60%     2.50MB  2.02%  testing.(*B).launch
         0     0% 99.60%     2.50MB  2.02%  testing.(*B).runN
(pprof) exit
```

```
go test -v -bench=. -benchmem -cpuprofile=cpu.prof
go tool pprof cpu.prof
(pprof) top10 -cum
Showing nodes accounting for 1.89s, 40.04% of 4.72s total
Dropped 71 nodes (cum <= 0.02s)
Showing top 10 nodes out of 84
      flat  flat%   sum%        cum   cum%
     0.14s  2.97%  2.97%      3.61s 76.48%  github.com/Konstantin8105/mms.Benchmark.func1.1
     0.03s  0.64%  3.60%      2.64s 55.93%  math/rand.Float64
     0.07s  1.48%  5.08%      2.61s 55.30%  math/rand.(*Rand).Float64
     0.03s  0.64%  5.72%      2.54s 53.81%  math/rand.(*Rand).Int63
     0.23s  4.87% 10.59%      2.51s 53.18%  math/rand.(*lockedSource).Int63
     0.52s 11.02% 21.61%      1.25s 26.48%  sync.(*Mutex).Lock
     0.44s  9.32% 30.93%      0.79s 16.74%  sync.(*Mutex).Unlock
     0.42s  8.90% 39.83%      0.73s 15.47%  sync.(*Mutex).lockSlow
     0.01s  0.21% 40.04%      0.65s 13.77%  runtime.mcall
         0     0% 40.04%      0.57s 12.08%  runtime.park_m
(pprof) exit
```

