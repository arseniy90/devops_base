package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func IsPalindrome(s string) bool {
	if len(s) == 0 {
		return true
	}

	left, right := 0, len(s)-1
	for left < right {
		lRune := rune(s[left])
		rRune := rune(s[right])

		if !isAlnum(lRune) {
			left++
		} else if !isAlnum(rRune) {
			right--
		} else if unicode.ToLower(lRune) == unicode.ToLower(rRune) {
			left++
			right--
		} else {
			return false
		}
	}
	return true
}

func isAlnum(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsDigit(char)
}

func main() {
	in := bufio.NewReader(os.Stdin)

	line, _ := in.ReadString('\n')
	line = strings.TrimSpace(line)
	if IsPalindrome(line) {
		fmt.Println("True")
		return
	}

	fmt.Println("False")
}
