# Micro-benchmark for eval vs. compile/run of expr ASTs

The code implement an AST (`ast.go`) with 3 possible nodes:

- simple int value
- binary operation
- function application (from 0 to 5 arguments0

Next to this it implements both:

- evaluation using recursive-descent interpretation (`ast.go`: `eval` method)
- evaluation by compiling first the AST to a stack machine program (`comp.go`: `compile`) then execution (`run.go`: `run`)

The main program is just an interactive test to make experiments.

The main benchmark is in `bench_test.go` and can be run as follows:

1. edit the `doRun` constant to set it to `false` (use `eval`)
2. run with `go test -bench . -benchmem -count 3 | tee bench-eval.log`
3. edit the `doRun` constant to set it to `true` (use compile/run)
4. run with `go test -bench . -benchmem -count 3 | tee bench-run.log`
5. `benchstat bench-eval.log bench-run.log`

For each "benchmark" the name indicates what is being benchmarked. There are 3 fields:

- simple vs. complex. "Simple" expressions only contain simple values
  and binary operations (no function applications). "Complex"
  expressions contain at least one function application.
- "size". This is an approximate number of items in the expression and
  thus determines approximate expression complexity. (it's the argument given
  to the `gen` expr auto-generator - check the code to see how it is
  used).
- "node": which expression is being tested. For each size, `numExpr`
  expressions are pre-generated using a fixed RNG seed, so that
  different invocations of `go test` (with e.g. different `doRun`
  values) use the same expressions. `numExpr = 3` in the code, but can
  be modified manually to run more benchmarks.

Then:

- when `doRun = false`, the expression is evaluated N times
  (`testing.B.N`) using its `eval` method.
- when `doRun = true`, the expression is compiled once using `compile`
  then executed N times using `run`.

On my machine I observe the following differences:

```
name                        old time/op    new time/op    delta
Blah/simple-sz1-node0-4       7.36ns ± 2%   13.50ns ± 1%   +83.40%  (p=0.000 n=9+10)
Blah/simple-sz5-node0-4        102ns ± 3%      43ns ± 1%   -57.47%  (p=0.000 n=10+9)
Blah/simple-sz5-node1-4        149ns ± 2%      59ns ± 2%   -60.29%  (p=0.000 n=9+10)
Blah/simple-sz5-node2-4        102ns ± 4%      43ns ± 1%   -57.43%  (p=0.000 n=10+10)
Blah/simple-sz10-node0-4       247ns ± 1%      94ns ±10%   -62.06%  (p=0.000 n=10+10)
Blah/simple-sz10-node1-4       241ns ± 1%      91ns ± 1%   -62.49%  (p=0.000 n=8+8)
Blah/simple-sz10-node2-4       288ns ± 1%     111ns ±10%   -61.44%  (p=0.000 n=8+10)
Blah/simple-sz100-node0-4     2.80µs ± 3%    0.88µs ± 0%   -68.39%  (p=0.000 n=10+10)
Blah/simple-sz100-node1-4     2.74µs ± 3%    0.88µs ± 3%   -68.07%  (p=0.000 n=10+9)
Blah/simple-sz100-node2-4     3.12µs ± 3%    1.02µs ± 9%   -67.26%  (p=0.000 n=10+10)
Blah/complex-sz1-node0-4      48.3ns ± 2%    25.4ns ±12%   -47.34%  (p=0.000 n=10+10)
Blah/complex-sz5-node0-4       232ns ± 6%      70ns ± 1%   -69.71%  (p=0.000 n=10+8)
Blah/complex-sz5-node1-4       142ns ± 5%      58ns ±10%   -59.51%  (p=0.000 n=9+10)
Blah/complex-sz5-node2-4       225ns ± 2%      70ns ± 1%   -68.85%  (p=0.000 n=10+10)
Blah/complex-sz10-node0-4      498ns ± 4%     137ns ± 0%   -72.59%  (p=0.000 n=10+9)
Blah/complex-sz10-node1-4      403ns ± 6%     105ns ± 3%   -73.89%  (p=0.000 n=10+9)
Blah/complex-sz10-node2-4      274ns ± 2%      91ns ± 3%   -66.84%  (p=0.000 n=10+10)
Blah/complex-sz100-node0-4    4.30µs ± 2%    0.95µs ± 1%   -78.04%  (p=0.000 n=10+8)
Blah/complex-sz100-node1-4    4.35µs ± 5%    0.99µs ± 3%   -77.23%  (p=0.000 n=10+9)
Blah/complex-sz100-node2-4    4.48µs ± 2%    1.00µs ±12%   -77.65%  (p=0.000 n=10+10)

name                        old alloc/op   new alloc/op   delta
Blah/simple-sz1-node0-4        0.00B          0.00B           ~     (all equal)
Blah/simple-sz5-node0-4        16.0B ± 0%      0.0B       -100.00%  (p=0.000 n=10+10)
Blah/simple-sz5-node1-4        24.0B ± 0%      0.0B       -100.00%  (p=0.000 n=10+10)
Blah/simple-sz5-node2-4        16.0B ± 0%      0.0B       -100.00%  (p=0.000 n=10+10)
Blah/simple-sz10-node0-4       40.0B ± 0%      0.0B       -100.00%  (p=0.000 n=10+10)
Blah/simple-sz10-node1-4       40.0B ± 0%      0.0B       -100.00%  (p=0.000 n=10+10)
Blah/simple-sz10-node2-4       48.0B ± 0%      0.0B       -100.00%  (p=0.000 n=10+10)
Blah/simple-sz100-node0-4       440B ± 0%        0B       -100.00%  (p=0.000 n=10+10)
Blah/simple-sz100-node1-4       432B ± 0%        0B       -100.00%  (p=0.000 n=10+10)
Blah/simple-sz100-node2-4       472B ± 0%        0B       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz1-node0-4       8.00B ± 0%     0.00B       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz5-node0-4       48.0B ± 0%      0.0B       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz5-node1-4       24.0B ± 0%      0.0B       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz5-node2-4       48.0B ± 0%      0.0B       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz10-node0-4       112B ± 0%        0B       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz10-node1-4       120B ± 0%        0B       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz10-node2-4      56.0B ± 0%      0.0B       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz100-node0-4    1.00kB ± 0%    0.00kB       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz100-node1-4    1.02kB ± 0%    0.00kB       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz100-node2-4    1.18kB ± 0%    0.00kB       -100.00%  (p=0.000 n=10+10)

name                        old allocs/op  new allocs/op  delta
Blah/simple-sz1-node0-4         0.00           0.00           ~     (all equal)
Blah/simple-sz5-node0-4         2.00 ± 0%      0.00       -100.00%  (p=0.000 n=10+10)
Blah/simple-sz5-node1-4         3.00 ± 0%      0.00       -100.00%  (p=0.000 n=10+10)
Blah/simple-sz5-node2-4         2.00 ± 0%      0.00       -100.00%  (p=0.000 n=10+10)
Blah/simple-sz10-node0-4        5.00 ± 0%      0.00       -100.00%  (p=0.000 n=10+10)
Blah/simple-sz10-node1-4        5.00 ± 0%      0.00       -100.00%  (p=0.000 n=10+10)
Blah/simple-sz10-node2-4        6.00 ± 0%      0.00       -100.00%  (p=0.000 n=10+10)
Blah/simple-sz100-node0-4       55.0 ± 0%       0.0       -100.00%  (p=0.000 n=10+10)
Blah/simple-sz100-node1-4       54.0 ± 0%       0.0       -100.00%  (p=0.000 n=10+10)
Blah/simple-sz100-node2-4       59.0 ± 0%       0.0       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz1-node0-4        1.00 ± 0%      0.00       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz5-node0-4        5.00 ± 0%      0.00       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz5-node1-4        3.00 ± 0%      0.00       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz5-node2-4        5.00 ± 0%      0.00       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz10-node0-4       10.0 ± 0%       0.0       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz10-node1-4       7.00 ± 0%      0.00       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz10-node2-4       6.00 ± 0%      0.00       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz100-node0-4      75.0 ± 0%       0.0       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz100-node1-4      79.0 ± 0%       0.0       -100.00%  (p=0.000 n=10+10)
Blah/complex-sz100-node2-4      78.0 ± 0%       0.0       -100.00%  (p=0.000 n=10+10)
```
