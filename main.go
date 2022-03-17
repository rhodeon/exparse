package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrMalformedExp = errors.New("malformed expression")
var ErrIllegalCharacter = errors.New("illegal character detected")

// operators
const (
	openParenthesis  = "("
	closeParenthesis = ")"
	add              = "+"
	subtract         = "-"
	multiply         = "*"
	divide           = "/"
	decimal          = "."
)

func main() {
	//parseParens("2+4+ (5+6)2 * 8")
	res, err := evaluate("+ 4/-3*2")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%#v", res)
	}
}

func isLegal(value string) bool {
	switch value {
	case add, subtract, multiply, divide:
		return true
	default:
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return false
		}
		return true
	}
}

// evaluate computes a simplified expression list (with only addition and subtraction operators)
// into a final float result.
func evaluate(expr string) (float64, error) {
	simplified, err := simplify(expr)
	if err != nil {
		return 0, err
	}

	var result float64

	// the flags determine the next operation
	var addFlag = true
	var subtractFlag bool

	for _, value := range simplified {
		switch value {
		case add:
			addFlag = true
			subtractFlag = false

		case subtract:
			subtractFlag = true
			addFlag = false

		default:
			valueDigit, _ := strconv.ParseFloat(value, 64)

			// add or subtract based on active flag
			if addFlag {
				result = result + valueDigit
			} else if subtractFlag {
				result = result - valueDigit
			}

			addFlag = false
			subtractFlag = false
		}
	}

	return result, nil
}

// simplify reduces an expression to have only addition and subtraction operators.
// An error is returned if a non-digit or non-operator character is found.
func simplify(expr string) ([]string, error) {
	// remove whitespaces and guard against illegal character placements
	legalExpr, err := normalize(expr)
	if err != nil {
		return []string{}, err
	}

	var simplified []string
	var skipFlag bool

	for pos, value := range legalExpr {
		switch value {
		// simplify multiplications and divisions by
		// evaluating the preceding and succeeding operands
		case multiply, divide:

			// all ParseFloats are guaranteed to work as illegal characters have
			// already been guarded against
			prev, _ := strconv.ParseFloat(simplified[len(simplified)-1], 64)
			next, _ := strconv.ParseFloat(legalExpr[pos+1], 64)

			// replace previously added value in legalExp with the multiplication/division evaluation
			if value == multiply {
				simplified[len(simplified)-1] = strconv.FormatFloat(prev*next, 'f', -1, 64)
			} else {
				simplified[len(simplified)-1] = strconv.FormatFloat(prev/next, 'f', -1, 64)
			}
			skipFlag = true

		default:
			// skip current operand if previous operation was multiplication or division
			// as it has already been evaluated
			if !skipFlag {
				simplified = append(simplified, value)
			}
			skipFlag = false
		}
	}

	return simplified, nil
}

// normalize removes all whitespaces from the expression and returns a list of the remaining valid values.
// An error is returned if a non-digit or non-operator character is found.
// An error is returned if the expression begins with a multiplication, division or decimal.
// An error is returned if the expression begins with an addition, subtraction, multiplication, division or decimal.
func normalize(expr string) ([]string, error) {
	// remove whitespaces in expression
	trimmedExpr := strings.Join(strings.Fields(expr), "")

	//resolve consecutive addition and subtraction operators
	replacer := strings.NewReplacer("+-", "-", "--", "+")
	trimmedExpr = replacer.Replace(trimmedExpr)

	var legalExp []string

	// TODO: Guard against consecutive operators

	// the operand is accumulated and added to legalExp when an operator is found
	var operand string

	for pos, char := range trimmedExpr {
		value := string(char)

		if !isLegal(value) && value != decimal {
			return []string{}, ErrIllegalCharacter
		}

		if pos == 0 {
			switch value {
			// allow only addition and subtraction as starting operators
			case add, subtract:
				legalExp = append(legalExp, value)

			case multiply, divide, decimal:
				return []string{}, ErrIllegalCharacter

			default:
				operand += value
			}
		} else if pos == len(trimmedExpr)-1 {
			switch value {
			case add, subtract, multiply, divide, decimal:
				// operators cannot end an expression
				return []string{}, ErrIllegalCharacter

			default:
				// accumulate final operand and append to legalExp
				operand += value
				legalExp = append(legalExp, operand)
			}
		} else {
			switch value {
			case add, multiply, divide:
				// append and reset operand
				legalExp = append(legalExp, operand)
				operand = ""
				// append operator
				legalExp = append(legalExp, value)

			case subtract:
				// group subtraction with the next operand as a unary operator

				if operand != "" {
					// a separate operand already precedes the subtraction sign
					legalExp = append(legalExp, operand)
				}

				operand = ""
				operand += value

			default:
				// accumulate operand
				operand += value
			}
		}
	}

	return legalExp, nil
}

// maxDepth returns the deepest level of nested brackets an expression has.
// An error is returned if the string doesn't have an equivalent closing bracket for each opening bracket.
func maxDepth(expression string) (int, error) {
	stack := bracketStack{}

	for _, char := range expression {
		if string(char) == closeParenthesis {
			stack.currentDepth--
			if stack.currentDepth < 0 {
				// more ")" than "("
				return -1, ErrMalformedExp
			}
		} else if string(char) == openParenthesis {
			stack.currentDepth++
			if stack.currentDepth > stack.maxDepth {
				stack.maxDepth = stack.currentDepth
			}
		}
	}

	if stack.currentDepth != 0 {
		// more "(" than ")"
		return -1, ErrMalformedExp
	}

	return stack.maxDepth, nil
}

type bracketStack struct {
	history      []rune
	currentDepth int
	maxDepth     int
}
