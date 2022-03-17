package main

import (
	"errors"
	"fmt"
	"log"
)

var ErrMalformedExp = errors.New("malformed expression")

const (
	openParenthesis  = "("
	closeParenthesis = ")"
)

func main() {
	max, err := maxDepth("()()()")
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Printf("Max depth: %d", max)
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
