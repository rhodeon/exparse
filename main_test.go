package main

import (
	"reflect"
	"testing"
)

func Test_maxDepth(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantDepth  int
		wantErr    error
	}{
		{"valid expression", "", 0, nil},
		{"valid expression", "()", 1, nil},
		{"valid expression", "()()", 1, nil},
		{"valid expression", "(())", 2, nil},
		{"valid expression", "()((()))()(())", 3, nil},
		{"malformed expression", "(", -1, ErrMalformedExp},
		{"malformed expression", ")", -1, ErrMalformedExp},
		{"malformed expression", "(()", -1, ErrMalformedExp},
		{"malformed expression", "())", -1, ErrMalformedExp},
		{"valid expression", "()((()))()((())", -1, ErrMalformedExp},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			depth, err := maxDepth(tt.expression)

			if err != tt.wantErr {
				t.Errorf("\nGot:\t%v\nWant:\t%v", err, tt.wantErr)
			}

			if depth != tt.wantDepth {
				t.Errorf("\nGot:\t%d\nWant:\t%d", depth, tt.wantDepth)
			}
		})
	}
}

func Test_simplify(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantResult []string
		wantErr    error
	}{
		{"pure addition", "3 + 5 + 8", []string{"3", "+", "5", "+", "8"}, nil},
		{"pure subtraction", "3 - 5 - 8", []string{"3", "-", "5", "-", "8"}, nil},
		{"pure multiplication", "3 * 5.5 * 8", []string{"132"}, nil},
		{"pure division", "3 / 5 / 8", []string{"0.075"}, nil},
		{"mixed", "3 + 5 - 2 * 7 / 3 + 4", []string{"3", "+", "5", "-", "4.666666666666667", "+", "4"}, nil},
		{"illegal character", "3 + b - 2", []string{}, ErrIllegalCharacter},
		{"addition prefix", "+ 2 / 3", []string{"+", "0.6666666666666666"}, nil},
		{"subtraction prefix", "- 2 * 3", []string{"-", "6"}, nil},
		{"multiplication prefix", "* 2 * 3", []string{}, ErrIllegalCharacter},
		{"division prefix", "/ 2 * 3", []string{}, ErrIllegalCharacter},
		{"end with operation", "2 * 3 +", []string{}, ErrIllegalCharacter},
		{"joined operands", "2*4", []string{"8"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := simplify(tt.expression)

			if err != tt.wantErr {
				t.Errorf("\nGot:\t%v\nWant:\t%v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(result, tt.wantResult) {
				t.Errorf("\nGot:\t%#v\nWant:\t%#v", result, tt.wantResult)
			}
		})
	}
}

func Test_evaluate(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		result     float64
	}{
		{"addition only", "8 + 6", 14},
		{"subtraction only", "8 - 6", 2},
		{"addition prefix", "+ 2 - 5.2", -3.2},
		{"subtraction prefix", "- 2 + 5.8", 3.8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := evaluate(tt.expression)

			if result != tt.result {
				t.Errorf("\nGot:\t%f\nWant:\t%f", result, tt.result)
			}
		})
	}
}
