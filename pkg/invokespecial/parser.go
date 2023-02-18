package invokespecial

import (
	"fmt"
)

// func Ok[T any](value T) (T, error) {
// 	return value, nil
// }
// func Fail[T any](msg string) (T, error) {
// 	var nilValue T
// 	return nilValue, errors.New(msg)
// }

type ParseContext struct {
	Text     string
	Position int
}

type Pair[T any, U any] struct {
	First  T
	Second U
}

func NewPair[T any, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{first, second}
}

type ParserFunc[T any] func(ctx *ParseContext) (T, error)

func Parse[T any](parser ParserFunc[T], text string) (T, error) {
	c := ParseContext{Text: text}
	r, err := parser(&c)
	var nilValue T
	if err != nil {
		return nilValue, fmt.Errorf("failed to parse %w", err)
	} else if c.Position != len(text) {
		return nilValue, fmt.Errorf("no errors but failed to parse entire string %w", err)
	}

	return r, nil
}
