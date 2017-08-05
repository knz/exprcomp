package main

import (
	"fmt"
	"io"
)

type ins struct {
	typ  int
	imm  int
	imm2 int
	ret  int
}

type prog struct {
	cst  []datum
	code []ins
	data []datum
	funs []func(ctx evalCtx, ret datum, args []datum) error
}

func (p *prog) disas(w io.Writer) {
	fmt.Fprintf(w, "constants: %+v\n", p.cst)
	fmt.Fprintf(w, "data: %+v\n", p.data)
	fmt.Fprintf(w, "functions: %+v\n", p.funs)
	fmt.Fprintln(w, "code:")
	for _, ins := range p.code {
		switch ins.typ {
		case 0:
			fmt.Fprintln(w, "loadcst", ins.imm)
		case 1:
			fmt.Fprintln(w, "addint")
		case 2:
			fmt.Fprintf(w, "call %d[%d]\n", ins.imm2, ins.imm)
		}
	}
}

func makeprog() prog {
	return prog{
		cst:  make([]datum, 0, 10),
		code: make([]ins, 0, 10),
		data: make([]datum, 0, 10),
		funs: make([]func(ctx evalCtx, res datum, args []datum) error, 0, 10),
	}
}

func (p *prog) reset() {
	p.cst = p.cst[:0]
	p.code = p.code[:0]
	p.data = p.data[:0]
}

func (p *prog) pushcode(c ins) {
	p.code = append(p.code, c)
}

func (p *prog) pushcst(c datum) int {
	cur := p.cst
	curlen := len(cur)
	p.cst = append(cur, c)
	return curlen
}

func (p *prog) pushdata(c datum) int {
	cur := p.data
	curlen := len(cur)
	p.data = append(cur, c)
	return curlen
}

func (p *prog) pushfunc(f func(ctx evalCtx, res datum, args []datum) error) int {
	cur := p.funs
	curlen := len(cur)
	p.funs = append(cur, f)
	return curlen
}

func (p *prog) compile(n node) {
	switch v := n.(type) {
	case datum:
		csti := p.pushcst(v)
		p.pushcode(ins{typ: 0, imm: csti})
	case *intadd:
		p.compile(v.left)
		p.compile(v.right)
		rv := dint(0)
		r := p.pushdata(&rv)
		p.pushcode(ins{typ: 1, imm: 0, ret: r})
	case *call:
		for _, a := range v.args {
			p.compile(a)
		}
		rv := v.makeres()
		r := p.pushdata(rv)
		fidx := p.pushfunc(v.fn)
		p.pushcode(ins{typ: 2, imm: len(v.args), imm2: fidx, ret: r})
	}
}
