package validator

import "errors"

var (
	ErrTooShort        = errors.New("should not be shorter than")
	ErrTooLong         = errors.New("should not be longer than")
	ErrRequired        = errors.New("is required")
	ErrRegexNotMatched = errors.New("should satisfy the regex")
)
