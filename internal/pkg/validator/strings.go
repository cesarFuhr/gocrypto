package validator

import (
	"fmt"
	"regexp"
)

type StringRule interface {
	Validate(string) error
}

type lengthRule struct {
	min int
	max int
}

func (r lengthRule) Validate(s string) error {
	if len(s) < r.min {
		return fmt.Errorf("%w %d", ErrTooShort, r.min)
	}
	if len(s) > r.max {
		return fmt.Errorf("%w %d", ErrTooLong, r.max)
	}

	return nil
}

type requiredRule struct{}

func (r requiredRule) Validate(s string) error {
	if len(s) == 0 {
		return ErrRequired
	}

	return nil
}

type regexRule struct {
	rg *regexp.Regexp
}

func (r regexRule) Validate(s string) error {
	if !r.rg.MatchString(s) {
		return fmt.Errorf("%w %v", ErrRegexNotMatched, r.rg)
	}

	return nil
}
