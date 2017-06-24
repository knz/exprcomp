package main

import (
	"fmt"
	"io"
	"math/rand"
)

func gen(max int, canfunc bool, w io.Writer) {
	s := 0
	for l := 0; l < max; {
		typ, arg := rand.Int(), rand.Int()
		switch typ % 3 {
		case 0:
			fmt.Fprintf(w, "%d.", arg%10)
			s++
			l++
		case 1:
			if s < 2 {
				continue
			}
			fmt.Fprintf(w, "+")
			s -= 1
			l++
		case 2:
			if !canfunc {
				continue
			}
			nargs := arg % 6
			if s < nargs {
				continue
			}
			fmt.Fprintf(w, "%c", byte('a'+nargs))
			s -= nargs - 1
			l++
		}
	}
	for ; s > 1; s-- {
		fmt.Fprintf(w, "+")
	}
	fmt.Fprintln(w)
}
