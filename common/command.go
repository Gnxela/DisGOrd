package common

import (
	"strconv"
	"strings"
)

type token interface {
	matches(str string) bool
	value(str string) interface{}
}

type AbsoluteToken struct {
	Pattern string
}

func (t *AbsoluteToken) value(str string) interface{} {
	return str
}

func (t *AbsoluteToken) matches(str string) bool {
	return t.Pattern == str
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

type NumericalToken struct{}

func (t *NumericalToken) value(str string) interface{} {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

func (t *NumericalToken) matches(str string) bool {
	_, err := strconv.ParseInt(str, 10, 64)
	return err == nil
}

type StringToken struct{}

func (t *StringToken) value(str string) interface{} {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

func (t *StringToken) matches(str string) bool {
	_, err := strconv.ParseInt(str, 10, 64)
	return err == nil
}

type Sequence struct {
	tokens []token
}

type Lexer struct {
	sequences []*Sequence
}

//Seperate token and input index

// Should REALLY be using trees here, a FSM will work for now
func (l *Lexer) ParseCommand(input string) (int, []interface{}) {
	splitInput := strings.Split(input, " ")
L:
	for i, sequence := range l.sequences {
		values := make([]interface{}, len(sequence.tokens))
		var inputIndex int
		var tokenIndex int
		for tokenIndex < len(sequence.tokens) {
			if inputIndex >= len(splitInput) {
				continue L
			}
			rawToken := sequence.tokens[tokenIndex]
			switch token := rawToken.(type) {
			case *OptionalToken:
				if !token.matches(splitInput[inputIndex]) {
					inputIndex-- //Don't consume input
				} else {
					values[tokenIndex] = token.value(splitInput[inputIndex])
				}
			default:
				if !token.matches(splitInput[inputIndex]) {
					continue L
				}
				values[tokenIndex] = token.value(splitInput[inputIndex])
			}
			inputIndex++
			tokenIndex++
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
