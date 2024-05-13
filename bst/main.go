package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type set[T comparable] map[T]bool

func (s set[T]) contains(e T) bool {
	return s[e]
}

func (s set[T]) add(e T) {
	s[e] = true
}

type traversal func() []int64

type node struct {
	v    int64
	l, r *node
}

type tree struct {
	root *node
}

func (t *tree) insert(i int64) {
	if t.root == nil {
		t.root = &node{
			v: i,
		}
		return
	}

	current := t.root
	for {
		if i > current.v {
			if current.r != nil {
				current = current.r
			} else {
				current.r = &node{
					v: i,
				}
				return
			}
		} else {
			if current.l != nil {
				current = current.l
			} else {
				current.l = &node{
					v: i,
				}
				return
			}
		}
	}
}

func (t tree) preorder() []int64 {
	var result []int64

	stack := []*node{t.root}
	for len(stack) > 0 {
		e := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		result = append(result, e.v)
		if e.r != nil {
			stack = append(stack, e.r)
		}
		if e.l != nil {
			stack = append(stack, e.l)
		}
	}

	return result
}

func (t tree) inorder() []int64 {
	var result []int64

	visited := set[int64]{}
	stack := []*node{t.root}
	for len(stack) > 0 {
		e := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited.contains(e.v) {
			result = append(result, e.v)
			if e.r != nil {
				stack = append(stack, e.r)
			}
			continue
		}

		visited.add(e.v)

		stack = append(stack, e)
		if e.l != nil {
			stack = append(stack, e.l)
		}
	}

	return result
}

func (t tree) postorder() []int64 {
	var result []int64

	visited := set[int64]{}
	stack := []*node{t.root}
	for len(stack) > 0 {
		e := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited.contains(e.v) {
			result = append(result, e.v)
			continue
		}

		visited.add(e.v)

		stack = append(stack, e)
		if e.r != nil {
			stack = append(stack, e.r)
		}
		if e.l != nil {
			stack = append(stack, e.l)
		}
	}

	return result
}

func (t tree) levelorder() []int64 {
	var result []int64

	q := []*node{t.root}
	for len(q) > 0 {
		e := q[0]
		q = q[1:]

		result = append(result, e.v)
		if e.l != nil {
			q = append(q, e.l)
		}
		if e.r != nil {
			q = append(q, e.r)
		}
	}

	return result
}

func main_io(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 1000000), 1000000)
	var inputs []string

	t := tree{}

	var n int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &n)

	scanner.Scan()
	inputs = strings.Split(scanner.Text(), " ")
	for i := 0; i < n; i++ {
		vi, _ := strconv.ParseInt(inputs[i], 10, 32)
		t.insert(vi)
	}

	traversals := []traversal{
		t.preorder,
		t.inorder,
		t.postorder,
		t.levelorder,
	}

	for _, trv := range traversals {
		values := trv()
		result := make([]string, len(values))
		for vi, v := range values {
			result[vi] = fmt.Sprint(v)
		}
		fmt.Fprintln(out, strings.Join(result, " "))
	}
}

func main() {
	main_io(os.Stdin, os.Stdout)
}
