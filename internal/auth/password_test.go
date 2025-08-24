package auth

import (
	"testing"
)

func testHashPassword(t *testing.T) {
	cases := []struct {
		input string
		expected string
	}{
		{
			input: "mypassword",
			expected: "$2a$10$2dvOknNA36o098qCcxpcpuK1Vx9ILImrwfTbENsHkXzDkg2gaBrP.",
		},
		{
			input: "maytheforcebewithyou",
			expected: "$2a$10$mqs3tQBeQ/qlUppvEqBdcOOddNm3Z4IgmjPjZSEpplxKnC3deEqo2",
		},
		{
			input: "supersecretpassword",
			expected: "$2a$10$ngeX5STNZcXEGoBAA91vDeadETgYR7eyvhKbGyEWvEI7xDBti9z3m",
		},
	}

	for _, c := range cases {
		actual, _ := HashPassword(c.input)
		if actual != c.expected {
			t.Errorf("Expected output %s, but got %s", c.expected, actual)
		}
	}
}

func testCheckPasswordHash(t *testing.T) {
	cases := []struct {
		input_password string
		input_hash string
		expected error
	}{
		{
			input_password: "mypassword",
			input_hash: "$2a$10$2dvOknNA36o098qCcxpcpuK1Vx9ILImrwfTbENsHkXzDkg2gaBrP.",
			expected: nil,
		},
		{
			input_password: "maytheforcebewithyou",
			input_hash: "$2a$10$mqs3tQBeQ/qlUppvEqBdcOOddNm3Z4IgmjPjZSEpplxKnC3deEqo2",
			expected: nil,
		},
	}

	for _, c := range cases {
		err := CheckPasswordHash(c.input_password, c.input_hash)
		if err != nil {
			t.Errorf("Expected no errors, but got %s", err.Error())
		}
	}
}
