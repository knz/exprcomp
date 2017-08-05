package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

var simpleExprs = map[int][]string{}
var complexExprs = map[int][]string{}

var simpleNodes = map[int][]node{}
var complexNodes = map[int][]node{}

const maxexprs = 3

var szs = []int{1, 5, 10, 100}

func init() {
	for _, n := range szs {
		var sexprs []string
		var snodes []node
		var cexprs []string
		var cnodes []node
		for i := 0; i < maxexprs; i++ {
			if i > 0 && n == 1 {
				continue
			}
			var buf bytes.Buffer
			gen(n, false, &buf)

			nd, err := parse(strings.NewReader(buf.String()))
			if err != nil {
				panic(err)
			}

			sexprs = append(sexprs, buf.String())
			snodes = append(snodes, nd)

			for {
				buf.Reset()
				gen(n, true, &buf)
				if strings.ContainsAny(buf.String(), "abcdef") {
					break
				}
			}
			nd, err = parse(strings.NewReader(buf.String()))
			if err != nil {
				panic(err)
			}
			cexprs = append(cexprs, buf.String())
			cnodes = append(cnodes, nd)
		}
		simpleExprs[n] = sexprs
		complexExprs[n] = cexprs
		simpleNodes[n] = snodes
		complexNodes[n] = cnodes
	}
	for ni, nodes := range []map[int][]node{simpleNodes, complexNodes} {
		name := "simple"
		if ni == 1 {
			name = "complex"
		}
		for _, sz := range szs {
			for i, n := range nodes[sz] {
				fmt.Fprintf(os.Stderr, "%s sz%d node%d: %s\n", name, sz, i, n)
			}
		}
	}
}

var count int
var prevsz int

func pick(b *testing.B, nodes map[int][]node, sz int) node {
	if b.N == 1 {
		if sz == prevsz {
			count++
		} else {
			prevsz = sz
			count = 0
		}
	}
	// fmt.Fprintf(os.Stderr, "sz %d N %d pick %d\n", len(nodes), b.N, count)
	return nodes[sz][count]
}

func bencheval(b *testing.B, n node) {
	var ctx evalCtx
	for i := 0; i < b.N; i++ {
		_, err := n.eval(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchcomprun(b *testing.B, n node) {
	prog := makeprog()
	prog.compile(n)
	var ctx evalCtx
	var st [100]datum
	for i := 0; i < b.N; i++ {
		err := prog.run(ctx, &st)
		if err != nil {
			b.Fatal(err)
		}
	}
}

const doRun = true

func BenchmarkBlah(b *testing.B) {
	benches := []struct {
		name  string
		nodes map[int][]node
	}{
		{"simple", simpleNodes},
		{"complex", complexNodes},
	}
	for _, bench := range benches {
		for _, sz := range szs {
			for count := 0; count < len(bench.nodes[sz]); count++ {
				b.Run(fmt.Sprintf("%s-sz%d-node%d", bench.name, sz, count), func(b *testing.B) {
					if doRun {
						benchcomprun(b, bench.nodes[sz][count])
					} else {
						bencheval(b, bench.nodes[sz][count])
					}
				})
			}
		}
	}
}
