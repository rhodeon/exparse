package main

import "testing"

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
