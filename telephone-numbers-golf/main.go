package main

import "fmt"

type n struct {
	c map[byte]n
}

func (n n) l() int {
	s := len(n.c)
	for _, c := range n.c {
		s += c.l()
	}
	return s
}

func (a *n) i(p string) {
	if len(p) == 0 {
		return
	}
	d := p[0]
	v, ok := a.c[d]
	if ok {
		v.i(p[1:])
		return
	}
	t := n{map[byte]n{}}
	t.i(p[1:])
	a.c[d] = t
}

func main() {
	N := 0
	S := fmt.Scan
	S(&N)
	r := &n{map[byte]n{}}
	for i := 0; i < N; i++ {
		t := ""
		S(&t)
		r.i(t)
	}
	fmt.Println(r.l())
}
