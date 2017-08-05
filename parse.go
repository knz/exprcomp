package main

import (
	"fmt"
	"io"
	"strconv"
)

func readbyte(in io.Reader) (byte, error) {
	var b [1]byte
	for {
		l, err := in.Read(b[:])
		if err != nil {
			return 0, err
		}
		if l > 0 {
			break
		}
	}
	return b[0], nil
}

func push(st []node, sp int, val node) ([]node, int) {
	if len(st) < sp+1 {
		return append(st, val), sp + 1
	}
	st[sp] = val
	return st, sp + 1
}

func parse(in io.Reader) (node, error) {
	var iv []byte

	var st []node
	sp := 0

	for {
		b, err := readbyte(in)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch b {
		case '+':
			if sp < 2 {
				return nil, fmt.Errorf("expected 2+ vals: %+v", st[:sp])
			}
			x, y := st[sp-2], st[sp-1]
			sp -= 2
			st, sp = push(st, sp, &intadd{left: x, right: y, addfunc: addint})
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			iv = append(iv, b)
			for {
				d, err := readbyte(in)
				if err != nil {
					return nil, err
				}
				if d == '.' {
					break
				}
				iv = append(iv, d)
			}
			n, err := strconv.ParseInt(string(iv), 10, 64)
			if err != nil {
				return nil, err
			}
			iv = iv[:0]
			nv := dint(n)
			st, sp = push(st, sp, &nv)
		case 'a', 'b', 'c', 'd', 'e', 'f':
			nargs := int(b - 'a')
			if sp < nargs {
				return nil, fmt.Errorf("expected %d+ vals: %+v", nargs, st[:sp])
			}
			args := append([]node(nil), st[sp-nargs:sp]...)
			sp -= nargs
			st, sp = push(st, sp, &call{args: args, fn: pseudofunc, makeres: pseudoalloc})
		default:
			// ignore
		}
	}

	if sp != 1 {
		return nil, fmt.Errorf("remaining values: %+v", st[:sp])
	}
	return st[0], nil
}
