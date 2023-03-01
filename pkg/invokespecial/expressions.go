package invokespecial

import "fmt"

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

// Haskell:  char c = satisfy (== c)
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

// Haskell: str s = Parser $ \t -> if s `isPrefixOf` t then Just (s, drop (length s) t) else Nothing
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

// Haskell: dot = satisfy $ const True
func Dot() ParserFunc[rune] {
	return func(ctx *ParseContext) (rune, error) {
		r := rune(ctx.Text[ctx.Position])
		ctx.Position++
		return r, nil
	}
}

// Haskell: eps = Parser $ \s -> Just ((), s)
func Eps() ParserFunc[interface{}] {
	return func(ctx *ParseContext) (interface{}, error) {
		return nil, nil
	}
}

// Haskell: anyOf s = satisfy (`elem` s)
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

// Haskell: notOf s = satisfy (`notElem` s)
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

// Haskell: eof = Parser $ \s -> if null s then Just ((), s) else Nothing
func EOF() ParserFunc[interface{}] {
	return func(ctx *ParseContext) (interface{}, error) {
		if ctx.IsAtEnd() {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to parse eof")
	}
}

// Haskell: neg p = Parser $ \s -> maybe (Just ((), s)) (const Nothing) (runParser p s)
func Negate[T any](p ParserFunc[T]) ParserFunc[interface{}] {
	return func(ctx *ParseContext) (interface{}, error) {
		if _, err := p(ctx); err == nil {
			return nil, fmt.Errorf("failed to parse negate")
		}
		return nil, nil
	}
}

// Haskell:  stry p = Parser $ \s -> fmap (\(_, s1) -> (take (length s - length s1) s, s1)) (runParser p s)
func Stry[T any](p ParserFunc[T]) ParserFunc[string] {
	return func(ctx *ParseContext) (string, error) {
		pos := ctx.Position
		_, err := p(ctx)
		if err != nil {
			return "", err
		}
		return ctx.Text[pos:ctx.Position], nil
	}
}

// Haskell: inter a b = (:) <$> a <*> many (b *> a)
func Inter[T any, U any](a ParserFunc[T], b ParserFunc[U]) ParserFunc[T] {
	return func(ctx *ParseContext) (T, error) {
		x, err := a(ctx)
		if err != nil {
			return x, err
		}
		for {
			_, err := b(ctx)
			if err != nil {
				break
			}
			x, err = a(ctx)
			if err != nil {
				break
			}
		}
		return x, nil
	}
}

// Haskell: dang a b = ((a `inter` b) <* optional b) <|> pure []
func Dangling[T any, U any](a ParserFunc[T], b ParserFunc[U]) ParserFunc[Pair[T, U]] {
	return Seq(Inter(a, b), Optional(b))
}

// Haskell: range a b = anyOf [a..b]
func Range(a, b rune) ParserFunc[rune] {
	return func(ctx *ParseContext) (rune, error) {
		r := rune(ctx.Text[ctx.Position])
		if r >= a && r <= b {
			ctx.Position++
			return r, nil
		}
		return 0, fmt.Errorf("failed to parse range. expected: %c-%c, actual: %c", a, b, r)
	}
}

func Many[T any](p ParserFunc[T]) ParserFunc[[]T] {
	return func(ctx *ParseContext) ([]T, error) {
		var result []T
		for {
			x, err := p(ctx)
			if err != nil {
				break
			}
			result = append(result, x)
		}
		return result, nil
	}
}

func Some[T any](p ParserFunc[T]) ParserFunc[[]T] {
	return func(ctx *ParseContext) ([]T, error) {
		x, err := p(ctx)
		if err != nil {
			return nil, err
		}
		result := []T{x}
		for {
			x, err := p(ctx)
			if err != nil {
				break
			}
			result = append(result, x)
		}
		return result, nil
	}
}

func Optional[T any](p ParserFunc[T]) ParserFunc[T] {
	return func(ctx *ParseContext) (T, error) {
		x, err := p(ctx)
		if err != nil {
			var nilValue T
			return nilValue, nil
		}
		return x, nil
	}
}

func Map[T any, U any](p ParserFunc[T], f func(T) U) ParserFunc[U] {
	return func(ctx *ParseContext) (U, error) {
		x, err := p(ctx)
		if err != nil {
			var nilValue U
			return nilValue, err
		}
		return f(x), nil
	}
}
