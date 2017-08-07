# S-Expression Configuration Language (SECL)

[**Github**](https://github.com/rls-moe/secl)

## Motivation

SECL is based on S-Expressions but reorients itself to being usuable and readable
as a generic configuration language.

The syntax should favor readability over compactness.

SECL should offer some basic built in functions to help programs achieve modular configuration without limitations or additional code.

## Design Goals

The Lexer/Tokenizer should be able to perform it's task using only forward-looking input. This enables usage under memory-constrained systems streaming the data

The parser itself should operate state-free, it's operation mode should only depend on the current node, never on previous or parent nodes.

SECL should be primary used via queries that return nodes on the parsed AST. MapLists make it difficult to properly unmarshal values into traditional objects in many languages.

## Reference Implementation

The reference implementation uses a simple stream-based lexer with a 3-phase parser, it demonstrates
the way SECL is intended to be used and parsed.

## Unicode and UTF-8

The SECL reference implementation fully supports the entire Unicode spectrum. Whitespace is as defined by the Unicode spec and basically the only important part of it.

## Parantheses

The entire file is wrapped within a `()` block, virtually, before parsing.

This allows one to omit the parantheses for first-level values.

## Types

The following atomic types are defined in SECL:

* Nil
* Booleans
* Strings
* Number
    * Integer Number
    * Decimal Number
* Bytes
* MapList

Numbers have arbitrary precision but parsers may limit them to platform-available
types if not otherwise possible.

Single Word Strings are all strings that do not contain a whitespace character, reserved characters like `:` or `!` or similar and are not a reserved keyword or function like `false` or `empty`, they may omit the quote characters.

Bytes are only a result of using the `decb64` or `loadb` functions and cannot be generated from SECL itself.

A MapList is a mix of a list/array and a dictionary/map.

### Predefined Values

* true, on, yes, allow
    * Truthy Values for Boolean Types
* false, off, no, deny
    * Falsy Values for Boolean Types
* nil
    * Empty Key
* empty, nothing
    * Empty MapList
* maybe
    * Is Randomly Truthy or Falsy
    * Truthy has a 50.1% chance
* randstr, randstr32, randstr64, randstr128, randstr256
    * Is a 42, 32, 64, 128 or 256 character string
    * Characters from base64 character set
    * Each occurence is a different string

### Nil

Nil is a non-special value, it is simply another value a key or list value might be set to. A key which is set to nil exists so there is no need to use a special type to express that some key might exist be nil, in that case the key simply exists and it's value happens to be nil, it could also happen to have any other valu.

## MapLists

SECL doesn't have arrays, lists or maps, instead it has one unified type: MapList

A MapList is a map and a list under one structure. If no map is specified for a list,
it is stored as a maplist without allocated map.

Maps are how SECL groups elements into sections.

A valid map key is a String and the value may be any other type.

## Functions

SECL provides some limited functions for ease of use:

* `loadd` - Loads all files from a given folder with a specific file extension (optional)
* `loadv` - Parse the specified file as a SECL atom value
* `loadf` - Load a single file
* `loadb` - Load a specified file as binary data
* `decb64` - Decode a base 64 string to binary
* `nop` - Accepts no parameters, returns nil, useful for testing
* `env` - Loads the specified environment variable
* `merge` - Merge all values in the parameter list into one single map

To execute functions, wrap them in parantheses and prepend a `!` to signal
the parser should expand the given block via execution.

Function execution is done in a seperate step performed by the user / application.
This prevents accidentally recursing endlessly into a directory or function.

### Merge

Merge will put several specified maplists into one.

As a general rule, if a specific key was reachable in a maplist, it cannot be made unreachable by merging it with another map. This means keys will not be deleted or have the complexity of their type reduced.

The reduction of type complexity goes as follows: `nil < everything but maplist < maplist`

Any value can replace a nil value but only a maplist can replace, for example, an integer.

Maplists are recursively merged under these rules, where existing keys are overwritten by newer values according to complexity rules.

If the complexity rules are broken, the merge must be aborted and the merge operation returns an error

### Not Implemented Yet

The following functions are not implemented yet:

* `loadd`
* `merge`


### Examples

```
additional-configuration: !(loadd "./conf.d" ".secl")
external-value: !(loadv "someconf.s")
binary-data: !(loadb "do_i_look_like_a.jpeg")
```

## Example

```
// Comments
# Comments
; Comments
/* Comments */

bool-value-1: true

string-value: "Hello World"

another-string-value: @"Hello World
You can easily use multiline strings,
only a \" can terminate strings.
    Leading whitespace however is trimmed."

raw-string: "
This is a raw string that begins with a \n and
    does not trim leading whitespace
"

nested-values: (
    section1: (
        bool-value: off
        switch-on-thing: maybe // 50.1/49.9 on/off chance
        encryption-key: randstr // Random Encryption Key
    )
    section2: (
        string-value: Hello!
        another-string: "Hello World!"
        int-value: 18
        more-ints: -12229
        or-really-big: 13231489146946239457269234761329471341985725234
        floats: 12312.121312
        actually: "we don't have floats, it's all decimals with arbitrary precision"
        decimals: 112.1121413043402374610471260327361203745103462037
        caseProblems?: "Not here!"
        special_charactersInKeys?: "No problemo!"
        JustDon'tUseTheSpacebarInKeys: Alrighto?
    )
    section3: (
        key1: Hallo
        key2: Welt!
        "Lists can be freely mixed into a map"
        "Keys may only be one word and need a appended :"
    )

    section4: empty
)
```

## Tests

You will find several test files in the reference implementations repository in the folder `/tests/`

The folder `/tests/must-parse` contains a series testfiles with .secl ending that must be parsed without error and a .expt file that contains debug output that should be reproduced in order to ensure the implementations have equivalent ideas about the internal structure of the data

A folder `/tests/must-bug` is planned, it will contain .secl files that must produce a specific or generic error when being parsed.

## EBNF (almost) Specification

Due to EBNF lacking in expressing Unicode ranges properly, the EBNF spec will define some things by simply including the functions of the reference code to express the validity of a character/rune.

```go
// isNotNull checks if the given rune is equal to 0x00
func isNotNull(cr rune) bool {
	return cr != 0
}

// isValidChunkRune checks if a given character can begin a chunk
func isValidChunkRune(cr rune) bool {
	return isNotNull(cr) && isValidInChunk(cr) && !unicode.IsNumber(cr) && !(cr == '!')
}

// isValidInChunk checks if a given character is valid to be used within a chunk
func isValidInChunk(cr rune) bool {
	return isNotNull(cr) && (!unicode.IsSpace(cr) && !(cr == '@' ||
		cr == ':' ||
		cr == '(' ||
		cr == ')' ||
		cr == '"'))
}

// isValidDigit checks if a given character can start a number token
func isValidDigit(cr rune) bool {
	return isNotNull(cr) && (unicode.IsNumber(cr) || (cr == '-' ||
		cr == '+'))
}

// isValidInNumber checks if a given character is valid within a number
func isValidInNumber(cr rune) bool {
	return isNotNull(cr) && (unicode.IsNumber(cr) || (cr == '+' ||
		cr == '-' ||
		cr == 'e' ||
		cr == '.' ||
		cr == '*' ||
		cr == '^' ||
		unicode.In(cr, unicode.Hex_Digit)))
}

func isValidAsSecondDigit(cr rune) bool {
	return isNotNull(cr) && (unicode.IsNumber(cr) || isValidInNumber(cr) || (cr == 'x' ||
		cr == 'o' ||
		cr == 'b' ||
		cr == 'X' ||
		cr == 'O' ||
		cr == 'B'))
}

func isValidStringBorder(cr rune) bool {
	return cr == '"'
}

func isEscapeInString(cr rune) bool {
	return cr == '\\'
}

func isValidCommentStarter(cr rune) bool {
	return cr == '/' || cr == ';' || cr == '#'
}

```

```ebnf
Comment = ( isValidCommentStarter , { isNotNull } , "\n" | "/*" , isNotNull , "*/" )

Number = isValidDigit , isValidAsSecondDigit , { isValidInNumber }
Symbol = isValidChunkRune , { isValidInChunk }
String = isValidStringBorder , { isNotNull | isEscapeInString , "\"" } , isValidStringBorder

MapList = "(" , Values , ")"
KeyValue = String , ":" , String

Root = Values

Values = { MapList | KeyValue | String | Number | Symbol | Comment }
```

The above definition is sufficient to construct a lexer that produces a flat list of tokens without validation if those are correctly formatted (only really relevant for numbers).