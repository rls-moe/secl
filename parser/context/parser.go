package context

import (
	"errors"
)

func RegisterFunction(name string, fn SECLFunc) error {
	if _, ok := DefaultFunctions[name]; ok {
		return errors.New("function already registered")
	}
	DefaultFunctions[name] = fn
	return nil
}
func MustRegisterFunction(name string, fn SECLFunc) {
	err := RegisterFunction(name, fn)
	if err != nil {
		panic(err)
	}
}

type Parser struct {
	DisableNumericTypes bool // Disables Integer
	Functions           map[string]SECLFunc
	Keywords            map[string]TokenType
	Symbols             *TokenTypeRegistry
}

func NewParserContext() *Parser {
	p := &Parser{
		DisableNumericTypes: false,
		Functions:           map[string]SECLFunc{},
		Keywords:            map[string]TokenType{},
		Symbols:             GetDefaultTokenTypes(),
	}
	for k, v := range DefaultKeywords {
		p.Keywords[k] = v
	}
	for k, v := range DefaultFunctions {
		p.Functions[k] = v
		p.Keywords[k] = DefaultTokenTypes.Function
	}
	return p
}

func (p *Parser) ToPhase1() *ParserPhase1 {
	return &ParserPhase1{
		DisableNumericTypes: p.DisableNumericTypes,
		Symbols:             p.Symbols,
	}
}

func (p *Parser) ToPhase3() *ParserPhase3 {
	return &ParserPhase3{
		Symbols: p.Symbols,
	}
}

func (p *Parser) ToLexer() *Lexer {
	return &Lexer{
		Symbols:  p.Symbols,
		Keywords: p.Keywords,
	}
}

func (p *Parser) ToRuntime() *Runtime {
	return &Runtime{
		Functions:     p.Functions,
		ParserContext: p,
	}
}

type Lexer struct {
	Symbols  *TokenTypeRegistry
	Keywords map[string]TokenType
}

// Phase1 contains only data relevant to Phase1 parsing
type ParserPhase1 struct {
	DisableNumericTypes bool // Disables Integer
	Symbols             *TokenTypeRegistry
}

type ParserPhase2 struct {
	KeywordList []string
}

type ParserPhase3 struct {
	Symbols *TokenTypeRegistry
}

type Runtime struct {
	Functions     map[string]SECLFunc
	ParserContext *Parser
}
