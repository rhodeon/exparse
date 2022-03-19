package solver

import (
	"strconv"
	"strings"
)

// Solve computes and returns the result of the given expression.
func Solve(expr string) string {
	// remove whitespaces in expression
	expr = strings.Join(strings.Fields(expr), "")
	resolved := resolveParentheses(expr)
	result := evaluate(resolved)
	return strconv.FormatFloat(result, 'f', -1, 64)

}

// resolveParentheses parses an expression with parentheses and returns its simplified form.
func resolveParentheses(expr string) string {
	simplified := simplifyParentheses(expr)
	evaluated := evaluateParentheses(simplified)
	return evaluated
}

// simplifyParentheses groups the parentheses in the expression
// and returns a list containing the groups as individual expressions
// alongside the expressions not covered by parentheses.
func simplifyParentheses(expr string) []string {
	var nested bool
	var result []string
	var buffer string

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

	return result
}

// evaluateParentheses computes the value of each sub-expression in the give list of expressions.
// Operators are not evaluated.
// The result is returned as a string with the simplified expression.
func evaluateParentheses(exprs []string) string {
	var evaluated []string

	for _, expr := range exprs {
		switch expr {
		case add, subtract, multiply, divide:
			// do not evaluate single operators
			evaluated = append(evaluated, expr)

		default:
			evaluatedExpr := evaluate(expr)
			evaluated = append(evaluated, strconv.FormatFloat(evaluatedExpr, 'f', -1, 64))
		}
	}

	result := strings.Join(evaluated, "")
	return result
}

// evaluate computes a simplified expression list (with only addition and subtraction operators)
// into a final float result.
func evaluate(expr string) float64 {
	// return single digits early
	if len(expr) == 1 {
		result, _ := strconv.ParseFloat(expr, 64)
		return result
	}

	simplified := simplify(expr)
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

	return result
}

// simplify reduces an expression to have only addition and subtraction operators.
// An error is returned if a non-digit or non-operator character is found.
func simplify(expr string) []string {
	legalExpr := normalize(expr)
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

	return simplified
}

// normalize returns a list of each element in the expression.
func normalize(expr string) []string {
	//Solve consecutive addition and subtraction operators
	replacer := strings.NewReplacer("+-", "-", "--", "+")
	expr = replacer.Replace(expr)
	var legalExp []string

	// the operand is accumulated and added to legalExp when an operator is found
	var operand string

	for pos, char := range expr {
		value := string(char)

		if pos == 0 {
			switch value {
			// allow only subtraction as a starting operator
			case subtract:
				legalExp = append(legalExp, value)

			default:
				operand += value
			}
		} else if pos == len(expr)-1 {
			// accumulate final operand and append to legalExp
			operand += value
			legalExp = append(legalExp, operand)
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

	return legalExp
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
