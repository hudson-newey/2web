package compiler

// all tokens are regex patterns

const compilerStartToken = "<script compiled>"
const compilerEndToken = "</script>"

const statementEndToken = ";"

const variableToken = "$"
const variableAssignmentToken = "="

const reactiveStartToken = "["
const reactiveEndToken = "]"

const mustacheStartToken = "{{ "
const mustacheEndToken = " }}"

const eventStartToken = "&"
const eventEndToken = "*"
