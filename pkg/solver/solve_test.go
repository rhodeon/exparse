package solver

import (
	"reflect"
	"testing"
)

func TestSolve(t *testing.T) {
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
			result, err := Solve(tt.expression)

			if err != tt.wantErr {
				t.Errorf("\nGot:\t%v\nWant:\t%v", err, tt.wantErr)
			}

			if result != tt.wantResult {
				t.Errorf("\nGot:\t%f\nWant:\t%f", result, tt.wantResult)
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
