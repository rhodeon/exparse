package solver

import "errors"

// legal characters
const (
	openParenthesis  = "("
	closeParenthesis = ")"
	add              = "+"
	subtract         = "-"
	multiply         = "*"
	divide           = "/"
	decimal          = "."
	whitespace       = " "
)

// errors
var (
	ErrMalformedExp               = errors.New("malformed expression")
	ErrIllegalCharacter           = errors.New("illegal character detected")
	ErrDepthExceeded              = errors.New("maximum parenthesis depth exceeded")
	ErrIllegalStart               = errors.New("expressions must begin only with a digit, '(' or '-'")
	ErrIllegalEnd                 = errors.New("expressions must end only a digit or ')'")
	ErrEmptyParentheses           = errors.New("empty parentheses are not allowed")
	ErrIllegalConsecutiveOperator = errors.New("illegal consecutive operators detected")
)
