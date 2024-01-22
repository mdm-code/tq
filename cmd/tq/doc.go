/*
tq - query TOML configuration files

Usage:

	tq [-q|--query arg...] [file...]

Options:

	-h, --help         show this help message and exit
	-q, --query        specify the query to run against the input data (default: '.')
	--tablesInline     emit all tables inline (default: false)
	--arraysMultiline  emit all arrays with one element per line (default: false)
	--indentSymbol     provide the string for the indentation level (default: '  ')
	--indentTables     indent tables and array tables literals (default: false)

Example:

	tq -q '["servers"][]["ip"]' <<EOF
	[servers]

	[servers.prod]
	ip = "10.0.0.1"
	role = "backend"

	[servers.staging]
	ip = "10.0.0.2"
	role = "backend"
	EOF

Output:

	'10.0.0.1'
	'10.0.0.2'

Tq is a tool for querying TOML configuration files with a sequence of intuitive
filters. It works as a regular Unix filter program reading input data from the
standard input and producing results to the standard output.
*/
package main
