/*
Package parser provides a Parser struct that parses lexer tokens from the tq
query into an AST.

GRAMMAR
-------
root        = query
query       = *(filter)
filter      = identity / selector
identity    = DOT
selector    = '[' *( STRING / INTEGER / span ) ']'
span        = *INTEGER ':' *INTEGER
*/
package parser
