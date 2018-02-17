package context

var DefaultFunctions = map[string]SECLFunc{}

var DefaultTokenTypes = GetDefaultTokenTypes()

var DefaultKeywords = map[string]TokenType{
	"nil":        DefaultTokenTypes.Nil,
	"empty":      DefaultTokenTypes.Empty,
	"nothing":    DefaultTokenTypes.Empty,
	"true":       DefaultTokenTypes.Bool,
	"on":         DefaultTokenTypes.Bool,
	"yes":        DefaultTokenTypes.Bool,
	"allow":      DefaultTokenTypes.Bool,
	"false":      DefaultTokenTypes.Bool,
	"off":        DefaultTokenTypes.Bool,
	"no":         DefaultTokenTypes.Bool,
	"deny":       DefaultTokenTypes.Bool,
	"maybe":      DefaultTokenTypes.Bool,
	"randstr":    DefaultTokenTypes.Randstr,
	"randstr32":  DefaultTokenTypes.Randstr,
	"randstr64":  DefaultTokenTypes.Randstr,
	"randstr128": DefaultTokenTypes.Randstr,
	"randstr256": DefaultTokenTypes.Randstr,
}
