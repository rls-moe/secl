# Breaking Changes

## v0.9.5

* Version `v0.9.5` will break serialization as the string type will now escape it's literal output on it's own. If you rely on Literal() to read a string, keep this in mind. This version also changes some semantics inside the parser which might lead to unexpected results in various applications but it should largely remain compatible with existing code. Additional this version removes the Position Information from String and Bool types as those values wheren't useful for development and complicated tests and code unnecessarily.

* Version `v0.9.5` breaks UnmarshalSECL interfaces; a bug in the unmarshaller caused it to call the function with the parent node instead of the correct subnode in certain circumstances. This has been fixed.

## v0.9.6

* Version `v0.9.6` breaks all internal interfaces and some internal behaviour. Functions gain an additional parameter, a runtime context, which allows for dynamic runtime environments (including runtime-defined functions). SECL files should still be parsed as normal as the internal changes are not exposed. Internal packages may no longer have access to function definitions in tests if the exec package was not loaded (this affects the lexer and phase1 tests)

* Maps will now be parsed in a rudamentary fashion, keys may only be strings and the MapList being unpacked must only have elements of one type, map[string]interface{} is not supported. Arrays remain unsupported