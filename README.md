# SECL

SECL is a S-Expression based configuration language, drawing inspiration and syntax from JSON and YAML.

It's main goal is to provide a configuration format that is easy to read and write with little overhead while
being simple and understandable in it's syntax and parsing.


```
list-entry

key: (
    another-list-entry
    key: "Hello World!"
    explicit: "UTF-8 support: 💩"
    
    other-features: (
        "arbitrary precision decimals": 29312345320423531230497123416079652.234234234
        "and integers": 343287640123745610476123012734613047236402374620376234023746
        "scientific notation": 9.1*10^22
        "multi-line strings": "
            HELLO WORLD!"
        "omit quotation around simple strings": Hello-Again
        "load other files from subfolders": !(loadd "folder" ".conf")
        "empty maps": empty
        
        "random values": (
            randstr128
            maybe
        )
        
        "mix lists and maps, they are the same in SECL"
    )

    new-features: (
        "loading environment variables": !(env "SOME_ENV_VAR")
        "with defaults": !(env "SOME_ENV_VAR" default: "Default Value")
        "merge": !(merge
            (you can now merge several maps)
            (the list part is appended to the previous map)
            (keys: are added)
            (a merge will never change a datatype, it will only overwrite existing keys if their types match)
        )
    )
)
```

## Project Status

SECL is considered in Version 0.9

Most features are now part of the reference implementation in this repository. The last missing function is `loadd`, merge has been implemented recently.

A non-formal specification can be found at [hackmd](https://hackmd.io/s/rylGmiXr-) or in this repository under SPECS.md (hackmd might be more up-to-date than SPECS.md)

## Introduction

### Strings

Strings that do not contain any whitespace (as recognized by UTF8) and do not start with a digit or contain reserved character can omit quotation marks;

```
HelloWorld

```

Reserved characters: `"!@:()`

If any of the reserved characters is present, the string must be wrapped with quotation marks. Additionally, if the string equals
a keyword or function name, it must be wrapped with quotation marks too. (ex. "false" or "empty" or "randstr256")

```
"Hello World"
```

Strings are multi-line by default, if you want to trim leading whitespace from each line and the beginning of the string, prepend a `@` symbol.

```
@"
Hello 
 World"
 
// equivalent to

"Hello
World"
```

The special values `randstr` and `randstr32` to `randstr256` allow the usage of random strings in the config. This can be useful for security purposes like HMAC keys.

### Comments

Single-line comments are started using `//`, `#` or `;`

Multiline Comments are wrapped in standard C multline comment slashes. `/* comment */`

### Numbers

Numbers are divided into integers and floats with no predefined precision.

Integers must be noted as digits in full, leading digits may be used. `0x`, `0o` and `0b` are used for hexadecimal, octal and binary notation respectively.

Floats can be noted in normal decimal notation, exponent notation or scientific notation:

```
0.001
1e-3
1.0*10^-3
```

### Keys and Lists and Maps

A string with an appended `:` is a map-key. These can be virtually used everywhere except after another map key.

The item after a map key is added to the current Map-List as a Map item. If an item has no key, it is added to the current Map-List as a List item.

The top level of a SECL file is considered a Map-List. A new Map-List is started by a parenthesis and must be terminated by one:

```
outmost-map

(
    new-map-here
)
```

The keywords `empty` and `nothing` are equivalent to `()` for readability

Internally, there is no difference between an list (also known as array) and a map, both use the same datatype and list elements can be mixed into dictionary entries.

Lists **must** be sorted exactly in the order that is observed in the input file. The list resulting from `(a b c)` **must** have `a` as first element and `c` as last. 

### Booleans

There are three categories of keywords for booleans:

* true
* false
* maybe

The true category contains `true`, `yes`, `on` and `allow`

The false category contains `false`, `no`, `off` and `deny`

The maybe category contains `maybe`. This keywords randomly evaluates to either true or false with a 50.1% chance of being true.

## Functions

Functions provide programmatic functionality to SECL, some functions are built-in but you can add your own, custom functions via exec.RegisterFunction.

Currently implemented functions:

* `nop` - does not parse any parameter, returns Nil
* `env` - Accepts 1 string parameter and 1 named parameter, it reads the environment variable specified in the parameter and if empty, uses the named parameter `default:` instead, see `tests/test06-env.exec.secl` for an example
* `loadb` - Load binary data from a file
* `loadf` - Load a single file as SECL and return the result maplist
* `loadv` - Load single value (anything but a maplist) from a file
* `decb64` - Decode base64 data
* `merge` - Merge several maplists into one maplist, it does not accept any parameters except maplists. Maplists are merged in order, it is not allowed to change a datatype of an entry.

Planned functions:

* `loadd` - Load all files with a specific extension from a folder