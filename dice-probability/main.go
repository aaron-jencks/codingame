package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type stack struct {
	buff []interface{}
}

func (s *stack) Push(e interface{}) {
	s.buff = append(s.buff, e)
}

func (s *stack) Pop() interface{} {
	if len(s.buff) == 0 {
		return ""
	} else {
		e := s.buff[len(s.buff)-1]
		s.buff = s.buff[:len(s.buff)-1]
		return e
	}
}

func (s *stack) Peek() interface{} {
	if len(s.buff) == 0 {
		return ""
	}
	return s.buff[len(s.buff)-1]
}

func (s stack) Len() int {
	return len(s.buff)
}

type operand struct {
	value int
	die   bool
}

func (o operand) ToMap() map[int]int {
	if !o.die {
		return map[int]int{
			o.value: 1,
		}
	}

	result := map[int]int{}
	for i := 1; i <= o.value; i++ {
		result[i] = 1
	}
	return result
}

func parseOperand(op string) operand {
	if op[0] == 'd' {
		var sides int
		fmt.Sscan(op[1:], &sides)
		return operand{sides, true}
	}
	var value int
	fmt.Sscan(op, &value)
	return operand{value, false}
}

func convertToPostfix(expression string) []string {
	var result []string

	var operatorStack *stack = &stack{}

	operandStart := 0
	for ei, ec := range expression {
		if strings.ContainsRune("*+->()", ec) {
			// We found an operator

			if operandStart != ei {
				result = append(result, expression[operandStart:ei])
			}
			operandStart = ei + 1

			switch ec {
			case '(':
				operatorStack.Push("(")
			case ')':
				for operatorStack.Peek() != "(" {
					result = append(result, operatorStack.Pop().(string))
				}
				operatorStack.Pop() // pop off the opening bracket
			case '*':
				for operatorStack.Peek() == "*" {
					result = append(result, operatorStack.Pop().(string))
				}
				operatorStack.Push("*")
			case '+':
				fallthrough
			case '-':
				for strings.ContainsAny(operatorStack.Peek().(string), "+-*") {
					result = append(result, operatorStack.Pop().(string))
				}
				operatorStack.Push(fmt.Sprintf("%c", ec))
			case '>':
				for operatorStack.Len() > 0 && operatorStack.Peek().(string) != "(" {
					result = append(result, operatorStack.Pop().(string))
				}
				operatorStack.Push(">")
			}
		}
	}

	if operandStart < len(expression) {
		result = append(result, expression[operandStart:])
	}

	for operatorStack.Len() > 0 {
		result = append(result, operatorStack.Pop().(string))
	}

	return result
}

func parseExpression(expression []string) map[int]int {
	opStack := stack{}
	for _, exp := range expression {
		if len(exp) == 1 && strings.ContainsAny(exp, "*+->") {
			b := opStack.Pop().(map[int]int)
			a := opStack.Pop().(map[int]int)
			nm := map[int]int{}
			switch exp {
			case "*":
				for ka, va := range a {
					for kb, vb := range b {
						c := ka * kb
						_, ok := nm[c]
						if !ok {
							nm[c] = 0
						}
						nm[c] += va * vb
					}
				}
			case "+":
				for ka, va := range a {
					for kb, vb := range b {
						c := ka + kb
						_, ok := nm[c]
						if !ok {
							nm[c] = 0
						}
						nm[c] += va * vb
					}
				}
			case "-":
				for ka, va := range a {
					for kb, vb := range b {
						c := ka - kb
						_, ok := nm[c]
						if !ok {
							nm[c] = 0
						}
						nm[c] += va * vb
					}
				}
			case ">":
				for ka, va := range a {
					for kb, vb := range b {
						var c int = 0
						if ka > kb {
							c = 1
						}

						_, ok := nm[c]
						if !ok {
							nm[c] = 0
						}

						nm[c] += va * vb
					}
				}
			}
			opStack.Push(nm)
		} else {
			opStack.Push(parseOperand(exp).ToMap())
		}
	}

	return opStack.Pop().(map[int]int)
}

func parseProbabilities(counts map[int]int) map[int]float64 {
	nm := map[int]float64{}
	total := 0
	for _, v := range counts {
		total += v
	}
	for k, v := range counts {
		nm[k] = float64(v) / float64(total) * 100
	}
	return nm
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	scanner.Scan()
	expr := scanner.Text()

	probs := parseProbabilities(
		parseExpression(
			convertToPostfix(expr)))

	var keys []int
	for k := range probs {
		keys = append(keys, k)
	}

	skeys := sort.IntSlice(keys)
	skeys.Sort()

	for _, k := range skeys {
		fmt.Printf("%d %.2f\n", k, probs[k])
	}
}
