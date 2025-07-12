# 2Web Compiler

## 1-Preprocessing

During pre-processing, we perform preprocessing that does **not** require
lexical analysis.

E.g. Expanding document partials, concatenating files,modifying encoding
types, find & replace, etc...

## 2-Lexer

During the lexer stage, we convert the text content into nodes.

Note that during the lexer stage, we just perform local grouping.

## 3-Validation

During the validation stage, we validate that the structure returned by the
lexer is correct.

## 4-Parser

During the parsing stage, we construct an abstract syntax tree, and turn
combinations fo lexer tokens into usable structs.

## 5-Templating

During templating, we modify the original source code so that if behaves
correctly.

Note that because 2web templates are a superset of html, css, and JavaScript,
our templating is quite simple given that we don't have to do any language
translations or use an intermediate language.
