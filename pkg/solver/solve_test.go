package solver

import (
	"reflect"
	"testing"
)

func TestSolve(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantResult string
	}{
		{"1", "2+5 - 2(6+4) + 3", "-10"},
		{"2", "-2+5 - 2(6+4) + 3", "-14"},
		{"3", "2+5 - 2(6*4) + 3", "-38"},
		{"4", "50 * 4 + 9 * 3 - 6 * 40000000000000", "-239999999999773"},
		{"5", "2(3.54 * 2.00 -1000 /200) / (20 + 30 * 2)", "0.052000000000000005"},
		{"6", " 2(3.54 * 2.00 -1000 /200) (20 + 30 * 2)", "332.8"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Solve(tt.expression)
			if result != tt.wantResult {
				t.Errorf("\nGot:\t%s\nWant:\t%s", result, tt.wantResult)
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
			result := resolveParentheses(tt.expression)
			if result != tt.wantResult {
				t.Errorf("\nGot:\t%s\nWant:\t%s", result, tt.wantResult)
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
			result := evaluate(tt.expression)

			if result != tt.result {
				t.Errorf("\nGot:\t%f\nWant:\t%f", result, tt.result)
			}
		})
	}
}

func Test_simplify(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantResult []string
	}{
		{"pure addition", "3+5+8", []string{"3", "+", "5", "+", "8"}},
		{"pure subtraction", "3-5-8", []string{"3", "-5", "-8"}},
		{"pure multiplication", "3*5.5*8", []string{"132"}},
		{"pure division", "3/5/8", []string{"0.075"}},
		{"mixed", "3+5-2*7/3+4", []string{"3", "+", "5", "-4.666666666666667", "+", "4"}},
		{"subtraction prefix", "-2*3", []string{"-", "6"}},
		{"joined operands", "2*4", []string{"8"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := simplify(tt.expression)
			if !reflect.DeepEqual(result, tt.wantResult) {
				t.Errorf("\nGot:\t%#v\nWant:\t%#v", result, tt.wantResult)
			}
		})
	}
}

func Test_normalize(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantResult []string
	}{
		{"addition", "4+3+2", []string{"4", "+", "3", "+", "2"}},
		{"addition and subtraction", "4+3-2", []string{"4", "+", "3", "-2"}},
		{"addition negation", "4+3+-2", []string{"4", "+", "3", "-2"}},
		{"subtraction negation", "4+3--2", []string{"4", "+", "3", "+", "2"}},
		{"combined multiplication and subtraction", "4+3*-2", []string{"4", "+", "3", "*", "-2"}},
		{"combined division and subtraction", "4+3/-2", []string{"4", "+", "3", "/", "-2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalize(tt.expression)
			if !reflect.DeepEqual(result, tt.wantResult) {
				t.Errorf("\nGot:\t%#v\nWant:\t%#v", result, tt.wantResult)
			}
		})
	}
}
