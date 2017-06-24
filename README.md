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
Blah/simple-sz1-node0-4       7.47ns ± 1%   13.59ns ± 2%  +81.85%         (p=0.007 n=3+10)
Blah/simple-sz5-node0-4       96.0ns ± 1%    88.3ns ± 3%   -8.03%         (p=0.007 n=3+10)
Blah/simple-sz5-node1-4        142ns ± 1%     124ns ± 2%  -12.47%         (p=0.000 n=3+10)
Blah/simple-sz5-node2-4       96.7ns ± 1%    88.0ns ± 3%   -9.03%         (p=0.007 n=3+10)
Blah/simple-sz10-node0-4       239ns ± 1%     198ns ± 2%  -17.11%         (p=0.007 n=3+10)
Blah/simple-sz10-node1-4       238ns ± 1%     198ns ± 2%  -16.89%         (p=0.007 n=3+10)
Blah/simple-sz10-node2-4       283ns ± 1%     236ns ± 1%  -16.78%         (p=0.007 n=3+10)
Blah/simple-sz100-node0-4     2.74µs ± 1%    2.10µs ± 1%  -23.36%         (p=0.007 n=3+10)
Blah/simple-sz100-node1-4     2.70µs ± 0%    2.12µs ± 5%  -21.68%         (p=0.007 n=3+10)
Blah/simple-sz100-node2-4     3.10µs ± 1%    2.28µs ± 8%  -26.25%         (p=0.007 n=3+10)
Blah/complex-sz1-node0-4      45.3ns ± 1%    43.9ns ± 2%   -3.27%         (p=0.007 n=3+10)
Blah/complex-sz5-node0-4       219ns ± 0%     154ns ± 2%  -29.45%         (p=0.007 n=3+10)
Blah/complex-sz5-node1-4       138ns ± 0%     119ns ± 2%  -13.20%         (p=0.000 n=3+10)
Blah/complex-sz5-node2-4       221ns ± 0%     155ns ± 3%  -30.11%         (p=0.000 n=3+10)
Blah/complex-sz10-node0-4      500ns ± 2%     306ns ± 2%  -38.92%         (p=0.007 n=3+10)
Blah/complex-sz10-node1-4      392ns ± 2%     210ns ± 3%  -46.58%         (p=0.007 n=3+10)
Blah/complex-sz10-node2-4      271ns ± 0%     200ns ± 1%  -26.29%         (p=0.007 n=3+10)
Blah/complex-sz100-node0-4    4.29µs ± 1%    2.17µs ± 2%  -49.52%         (p=0.007 n=3+10)
Blah/complex-sz100-node1-4    4.30µs ± 2%    2.31µs ± 2%  -46.23%         (p=0.007 n=3+10)
Blah/complex-sz100-node2-4    4.48µs ± 1%    2.12µs ± 1%  -52.66%         (p=0.007 n=3+10)

