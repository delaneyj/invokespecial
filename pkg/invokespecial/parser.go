package invokespecial

import (
	"fmt"
)

type ParserFunc[T any] func(ctx *ParseContext) (T, error)

func Parse[T any](parser ParserFunc[T], text string) (T, error) {
	pc := NewParseContext(text)
	r, err := parser(pc)
	var nilValue T
	if err != nil {
		return nilValue, fmt.Errorf("failed to parse %w", err)
	} else if !pc.IsAtEnd() {
		return nilValue, fmt.Errorf("no errors but failed to parse entire string %w", err)
	}

	return r, nil
}

// <*> in the Haskell world
func Seq[T any, U any](left ParserFunc[T], right ParserFunc[U]) ParserFunc[Pair[T, U]] {
	return func(ctx *ParseContext) (Pair[T, U], error) {
		l, err := left(ctx)
		if err != nil {
			var nilValue Pair[T, U]
			return nilValue, fmt.Errorf("failed to parse left value in seq %w", err)
		}
		r, err := right(ctx)
		if err != nil {
			var nilValue Pair[T, U]
			return nilValue, fmt.Errorf("failed to parse right value in seq %w", err)
		}
		return NewPair(l, r), nil
	}
}

type ParseContext struct {
	Text     string
	Position int
}

func NewParseContext(text string) *ParseContext {
	return &ParseContext{Text: text}
}

func (pc *ParseContext) IsAtEnd() bool {
	return pc.Position == len(pc.Text)
}

type Pair[T any, U any] struct {
	First  T
	Second U
}

func NewPair[T any, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{first, second}
}
