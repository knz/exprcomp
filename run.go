package main

import "fmt"

func (p *prog) run(ctx evalCtx, lst *[100]datum) error {
	code := p.code
	clen := len(code)
	cst := p.cst
	cstlen := len(cst)
	dst := 0
	data := p.data
	st := lst
	var res datum
	var err error
	for csp := 0; csp < clen; csp++ {
		ins := code[csp]
		switch ins.typ {
		case 0: // loadcst
			if ins.imm > cstlen {
				return fmt.Errorf("no such constant: %d", ins.imm)
			}
			res = cst[ins.imm]
		case 1: // add
			x, y := st[dst-2], st[dst-1]
			res := data[ins.ret]
			err = addint(ctx, res, x, y)
			if err != nil {
				return err
			}
			dst -= 2
		case 2: // call
			args := st[dst-ins.imm : dst]
			res = data[ins.ret]
			err = p.funs[ins.imm2](ctx, res, args)
			if err != nil {
				return err
			}
			dst -= ins.imm
		}
		st[dst] = res
		dst++
	}
	if dst != 1 {
		return fmt.Errorf("left over values: %+v", st)
	}
	return nil
}
