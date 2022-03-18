package solver

import "testing"

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
			_, err := Validate(tt.expression)

			if err != tt.wantErr {
				t.Errorf("\nGot:\t%v\nWant:\t%v", err, tt.wantErr)
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
