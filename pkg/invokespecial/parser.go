package invokespecial

import (
	"encoding/json"
	"fmt"
	"strings"
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

type ParseContext struct {
	Text                   string
	Position, FailPosition int
	Errors                 []error
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

func assert[T any](pc ParseContext, condition bool, message string, len int, result T) (T, error) {
	if !condition {
		pc.Position += len
		return result, nil
	}
	if pc.Position > pc.FailPosition {
		pc.FailPosition = pc.Position
		pc.Errors = []error{}
	}
	if pc.Position == pc.FailPosition {
		pc.Errors = append(pc.Errors, fmt.Errorf(message))
	}
	return result, fmt.Errorf(message)
}

func getError(pc *ParseContext) (string, error) {
	prev := strings.Split(pc.Text[0:pc.FailPosition], "\n")
	begin := prev[len(prev)-1]
	row, col := len(prev), len(begin)

	sb := strings.Builder{}
	sb.WriteString(begin + "\n")

	for i := 0; i < len(begin)+1; i++ {
		sb.WriteString(".")
	}
	sb.WriteString(strings.Split(pc.Text, "\n")[row-1][0:len(begin)])

	b, err := json.MarshalIndent(pc.Text[pc.FailPosition], "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal error %w", err)
	}

	errs := make([]string, len(pc.Errors))
	for i, e := range pc.Errors {
		errs[i] = e.Error()
	}

	return fmt.Sprintf(
		"%d:%d: got %s, expected %s\n%s",
		row, col,
		string(b), strings.Join(errs, ","),
		sb.String(),
	), nil

}
