/*
tq - query TOML configuration files

Usage:

	tq [-qtmsi] [file...]

Options:

	-h, --help              show this help message and exit
	-q, --query             query to run against the input data (default: '.')
	-t, --tables-inline     emit tables inline (default: false)
	-m, --arrays-multiline  emit arrays one element per line (default: false)
	-s, --indent-symbol     provide the indentation string (default: '  ')
	-i, --indent-tables     indent tables and array tables (default: false)

Example:

	<<EOF tq -q .servers[].ip
	[servers]

	[servers.prod]
	ip = "10.0.0.1"
	role = "backend"

	[servers.staging]
	ip = "10.0.0.2"
	role = "backend"
	EOF

Output:

	10.0.0.1
	10.0.0.2

Tq is a tool for querying TOML configuration files with a sequence of intuitive
filters. It works as a regular Unix filter program reading input data from the
standard input and producing results to the standard output.
*/
package main
