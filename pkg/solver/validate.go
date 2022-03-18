package solver

import (
	"regexp"
	"strconv"
)

// Validate checks if the given expression is well-formed.
func Validate(expr string) (string, error) {
	for _, char := range expr {
		value := string(char)
		if !isLegal(value) {
			return value, ErrIllegalCharacter
		}
	}

	// a maximum parentheses depth of 1 is allowed
	depth, err := maxDepth(expr)
	if err != nil {
		return "", err
	}

	if depth > 1 {
		return "", ErrDepthExceeded
	}

	// expressions must begin only with a digit, '(' or '-'
	startPattern := regexp.MustCompile("^[\\d\\-(\\s]")
	if !startPattern.MatchString(expr) {
		return "", ErrIllegalStart
	}

	// expressions must end only a digit or ')'
	endPattern := regexp.MustCompile("[\\d)\\s]$")
	if !endPattern.MatchString(expr) {
		return "", ErrIllegalEnd
	}

	// empty parentheses are not allowed
	parenthesisPattern := regexp.MustCompile("\\([)]")
	if parenthesisPattern.MatchString(expr) {
		return "", ErrEmptyParentheses
	}

	// only minus can appear consecutively after another operator
	operatorPattern := regexp.MustCompile("[+*/.]2+")
	if operatorPattern.MatchString(expr) {
		return "", ErrIllegalConsecutiveOperator
	}

	operatorPattern = regexp.MustCompile("-[+*/.]")
	if operatorPattern.MatchString(expr) {
		return "", ErrIllegalConsecutiveOperator
	}

	operatorPattern = regexp.MustCompile("---")
	if operatorPattern.MatchString(expr) {
		return "", ErrIllegalConsecutiveOperator
	}

	return "", nil
}

func isLegal(value string) bool {
	switch value {
	case add, subtract, multiply, divide, openParenthesis, closeParenthesis, decimal, whitespace:
		return true
	default:
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return false
		}
		return true
	}
}

// maxDepth returns the deepest level of nested brackets an expression has.
// An error is returned if the string doesn't have an equivalent closing bracket for each opening bracket.
func maxDepth(expression string) (int, error) {
	var currentDepth int
	var maxDepth int

	for _, char := range expression {
		if string(char) == closeParenthesis {
			currentDepth--
			if currentDepth < 0 {
				// more ")" than "("
				return -1, ErrMalformedExp
			}
		} else if string(char) == openParenthesis {
			currentDepth++
			if currentDepth > maxDepth {
				maxDepth = currentDepth
			}
		}
	}

	if currentDepth != 0 {
		// more "(" than ")"
		return -1, ErrMalformedExp
	}

	return maxDepth, nil
}
