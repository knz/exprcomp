package main

import (
	"bytes"
	"fmt"
	"strconv"
)

type evalCtx struct {
	somedata [100]byte
}

type node interface {
	fmt.Stringer
	eval(ctx evalCtx) (datum, error)
}

type datum interface {
	fmt.Stringer
	isDatum()
}

type dint int64

func (*dint) isDatum()                        {}
func (d *dint) String() string                { return strconv.FormatInt(int64(*d), 10) }
func (d *dint) eval(_ evalCtx) (datum, error) { return d, nil }

type intadd struct {
	left, right node
	addfunc     func(ctx evalCtx, res, x, y datum) error
}

func (b *intadd) String() string {
	return "(" + b.left.String() + "+" + b.right.String() + ")"
}

func addint(ctx evalCtx, res, x, y datum) error {
	ri := res.(*dint)
	*ri = dint(int64(*(x.(*dint))) + int64(*(y.(*dint))))
	return nil
}

func (b *intadd) eval(ctx evalCtx) (datum, error) {
	x, err := b.left.eval(ctx)
	if err != nil {
		return nil, err
	}
	y, err := b.right.eval(ctx)
	if err != nil {
		return nil, err
	}
	res := dint(0)
	return &res, b.addfunc(ctx, &res, x, y)
}

type call struct {
	args    []node
	makeres func() datum
	fn      func(ctx evalCtx, res datum, args []datum) error
}

func pseudoalloc() datum {
	res := dint(0)
	return &res
}

func pseudofunc(ctx evalCtx, res datum, args []datum) error {
	ri := res.(*dint)
	*ri = dint(int64(len(args)))
	return nil
}

func (c *call) String() string {
	var buf bytes.Buffer
	buf.WriteByte('a' + byte(len(c.args)))
	buf.WriteByte('(')
	for i, a := range c.args {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(a.String())
	}
	buf.WriteByte(')')
	return buf.String()
}

func (c *call) eval(ctx evalCtx) (datum, error) {
	args := make([]datum, len(c.args))
	for i, v := range c.args {
		d, err := v.eval(ctx)
		if err != nil {
			return nil, err
		}
		args[i] = d
	}
	res := c.makeres()
	return res, c.fn(ctx, res, args)
}
