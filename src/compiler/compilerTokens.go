package compiler

// all tokens are regex patterns

const compilerStartToken = "<script compiled>"
const compilerEndToken = "</script>"

const statementEndToken = ";"

// notice that there is a trailing space. This means that all variables must be
// defined in the format $ varName = "varValue";
const variableToken = "$"
const variableAssignmentToken = "="

const reactiveStartToken = "["
const reactiveEndToken = "]"

const eventToken = "@"
