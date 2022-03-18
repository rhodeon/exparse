package main

import (
	"reflect"
	"testing"
)

func Test_validate(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantErr    error
	}{
		{"valid expression", "2 + 5 + 6", nil},
		{"illegal character", "3 + b - 2", ErrIllegalCharacter},
		{"legal start (digit)", "2/3", nil},
		{"legal start (parenthesis)", "(2/3)", nil},
		{"legal start (subtraction)", "-2/3", nil},
		{"legal start (whitespace)", " -2/3", nil},
		{"illegal start (addition)", "+2/3", ErrIllegalStart},
		{"illegal start (multiplication)", "*2/3", ErrIllegalStart},
		{"illegal start (division)", "/2/3", ErrIllegalStart},
		{"illegal start (decimal)", ".2/3", ErrIllegalStart},
		{"legal start (digit)", "2/3", nil},
		{"legal start (parenthesis)", "(2/3)", nil},
		{"illegal end (addition)", "2/3+", ErrIllegalEnd},
		{"illegal end (subtraction)", "2/3-", ErrIllegalEnd},
		{"legal end (whitespace)", "-2/3 ", nil},
		{"illegal end (multiplication)", "2/3*", ErrIllegalEnd},
		{"illegal end (division)", "2/3/", ErrIllegalEnd},
		{"illegal end (decimal)", "2/3.", ErrIllegalEnd},
		{"empty parenthesis", "2/3 + ()", ErrEmptyParentheses},
		{"legal consecutive operator (double minus)", "4--2", nil},
		{"illegal consecutive operator (triple minus)", "4---2", ErrIllegalConsecutiveOperator},
		{"illegal consecutive operator (addition)", "4++2", ErrIllegalConsecutiveOperator},
		{"illegal consecutive operator (multiplication)", "4**2", ErrIllegalConsecutiveOperator},
		{"illegal consecutive operator (division)", "4//2", ErrIllegalConsecutiveOperator},
		{"illegal consecutive operator (decimal)", "4..2", ErrIllegalConsecutiveOperator},
		{"illegal consecutive operator (mixed)", "4/+2", ErrIllegalConsecutiveOperator},
		{"illegal consecutive operator (mixed)", "4*/2", ErrIllegalConsecutiveOperator},
		{"illegal consecutive operator (minus then plus)", "4-+2", ErrIllegalConsecutiveOperator},
		{"illegal consecutive operator (minus then multiplication)", "4-/2", ErrIllegalConsecutiveOperator},
		{"illegal consecutive operator (minus then division)", "4-*2", ErrIllegalConsecutiveOperator},
		{"illegal consecutive operator (minus then decimal)", "4-.2", ErrIllegalConsecutiveOperator},
		{"illegal consecutive operator (multiplication then decimal)", "4*.2", ErrIllegalConsecutiveOperator},
		{"illegal consecutive operator (division then decimal)", "4/.2", ErrIllegalConsecutiveOperator},
		{"illegal consecutive operator (decimal then division)", "4./2", ErrIllegalConsecutiveOperator},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validate(tt.expression)

			if err != tt.wantErr {
				t.Errorf("\nGot:\t%v\nWant:\t%v", err, tt.wantErr)
			}
		})
	}
}

func Test_resolve(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantResult float64
		wantErr    error
	}{
		{"1", "2+5 - 2(6+4) + 3", -10, nil},
		{"2", "-2+5 - 2(6+4) + 3", -14, nil},
		{"3", "2+5 - 2(6*4) + 3", -38, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolve(tt.expression)

			if err != tt.wantErr {
				t.Errorf("\nGot:\t%v\nWant:\t%v", err, tt.wantErr)
			}

			if result != tt.wantResult {
				t.Errorf("\nGot:\t%f\nWant:\t%f", result, tt.wantResult)
			}
		})
	}
}

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
		{"pure addition", "3+5+8", []string{"3", "+", "5", "+", "8"}, nil},
		{"pure subtraction", "3-5-8", []string{"3", "-5", "-8"}, nil},
		{"pure multiplication", "3*5.5*8", []string{"132"}, nil},
		{"pure division", "3/5/8", []string{"0.075"}, nil},
		{"mixed", "3+5-2*7/3+4", []string{"3", "+", "5", "-4.666666666666667", "+", "4"}, nil},
		{"subtraction prefix", "-2*3", []string{"-", "6"}, nil},
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
		{"single positive operand", "4", 4},
		{"single negative operand", "-4", -4},
		{"addition only", "8+6", 14},
		{"subtraction only", "8-6", 2},
		{"addition prefix", "+2-5.2", -3.2},
		{"subtraction prefix", "-2+5.8", 3.8},
		{"subtraction prefix", "-2-5.8", -7.8},
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

func Test_normalize(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantResult []string
		wantErr    error
	}{
		{"addition", "4+3+2", []string{"4", "+", "3", "+", "2"}, nil},
		{"addition and subtraction", "4+3-2", []string{"4", "+", "3", "-2"}, nil},
		{"addition negation", "4+3+-2", []string{"4", "+", "3", "-2"}, nil},
		{"subtraction negation", "4+3--2", []string{"4", "+", "3", "+", "2"}, nil},
		{"combined multiplication and subtraction", "4+3*-2", []string{"4", "+", "3", "*", "-2"}, nil},
		{"combined division and subtraction", "4+3/-2", []string{"4", "+", "3", "/", "-2"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := normalize(tt.expression)

			if err != tt.wantErr {
				t.Errorf("\nGot:\t%v\nWant:\t%v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(result, tt.wantResult) {
				t.Errorf("\nGot:\t%#v\nWant:\t%#v", result, tt.wantResult)
			}
		})
	}
}

func Test_resolveParentheses(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantResult string
	}{
		{"1", "2+5-2(6+4)+3", "2+5-2*10+3"},
		{"2", "-2+5-2(6+4)+3", "-2+5-2*10+3"},
		{"3", "2+5-2(6*4)+3", "2+5-2*24+3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := resolveParentheses(tt.expression)

			if result != tt.wantResult {
				t.Errorf("\nGot:\t%s\nWant:\t%s", result, tt.wantResult)
			}
		})
	}
}
