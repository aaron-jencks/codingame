package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostFix(t *testing.T) {
	testCases := []struct {
		name    string
		infix   string
		postfix []string
	}{
		{
			name:    "normal simple",
			infix:   "3*2+5",
			postfix: []string{"3", "2", "*", "5", "+"},
		},
		{
			name:    "normal parentheses",
			infix:   "3*(2+5*8-2)",
			postfix: []string{"3", "2", "5", "8", "*", "+", "2", "-", "*"},
		},
		{
			name:    "dice",
			infix:   "d6>d6-7",
			postfix: []string{"d6", "d6", "7", "-", ">"},
		},
		{
			name:    "dice parentheses",
			infix:   "(d6>d6)-7",
			postfix: []string{"d6", "d6", ">", "7", "-"},
		},
		{
			name:    "nested parentheses",
			infix:   "((d6+(d20-4)*8+9)>(d6+2))-7",
			postfix: []string{"d6", "d20", "4", "-", "8", "*", "+", "9", "+", "d6", "2", "+", ">", "7", "-"},
		},
	}

	for ti, tt := range testCases {
		t.Run(fmt.Sprintf("%d: %s", ti, tt.name), func(t *testing.T) {
			postfix := convertToPostfix(tt.infix)
			for pi, p := range postfix {
				assert.Equal(t, tt.postfix[pi], p, "postfix should be converted correctly")
			}
		})
	}
}

func TestOperandParse(t *testing.T) {
	testCases := []struct {
		name string
		op   string
		v    int
		die  bool
	}{
		{
			name: "int",
			op:   "567",
			v:    567,
			die:  false,
		},
		{
			name: "neg int",
			op:   "-567",
			v:    -567,
			die:  false,
		},
		{
			name: "die",
			op:   "d654",
			v:    654,
			die:  true,
		},
	}

	for ti, tt := range testCases {
		t.Run(fmt.Sprintf("%d: %s", ti, tt.name), func(t *testing.T) {
			operand := parseOperand(tt.op)
			assert.Equal(t, tt.v, operand.value, "values should match")
			assert.Equal(t, tt.die, operand.die, "die values should match")
		})
	}
}

func TestParseExpression(t *testing.T) {
	testCases := []struct {
		name       string
		expression []string
		result     map[int]int
	}{
		{
			name:       "2 d6s",
			expression: []string{"d6", "d6", "+"},
			result: map[int]int{
				2:  1,
				3:  2,
				4:  3,
				5:  4,
				6:  5,
				7:  6,
				8:  5,
				9:  4,
				10: 3,
				11: 2,
				12: 1,
			},
		},
		{
			name:       "d6*d6",
			expression: []string{"d6", "d6", "*"},
			result: map[int]int{
				1:  1,
				2:  2,
				3:  2,
				4:  3,
				5:  2,
				6:  4,
				8:  2,
				9:  1,
				10: 2,
				12: 4,
				15: 2,
				16: 1,
				18: 2,
				20: 2,
				24: 2,
				25: 1,
				30: 2,
				36: 1,
			},
		},
		{
			name:       "sample",
			expression: []string{"1", "d4", "+", "1", "+"},
			result: map[int]int{
				3: 1,
				4: 1,
				5: 1,
				6: 1,
			},
		},
	}

	for ti, tt := range testCases {
		t.Run(fmt.Sprintf("%d: %s", ti, tt.name), func(t *testing.T) {
			result := parseProbabilities(parseExpression(tt.expression))
			assert.Equal(t, len(tt.result), len(result), "number of elements in each results should match")
			for k, ev := range tt.result {
				v, ok := result[k]
				assert.True(t, ok, "keys should exist in result")
				assert.Equal(t, ev, v, "Values for %d should be equal", k)
			}
		})
	}
}

func TestParseFull(t *testing.T) {
	testCases := []struct {
		name   string
		expr   string
		result map[int]float64
	}{
		{
			name: "simple comparison",
			expr: "(2>5)+2*(5>2)+4*(10>5)",
			result: map[int]float64{
				6: 100,
			},
		},
	}

	for ti, tt := range testCases {
		t.Run(fmt.Sprintf("%d: %s", ti, tt.name), func(t *testing.T) {
			result := parseProbabilities(parseExpression(convertToPostfix(tt.expr)))
			assert.Equal(t, len(tt.result), len(result), "number of elements in each results should match")
			for k, ev := range tt.result {
				v, ok := result[k]
				assert.True(t, ok, "keys should exist in result")
				assert.Equal(t, ev, v, "Values for %d should be equal", k)
			}
		})
	}
}
