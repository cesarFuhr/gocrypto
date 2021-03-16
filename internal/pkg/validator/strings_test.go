package validator

import (
	"errors"
	"regexp"
	"testing"
)

func TestLengthRule(t *testing.T) {
	type in struct {
		s        string
		min, max int
	}
	cases := []struct {
		d     string
		input in
		want  error
	}{
		{d: "returns an error if the string is shorter than min", input: in{"12", 3, 5}, want: ErrTooShort},
		{d: "returns an error if the string is longer than max", input: in{"121212", 3, 5}, want: ErrTooLong},
		{d: "returns nil if the string is inside the minimum length", input: in{"123", 3, 5}, want: nil},
		{d: "returns nil if the string is inside the maximum length", input: in{"12345", 3, 5}, want: nil},
	}

	for _, c := range cases {
		t.Run(c.d, func(t *testing.T) {
			rule := lengthRule{c.input.min, c.input.max}
			got := rule.Validate(c.input.s)

			if !errors.Is(got, c.want) {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}

func BenchmarkLengthSuccess(b *testing.B) {
	rule := lengthRule{1, 2}
	s := "1"
	for i := 0; i < b.N; i++ {
		rule.Validate(s)
	}
}
func BenchmarkLengthError(b *testing.B) {
	rule := lengthRule{1, 2}
	s := "1123"
	for i := 0; i < b.N; i++ {
		rule.Validate(s)
	}
}

func TestRequiredRule(t *testing.T) {
	cases := []struct {
		d     string
		input string
		want  error
	}{
		{d: "returns an error if the string is empty", input: "", want: ErrRequired},
		{d: "returns nil if the string is not empty", input: "1", want: nil},
	}

	rule := requiredRule{}
	for _, c := range cases {
		t.Run(c.d, func(t *testing.T) {
			got := rule.Validate(c.input)

			if !errors.Is(got, c.want) {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}

func BenchmarkRequiredSuccess(b *testing.B) {
	rule := requiredRule{}
	s := "1"
	for i := 0; i < b.N; i++ {
		rule.Validate(s)
	}
}

func BenchmarkRequiredError(b *testing.B) {
	rule := requiredRule{}
	s := ""
	for i := 0; i < b.N; i++ {
		rule.Validate(s)
	}
}

func TestRegexRule(t *testing.T) {
	type in struct {
		s  string
		rg *regexp.Regexp
	}
	rgxTest, _ := regexp.Compile("123")
	cases := []struct {
		d     string
		input in
		want  error
	}{
		{d: "returns an error if the regexp did not match", input: in{"321", rgxTest}, want: ErrRegexNotMatched},
		{d: "returns nil if the regexp have matched", input: in{"123", rgxTest}, want: nil},
	}

	for _, c := range cases {
		t.Run(c.d, func(t *testing.T) {
			rule := regexRule{c.input.rg}
			got := rule.Validate(c.input.s)

			if !errors.Is(got, c.want) {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}

func BenchmarkRegexpSuccess(b *testing.B) {
	rgx, _ := regexp.Compile("string")
	rule := regexRule{rgx}
	s := "string"
	for i := 0; i < b.N; i++ {
		rule.Validate(s)
	}
}

func BenchmarkRegexpError(b *testing.B) {
	rgx, _ := regexp.Compile("string")
	rule := regexRule{rgx}
	s := "123456"
	for i := 0; i < b.N; i++ {
		rule.Validate(s)
	}
}
