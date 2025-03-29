package compiler

const compilerStartToken = "<script compiled>"
const compilerEndToken = "</script>"

// notice that there is a trailing space. This means that all variables must be
// defined in the format $ varName = "varValue"
const variableToken = "$ "
const variableAssignmentToken = "="

const eventToken = "@"
