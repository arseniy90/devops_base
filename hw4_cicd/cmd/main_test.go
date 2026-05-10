package main

import "testing"

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Empty string", "", true},
		{"Single character", "a", true},
		{"Simple palindrome", "racecar", true},
		{"Numeric palindrome", "12321", true},
		{"Case insensitive", "RaceCar", true},
		{"With spaces and symbols", "A man, a plan, a canal: Panama", true},
		{"Not a palindrome", "race a car", false},
		{"Mixed alphanumeric", "a car1 n 1rac a", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPalindrome(tt.input)
			if result != tt.expected {
				t.Errorf("IsPalindrome(%q) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
