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
	result, err := resolveParentheses("2+5 - 2(6+4) + 3")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%#v", result)
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
	simplified = splitUnary(simplified)

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

// splitUnary separates all subtraction signs in the expression from their operands.
func splitUnary(expr []string) []string {
	var result []string

	for pos, value := range expr {
		if pos == 0 {
			// allow only the first operand with a negation retain its sign
			result = append(result, value)
		} else {
			if string(value[0]) == subtract {
				result = append(result, subtract)
				result = append(result, value[1:])
			} else {
				result = append(result, value)
			}
		}
	}

	return result
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

	// ensure only numeric single characters are parsed
	if len(trimmedExpr) == 1 {
		_, err := strconv.ParseFloat(trimmedExpr, 64)
		if err != nil {
			return []string{}, ErrIllegalCharacter
		}
		return []string{trimmedExpr}, nil
	}

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

// resolveParentheses parses an expression with parentheses and returns its simplified form.
func resolveParentheses(expr string) (string, error) {
	simplified, err := simplifyParentheses(expr)
	if err != nil {
		return "", err
	}

	evaluated, err := evaluateParentheses(simplified)
	if err != nil {
		return "", err
	} else {
		return evaluated, nil
	}
}

// simplifyParentheses groups the parentheses in the expression
// and returns a list containing the groups as individual expressions
// alongside the expressions not covered by parentheses.
func simplifyParentheses(expr string) ([]string, error) {
	var nested bool
	var result []string
	var buffer string

	expr = strings.Join(strings.Fields(expr), "")

	for pos, char := range expr {
		value := string(char)

		switch value {
		case openParenthesis:
			nested = true
			buffer = ""

		case closeParenthesis:
			// flush all values contained within the parentheses
			nested = false
			result = append(result, buffer)
			buffer = ""

		case add, subtract, multiply, divide:
			if !nested {
				if buffer != "" {
					// flush any existing preceding operand to the result
					result = append(result, buffer)
					buffer = ""
				}

				// add the operator to the result
				result = append(result, value)
			} else {
				// accumulate operators as part of sub-expressions
				buffer += value
			}

		default:
			if !isLegal(value) && value != decimal {
				return []string{}, ErrIllegalCharacter
			}

			// accumulate digits, nested or not
			buffer += value

			if !nested {
				if pos == len(expr)-1 {
					// flush the final contents of the buffer
					result = append(result, buffer)
				} else if string(expr[pos+1]) == openParenthesis {
					// look ahead to append a multiplication if no operator
					// is found before the next parenthesis
					result = append(result, buffer)
					result = append(result, multiply)
					buffer = ""
				}
			}
		}
	}

	return result, nil
}

// evaluateParentheses computes the value of each sub-expression in the give list of expressions.
// Operators are not evaluated.
// The result is returned as a string with the simplified expression.
func evaluateParentheses(exprs []string) (string, error) {
	var evaluated []string

	for _, expr := range exprs {
		switch expr {
		case add, subtract, multiply, divide:
			// do not evaluate single operators
			evaluated = append(evaluated, expr)

		default:
			evaluatedExpr, err := evaluate(expr)
			if err != nil {
				return "", err
			}
			evaluated = append(evaluated, strconv.FormatFloat(evaluatedExpr, 'f', -1, 64))
		}
	}

	result := strings.Join(evaluated, "")
	return result, nil
}

func trimParentheses(expr string) string {
	expr = strings.TrimLeft(expr, "(")
	expr = strings.TrimRight(expr, ")")
	return expr
}
