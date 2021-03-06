package common

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	stringTokenRegex *regexp.Regexp = regexp.MustCompile(`(?m)^(?:'(.*)'|"(.*)"|([^\s]*))`)
)

type token interface {
	matches(str string) bool
	value(str string) interface{}
	consume(str string) string
}

type AbsoluteToken struct {
	Pattern string
}

func (t *AbsoluteToken) value(str string) interface{} {
	return t.Pattern
}

func (t *AbsoluteToken) matches(str string) bool {
	if len(t.Pattern) > len(str) {
		return false
	}
	return t.Pattern == str[:len(t.Pattern)]
}

func (t *AbsoluteToken) consume(str string) string {
	return str[len(t.Pattern):]
}

type OptionalToken struct {
	Token token
}

func (t *OptionalToken) value(str string) interface{} {
	return t.Token.value(str)
}

func (t *OptionalToken) matches(str string) bool {
	return t.Token.matches(str)
}

func (t *OptionalToken) consume(str string) string {
	return t.Token.consume(str)
}

type NumericalToken struct{}

func (t *NumericalToken) value(str string) interface{} {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

func (t *NumericalToken) matches(str string) bool {
	_, err := strconv.ParseInt(str, 10, 64)
	return err == nil
}

func (t *NumericalToken) consume(str string) string {
	return str[len(string(t.value(str).(int64))):]
}

// Eventually should allow for quoting
type StringToken struct{}

func (t *StringToken) value(str string) interface{} {
	match := stringTokenRegex.FindStringSubmatch(str)
	if match[1] != "" {
		return match[1]
	} else if match[2] != "" {
		return match[2]
	} else {
		return match[3]
	}
}

func (t *StringToken) matches(str string) bool {
	return stringTokenRegex.MatchString(str)
}

func (t *StringToken) consume(str string) string {
	loc := stringTokenRegex.FindStringIndex(str)
	return str[loc[1]:]
}

type Sequence struct {
	tokens []token
}

type Lexer struct {
	sequences []*Sequence
}

// Should REALLY be using trees here, a FSM will work for now
func (l *Lexer) ParseCommand(command string) (int, []interface{}) {
L:
	for i, sequence := range l.sequences {
		input := command
		values := make([]interface{}, len(sequence.tokens))
		for tokenIndex, rawToken := range sequence.tokens {
			input = strings.TrimSpace(input) //Should not be doing this. Change later
			switch token := rawToken.(type) {
			case *OptionalToken:
				if token.matches(input) {
					values[tokenIndex] = token.value(input)
					input = token.consume(input)
				}
			default:
				if !token.matches(input) {
					continue L
				}
				values[tokenIndex] = token.value(input)
				input = token.consume(input)
			}
		}
		return i, values
	}
	return -1, nil
}

func CreateLexer(sequences ...*Sequence) *Lexer {
	return &Lexer{sequences}
}

func CreateSequence(tokens ...token) *Sequence {
	return &Sequence{tokens}
}
