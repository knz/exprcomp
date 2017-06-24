package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) == 0 {
		return
	}
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "usage: %s [gen N|use]\n", os.Args[0])
		os.Exit(1)
	}
	switch os.Args[1] {
	case "gen":
		max, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, "arg error: %v", err)
			os.Exit(1)
		}
		canfunc := true
		if len(os.Args) > 3 {
			canfunc = false
		}
		gen(max, canfunc, os.Stdout)
	case "use":
		node, err := parse(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "parse error:", err)
			os.Exit(1)
		}
		fmt.Println("parsed:", node)
		var ctx evalCtx
		res, err := node.eval(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, "eval error:", err)
			os.Exit(1)
		}
		fmt.Println("eval result:", int64(*res.(*dint)))

		prog := makeprog()
		prog.compile(node)
		fmt.Println("compiled prog:")
		prog.disas(os.Stdout)

		var st [100]datum
		err = prog.run(ctx, &st)
		if err != nil {
			fmt.Fprintln(os.Stderr, "exec error:", err)
			os.Exit(1)
		}
		fmt.Println("exec result:", st[0])

	default:
		fmt.Fprintln(os.Stderr, "unknown command: %s", os.Args[1])
		os.Exit(1)
	}
}
