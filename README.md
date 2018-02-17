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
            SELECT * FROM table
            WHERE a = \"b\"
            "
        "omit quotation around simple strings": Hello-Again
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

## Project Status

SECL is in Version "0.9.5".

Most features are now part of the reference implementation in this repository. The last missing function is `loadd`, merge has been implemented recently.

A non-formal specification can be found at [hackmd](https://hackmd.io/s/rylGmiXr-) or in this repository under SPECS.md (hackmd might be more up-to-date than SPECS.md)

## SECL vs JSON vs XML

SECL is primarly focused on configuration files, it may not be useful in other contexts but it can be.

As configuration files are typically not loaded repeatedly and performance is a secondary problem, SECL development does not focus on beating JSON and XML parsers in performance. Instead, they are beaten in feature sets.

Features that make SECL subjectively better than JSON or XML include;

* Simple multiline strings 
* Unified datatype for Arrays and Maps
* Numbers without fixed precision or limits
* Functions to include other SECL or even entire folders
* Functions to read environment variables
* Special keywords for cryptographically random strings
* UTF-8 Support in any string by default

## Introduction

To read up on how to write SECL, visit [INTRODUCTION.md](/INTRODUCTION.md). This file will contain detailed documentation and introduction into the SECL language.

## Contribution & License

SECL is licensed under the MPL 2.0, in essence, it's GPL2 but it will not "infect" any applications you write using SECL. By contributing you accept that your code is licensed under MPL2.0. A CLA is not necessary.

Until Version 1.0.0 SECL will not be considered stable, however, even in the current phase of the project, all breaking changes are noted in [BREAKING_CHANGES](/BREAKING_CHANGES.md). Note that we cannot possible know any and all ways a change may affect code. Certain edgecases may not be included at the discretion of the project leadership.

Once Version 1.0.0 is reached we will not break any dependent code using existing functions of the primary package (`go.rls.moe/secl`) and packages that are to be consumed by external packages. Internal breaking changes may be recorded but no stability is guaranteed.