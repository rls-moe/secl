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
        "scientific notation": 9*10^22
        "multi-line strings": "
            HELLO WORLD!"
        "omit quotation around simple strings": Hello-Again!
        "load other files from subfolders": !(loadd "folder" ".conf")
        "empty maps": empty
        
        "random values": (
            randstr128
            maybe
        )
        
        "mix lists and maps, they are the same in SECL"
    )
)
```

## Introduction

### Strings

Strings that do not contain any whitespace (as recognized by UTF8) and do not start with a digit or reserved character can omit quotation marks;

```
HelloWorld

```

Reserved characters: `"!@:()`

If any of the reserved characters is present, the string must be wrapped with quotation marks.

```
"Hello World"
```

Strings are multi-line by default, if you want to trim leading whitespace from each line and the beginning of the string, prepend a `@` symbol.

```bash
@"
Hello 
 World"
 
// equivalent to

"Hello
World"
```

The special values `randstr` and `randstr32` to `randstr256` allow the usage of random strings in the config. This can be useful for security purposes.

### Comments

Single-line comments are started using `//`, `#` or `;`

Multiline Comments are wrapped in standard C multline comment slashes. `/* comment */`

### Numbers

Numbers are divided into integers and floats with no predefined precision.

Integers must be noted as digits in full.

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

```bash
outmost-map

(
    new-map-here
)
```

The keywords `empty` and `nothing` are equivalent to `()` for readability

### Booleans

There are three categories of keywords for booleans:

* true
* false
* maybe

The true category contains `true`, `yes`, `on` and `allow`

The false category contains `false`, `no`, `off` and `deny`

The maybe category contains `maybe`. This keywords randomly evaluates to either true or false with a 50.1% chance of being true.