name                        old alloc/op   new alloc/op   delta
Blah/simple-sz1-node0-4       0.00B ±NaN%    0.00B ±NaN%     ~     (all samples are equal)
Blah/simple-sz5-node0-4        16.0B ± 0%     16.0B ± 0%     ~     (all samples are equal)
Blah/simple-sz5-node1-4        24.0B ± 0%     24.0B ± 0%     ~     (all samples are equal)
Blah/simple-sz5-node2-4        16.0B ± 0%     16.0B ± 0%     ~     (all samples are equal)
Blah/simple-sz10-node0-4       40.0B ± 0%     40.0B ± 0%     ~     (all samples are equal)
Blah/simple-sz10-node1-4       40.0B ± 0%     40.0B ± 0%     ~     (all samples are equal)
Blah/simple-sz10-node2-4       48.0B ± 0%     48.0B ± 0%     ~     (all samples are equal)
Blah/simple-sz100-node0-4       440B ± 0%      440B ± 0%     ~     (all samples are equal)
Blah/simple-sz100-node1-4       432B ± 0%      432B ± 0%     ~     (all samples are equal)
Blah/simple-sz100-node2-4       472B ± 0%      472B ± 0%     ~     (all samples are equal)
Blah/complex-sz1-node0-4       8.00B ± 0%     8.00B ± 0%     ~     (all samples are equal)
Blah/complex-sz5-node0-4       48.0B ± 0%     32.0B ± 0%  -33.33%         (p=0.000 n=3+10)
Blah/complex-sz5-node1-4       24.0B ± 0%     24.0B ± 0%     ~     (all samples are equal)
Blah/complex-sz5-node2-4       48.0B ± 0%     32.0B ± 0%  -33.33%         (p=0.000 n=3+10)
Blah/complex-sz10-node0-4       112B ± 0%       64B ± 0%  -42.86%         (p=0.000 n=3+10)
Blah/complex-sz10-node1-4       120B ± 0%       40B ± 0%  -66.67%         (p=0.000 n=3+10)
Blah/complex-sz10-node2-4      56.0B ± 0%     40.0B ± 0%  -28.57%         (p=0.000 n=3+10)
Blah/complex-sz100-node0-4    1.00kB ± 0%    0.46kB ± 0%  -54.40%         (p=0.000 n=3+10)
Blah/complex-sz100-node1-4    1.02kB ± 0%    0.49kB ± 0%  -51.97%         (p=0.000 n=3+10)
Blah/complex-sz100-node2-4    1.18kB ± 0%    0.44kB ± 0%  -62.59%         (p=0.000 n=3+10)

name                        old allocs/op  new allocs/op  delta
Blah/simple-sz1-node0-4        0.00 ±NaN%     0.00 ±NaN%     ~     (all samples are equal)
Blah/simple-sz5-node0-4         2.00 ± 0%      2.00 ± 0%     ~     (all samples are equal)
Blah/simple-sz5-node1-4         3.00 ± 0%      3.00 ± 0%     ~     (all samples are equal)
Blah/simple-sz5-node2-4         2.00 ± 0%      2.00 ± 0%     ~     (all samples are equal)
Blah/simple-sz10-node0-4        5.00 ± 0%      5.00 ± 0%     ~     (all samples are equal)
Blah/simple-sz10-node1-4        5.00 ± 0%      5.00 ± 0%     ~     (all samples are equal)
Blah/simple-sz10-node2-4        6.00 ± 0%      6.00 ± 0%     ~     (all samples are equal)
Blah/simple-sz100-node0-4       55.0 ± 0%      55.0 ± 0%     ~     (all samples are equal)
Blah/simple-sz100-node1-4       54.0 ± 0%      54.0 ± 0%     ~     (all samples are equal)
Blah/simple-sz100-node2-4       59.0 ± 0%      59.0 ± 0%     ~     (all samples are equal)
Blah/complex-sz1-node0-4        1.00 ± 0%      1.00 ± 0%     ~     (all samples are equal)
Blah/complex-sz5-node0-4        5.00 ± 0%      4.00 ± 0%  -20.00%         (p=0.000 n=3+10)
Blah/complex-sz5-node1-4        3.00 ± 0%      3.00 ± 0%     ~     (all samples are equal)
Blah/complex-sz5-node2-4        5.00 ± 0%      4.00 ± 0%  -20.00%         (p=0.000 n=3+10)
Blah/complex-sz10-node0-4       10.0 ± 0%       8.0 ± 0%  -20.00%         (p=0.000 n=3+10)
Blah/complex-sz10-node1-4       7.00 ± 0%      5.00 ± 0%  -28.57%         (p=0.000 n=3+10)
Blah/complex-sz10-node2-4       6.00 ± 0%      5.00 ± 0%  -16.67%         (p=0.000 n=3+10)
Blah/complex-sz100-node0-4      75.0 ± 0%      57.0 ± 0%  -24.00%         (p=0.000 n=3+10)
Blah/complex-sz100-node1-4      79.0 ± 0%      61.0 ± 0%  -22.78%         (p=0.000 n=3+10)
Blah/complex-sz100-node2-4      78.0 ± 0%      55.0 ± 0%  -29.49%         (p=0.000 n=3+10)
```
