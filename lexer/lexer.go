package lexer

import (
	"fmt"
	"regexp"
)

type regexHandler func(lex *Lexer, regex *regexp.Regexp)

type regexPattern struct {
	Regex   *regexp.Regexp
	handler regexHandler
}

type Lexer struct {
	patterns []regexPattern
	Tokens   []Token
	source   string
	pos      int
}

func (lexer *Lexer) advanceN(n int) {
	lexer.pos += n
}

func (lexer *Lexer) push(token Token) {
	lexer.Tokens = append(lexer.Tokens, token)
}

func (lexer *Lexer) at() byte {
	return lexer.source[lexer.pos]
}

func (lexer *Lexer) remainder() string {
	return lexer.source[lexer.pos:]
}

func (lexer *Lexer) atEOF() bool {
	return lexer.pos >= len(lexer.source)
}

func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *Lexer, regex *regexp.Regexp) {
		lex.advanceN(len(value))
		lex.push(NewToken(kind, value))
	}
}

func createLexer(source string) *Lexer {
	return &Lexer{
		pos:    0,
		source: source,
		Tokens: make([]Token, 0),
		patterns: []regexPattern{
			{
				Regex:   regexp.MustCompile(`\s+`),
				handler: skipHandler,
			},
			{
				Regex:   regexp.MustCompile(`[0-9]+(\.[0-9]+)?`),
				handler: numberTokenHandler,
			},
			{
				Regex:   regexp.MustCompile(`\[`),
				handler: defaultHandler(OPEN_BRACKET, "["),
			},
			{
				Regex:   regexp.MustCompile(`\]`),
				handler: defaultHandler(CLOSE_BRACKET, "]"),
			},
			{
				Regex:   regexp.MustCompile(`\{`),
				handler: defaultHandler(OPEN_CURLY, "{"),
			},
			{
				Regex:   regexp.MustCompile(`\}`),
				handler: defaultHandler(CLOSE_CURLY, "}"),
			},
			{
				Regex:   regexp.MustCompile(`\(`),
				handler: defaultHandler(OPEN_PAREN, "("),
			},
			{
				Regex:   regexp.MustCompile(`\)`),
				handler: defaultHandler(CLOSE_PAREN, ")"),
			},
			{
				Regex:   regexp.MustCompile(`==`),
				handler: defaultHandler(EQUALS, "=="),
			},
			{
				Regex:   regexp.MustCompile(`!=`),
				handler: defaultHandler(NOT_EQUALS, "!="),
			},
			{
				Regex:   regexp.MustCompile(`=`),
				handler: defaultHandler(ASSIGNMENT, "="),
			},
			{
				Regex:   regexp.MustCompile(`!`),
				handler: defaultHandler(NOT, "!"),
			},
			{
				Regex:   regexp.MustCompile(`<=`),
				handler: defaultHandler(LESS_EQUALS, "<="),
			},
			{
				Regex:   regexp.MustCompile(`<`),
				handler: defaultHandler(LESS, "<"),
			},
			{
				Regex:   regexp.MustCompile(`>=`),
				handler: defaultHandler(GREATER_EQUALS, ">="),
			},
			{
				Regex:   regexp.MustCompile(`>`),
				handler: defaultHandler(GREATER, ">"),
			},
			{
				Regex:   regexp.MustCompile(`\|\|`),
				handler: defaultHandler(OR, "||"),
			},
			{
				Regex:   regexp.MustCompile(`&&`),
				handler: defaultHandler(AND, "&&"),
			},
			{
				Regex:   regexp.MustCompile(`\.\.`),
				handler: defaultHandler(DOT_DOT, ".."),
			},
			{
				Regex:   regexp.MustCompile(`\.`),
				handler: defaultHandler(DOT, "."),
			},
			{
				Regex:   regexp.MustCompile(`;`),
				handler: defaultHandler(SEMI_COLON, ";"),
			},
			{
				Regex:   regexp.MustCompile(`:`),
				handler: defaultHandler(COLON, ":"),
			},
			{
				Regex:   regexp.MustCompile(`\?`),
				handler: defaultHandler(QUESTION, "?"),
			},
			{
				Regex:   regexp.MustCompile(`,`),
				handler: defaultHandler(COMMA, ","),
			},
			{
				Regex:   regexp.MustCompile(`\+\+`),
				handler: defaultHandler(PLUS_PLUS, "++"),
			},
			{
				Regex:   regexp.MustCompile(`\-\-`),
				handler: defaultHandler(MINUS_MINUS, "--"),
			},
			{
				Regex:   regexp.MustCompile(`\+=`),
				handler: defaultHandler(PLUS_EQUALS, "+="),
			},
			{
				Regex:   regexp.MustCompile(`\-=`),
				handler: defaultHandler(MINUS_EQUALS, "-="),
			},
			{
				Regex:   regexp.MustCompile(`\+`),
				handler: defaultHandler(PLUS, "+"),
			},
			{
				Regex:   regexp.MustCompile(`-`),
				handler: defaultHandler(MINUS, "-"),
			},
			{
				Regex:   regexp.MustCompile(`/`),
				handler: defaultHandler(SLASH, "/"),
			},
			{
				Regex:   regexp.MustCompile(`\*`),
				handler: defaultHandler(STAR, "*"),
			},
			{
				Regex:   regexp.MustCompile(`%`),
				handler: defaultHandler(PERCENT, "%"),
			},
			{
				Regex:   regexp.MustCompile(`\-\-`),
				handler: defaultHandler(MINUS_MINUS, "--"),
			},
			{
				Regex:   regexp.MustCompile(`\-\-`),
				handler: defaultHandler(MINUS_MINUS, "--"),
			},
			{
				Regex:   regexp.MustCompile(`\-\-`),
				handler: defaultHandler(MINUS_MINUS, "--"),
			},
			{
				Regex:   regexp.MustCompile(`\-\-`),
				handler: defaultHandler(MINUS_MINUS, "--"),
			},
		},
	}
}

func numberTokenHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(NewToken(NUMBER, match))
	lex.advanceN(len(match))
}

func skipHandler(lex *Lexer, regexp *regexp.Regexp) {
	match := regexp.FindStringIndex(lex.remainder())
	lex.advanceN(match[1])
}

func Tokenize(source string) []Token {
	lex := createLexer(source)
	for !lex.atEOF() {
		matched := false
		for _, pattern := range lex.patterns {
			location := pattern.Regex.FindStringIndex(lex.remainder())
			if location != nil && location[0] == 0 {
				pattern.handler(lex, pattern.Regex)
				matched = true
				break
			}
		}
		if !matched {
			panic(fmt.Sprintf("Lexer::Error -> unrecognized token near %s\n", lex.remainder()))
		}
	}
	lex.push(NewToken(EOF, "EOF"))
	return lex.Tokens
}
