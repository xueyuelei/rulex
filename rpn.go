package rulex

import (
	"fmt"
	"strings"
	"unicode"
)

var priorityTable = [6][6]rune{
	//        &    |    !	 (    )    0
	/* & */ {'>', '>', '<', '<', '>', '>'},
	/* | */ {'>', '>', '<', '<', '>', '>'},
	/* ! */ {'>', '>', '<', '<', '>', '>'},
	/* ( */ {'<', '<', '<', '<', '=', ' '},
	/* ) */ {' ', ' ', ' ', ' ', ' ', ' '},
	/* 0 */ {'<', '<', '<', '<', ' ', '='},
}

func getIndex(op rune) int {
	var idx int
	switch op {
	case '&':
		idx = 0
	case '|':
		idx = 1
	case '!':
		idx = 2
	case '(':
		idx = 3
	case ')':
		idx = 4
	case 0:
		idx = 5
	default:
		idx = -1
	}
	return idx
}

// ValidateOperandFunc is type of func to validate the name of operand
type ValidateOperandFunc func(string) bool

// RPN converts infix expression to Reverse Polish Notation expression
func RPN(expr string, fn ValidateOperandFunc) ([]string, error) {
	exprUTF8 := []rune(removeSpace(expr)) // compatible with utf-8
	exprUTF8 = append(exprUTF8, 0)        // sentinel 0

	opStk := NewStack()
	opStk.Push(rune(0))
	var exprRPN []string
	var begin int
	for !opStk.Empty() {
		end := getNext(exprUTF8, begin)

		seg := exprUTF8[begin:end]
		if len(seg) == 1 && getIndex(seg[0]) != -1 {
			op := seg[0]
			switch orderBetween(opStk.Top().(rune), op) {
			case '<':
				opStk.Push(op)
				begin = end
			case '>':
				op := opStk.Pop().(rune)
				exprRPN = append(exprRPN, string(op))
			case '=':
				opStk.Pop()
				begin = end
			default:
				return nil, fmt.Errorf("%w, no matching operand '%c', expr: %s", ErrInvalidSyntax, op, expr)
			}
		} else {
			condName := string(seg)
			if fn != nil && !fn(condName) {
				return nil, fmt.Errorf("%w, no condition name '%s', expr: %s", ErrCondNotMatch, condName, expr)
			}
			exprRPN = append(exprRPN, condName)
			begin = end
		}
	}
	if !validate(exprRPN) {
		return nil, fmt.Errorf("%w, mismatch between the number of operators and operands, expr: %s", ErrInvalidSyntax, expr)
	}

	return exprRPN, nil
}

func validate(rpn []string) bool {
	var left int
	for _, item := range rpn {
		if item == "!" {
		} else if item == "&" || item == "|" {
			left--
		} else {
			left++
		}
	}
	return left == 1
}

func removeSpace(s string) string {
	var builder strings.Builder
	for _, ch := range s {
		if unicode.IsSpace(ch) {
			continue
		}
		builder.WriteRune(ch)
	}
	return builder.String()
}

func orderBetween(left, right rune) rune {
	top_idx := getIndex(left)
	cur_idx := getIndex(right)

	return priorityTable[top_idx][cur_idx]
}

func getNext(expr []rune, idx int) int {
	for i, ch := range expr[idx:] {
		if getIndex(ch) != -1 {
			if i == 0 {
				return idx + i + 1
			}
			return idx + i
		}
	}
	return len(expr)
}
