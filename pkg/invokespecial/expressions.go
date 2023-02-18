package invokespecial

import "fmt"

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
func EOF() ParserFunc[rune] {
	return func(ctx *ParseContext) (rune, error) {
		if ctx.IsAtEnd() {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to parse eof")
	}
}

// Haskell: neg p = Parser $ \s -> maybe (Just ((), s)) (const Nothing) (runParser p s)
func Negate[T any](p ParserFunc[T]) ParserFunc[T] {
	return func(ctx *ParseContext) (T, error) {
		x, err := p(ctx)
		if err == nil {
			var nilValue T
			return nilValue, fmt.Errorf("failed to parse negate")
		}
		return x, nil
	}
}

// Haskell: stry :: Parser t a -> Parser t [t]
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

// Haskell: stry p = Parser $ \s -> fmap (\(_, s1) -> (take (length s - length s1) s, s1)) (runParser p s)
// Haskell: inter a b = (:) <$> a <*> many (b *> a)
// Haskell: dang a b = ((a `inter` b) <* optional b) <|> pure []
// Haskell: range a b = anyOf [a..b]
// Haskell: document = concat <$> many statement where
// Haskell: statement = func <|> value <|> text
// Haskell: func = gen <$> open <*> document <*> close where
//				open = template (stry (str "func " *> rest))
//				close = template (str "endfunc")
//				gen o b _ = concat $ map (++ "\n") [o ++ " {", "result := \"\"", b, "return result", "}"]
// Haskell: value = gen <$> template (str " " *> rest) where
//				gen s = "result += " ++ s ++ "\n"
// Haskell: rest = many (neg close *> dot)
// Haskell: template p = open *> p <* close
// Haskell: text = gen <$> some (neg open *> dot) where
//				gen s = "result += " ++ show s ++ "\n"
// open = str "{%"
// close = str "%}"
// space = many $ char ' '
