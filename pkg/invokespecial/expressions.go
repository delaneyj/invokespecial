package invokespecial

import "fmt"

func Char(expected rune) ParserFunc[rune] {
	return func(ctx *ParseContext) (rune, error) {
		r := rune(ctx.Text[ctx.Position])
		if r == expected {
			ctx.Position++
			return r, nil
		}
		return 0, fmt.Errorf("failed to parse char. expected: %c, actual: %c", expected, r)
	}
}

func Str(expected string) ParserFunc[string] {
	return func(ctx *ParseContext) (string, error) {
		rpos := ctx.Position + len(expected)
		actual := ctx.Text[ctx.Position:rpos]
		if actual == expected {
			ctx.Position = rpos
			return expected, nil
		}
		return "", fmt.Errorf("failed to parse string. expected: %s, actual: %s", expected, actual)
	}
}

func Dot() ParserFunc[rune] {
	return func(ctx *ParseContext) (rune, error) {
		r := rune(ctx.Text[ctx.Position])
		ctx.Position++
		return r, nil
	}
}

// Esp WAT

func AnyOf[T any](parsers ...ParserFunc[T]) ParserFunc[T] {
	return func(ctx *ParseContext) (T, error) {
		for _, p := range parsers {
			r, err := p(ctx)
			if err == nil {
				return r, nil
			}
		}
		var nilValue T
		return nilValue, fmt.Errorf("failed to parse any of")
	}
}

func NoneOf[T any](parsers ...ParserFunc[T]) ParserFunc[rune] {
	return func(ctx *ParseContext) (rune, error) {
		for _, p := range parsers {
			_, err := p(ctx)
			if err == nil {
				return 0, fmt.Errorf("failed to parse none of")
			}
		}
		return Dot()(ctx)
	}
}

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
