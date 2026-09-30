<h1 align="center">
  <div >
    <img
      src="https://raw.githubusercontent.com/mdm-code/mdm-code.github.io/main/tq_logo.png"
      alt="logo"
      style="object-fit: contain"
      width="30%"
    />
  </div>
</h1>

<h4 align="center">Query TOML configuration files with the Tq terminal utility</h4>

<div align="center">
<p>
    <a href="https://pkg.go.dev/github.com/mdm-code/tq/v2">
        <img alt="Go package docs" src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white">
    </a>
    <a href="https://opensource.org/licenses/MIT" rel="nofollow">
        <img alt="MIT license" src="https://img.shields.io/github/license/mdm-code/tq">
    </a>
    <a href="https://github.com/mdm-code/tq/actions?query=workflow%3ATest">
        <img alt="Tests" src="https://github.com/mdm-code/tq/workflows/Test/badge.svg">
    </a>
    <a href="https://github.com/mdm-code/tq/actions?query=workflow%3ALint">
        <img alt="Linters" src="https://github.com/mdm-code/tq/workflows/Lint/badge.svg">
    </a>
    <a href="https://app.codecov.io/gh/mdm-code/tq">
        <img alt="Code coverage" src="https://codecov.io/gh/mdm-code/tq/branch/main/graphs/badge.svg?branch=main">
    </a>
</p>
</div>

The `tq` program lets you query TOML configuration files with a sequence of
intuitive filters. It works as a regular Unix filter program reading input data
from the standard input and producing results to the standard output. Consult the
[package documentation](https://pkg.go.dev/github.com/mdm-code/tq/v2) or check the
[Usage](#usage) section to see how you can use `tq`.


## Installation

Install the program and use `tq` on the command-line to filter TOML files on
the terminal.

```sh
go install github.com/mdm-code/tq/v2/cmd/tq@latest
```

Here is how you can get the whole Go package downloaded to fiddle with, but
it exposes only the public interfaces for `tq` and the TOML adapter so that
the latter can be swapped out.

```sh
go get github.com/mdm-code/tq/v2
```


## Usage

Enter `tq -h` to get usage information and the list of options that can be used
with the command. Here is a table with the supported filter expressions and some
examples to get you going on how to use `tq` in your workflow.

Some effort has been made to make queries less clunky to type out on the
command line and the syntax for queries more aligned with the TOML syntax and
semantics. It's been decided to drop the requirement for square brackets for
selectors and quotation marks for bare strings. Queries can now span across
multiple lines so that they are still legible as their complexity increases.
Longer queries run in a shell script might benefit from it. As for quoted
strings, both inverted commas and quotes can be used. A note of caution though
that these should be used such that they do not interfere with shell quoting.


### Supported filters

| <a href="#supported-filters"><img width="1000" height="0"></a><p>Filter</p> | <a href="#supported-filters"><img width="1000" height="0"></a><p>Expression</p>                     |
| :-------------------------------------------------------------------------: | :-------------------------------------------------------------------------------------------------: |
| <kbd><b>identity</b></kbd>                                                  | <kbd><b>.</b></kbd>                                                                                 |
| <kbd><b>key</b></kbd>                                                       | <kbd><b>["string"]</b></kbd> or <kbd><b>"quoted string"</b></kbd> or <kbd><b>bare-string</b></kbd>  |
| <kbd><b>index</b></kbd>                                                     | <kbd><b>[0]</b></kbd>                                                                               |
| <kbd><b>iterator</b></kbd>                                                  | <kbd><b>[]</b></kbd>                                                                                |
| <kbd><b>span</b></kbd>                                                      | <kbd><b>[:]</b></kbd>                                                                               |

> [!NOTE]
> `?` is not a standalone filter. It is a modifier for the preceding filter
> expression (for example `[]?` or `.key?`) and makes that filter optional
> by skipping incompatible inputs instead of failing the whole query.


### Supported escape sequences for quoted strings

Commonly found characters are mapped onto often used escaped sequences. These
can be used in quoted strings mostly the same way one would use them in a TOML
file though the specification for the TOML language advises against the use of
funky keys unless there is a good reason to use them.

```txt
\b          - backspace
\t          - tab
\n          - linefeed
\f          - form feed
\r          - carriage return
\"          - double quote
\'          - single quote
\\          - backslash
\uhhhh      - short 16-bit hexadecimal form
\Uhhhhhhhh  - long 32-bit hexadecimal form
```


### Conversion caveats

Given the current implementation, *most* values are represented exactly the way
they are spelled out in the input file after they're queried for. The following
value notations will be converted to a different notation related to the
backing Go type:

```txt
1_000       => 1000   # Underscores are not retained.
0xFFFF      => 65535  # Hexadecimal is converted to decimal.
0o755       => 493    # Octal is converted to decimal.
0b1111_1111 => 255    # Binary is converted to decimal.
+100        => 100    # The plus sign is dropped.
5e-3        => 0.005  # The exponential notation is not kept.
```

> [!NOTE]
> Other relevant notations like date, time, date-time, with and without the
> offset, inf, nan, negative numbers, stay the way they're written in the
> input file.


### Simple query with optional modifier

Here's a sample TOML configuration file with a query handling missing keys
gracefully. The query works as follows: (1) iterates over `runners` in the
array of tables; (2) navigates to each `policy` inside `kubernetes.autoscaler`;
(3) selects the second element from the `periods` array with `?` modifier
indicating that the filter is optional so that it won't fail if the index
doesn't exist.

```sh
<< EOF tq -q '
.runners[]
  .kubernetes.autoscaler.policy[]
    .periods[1]?
'
[[runners]]
  name = "kubernetes-runner"
  executor = "kubernetes"
  [runners.kubernetes]
    namespace = "gitlab-runner"
    cpu_request = "500m"
    memory_request = "1Gi"
    [runners.kubernetes.autoscaler]
      max_pause_pods = 10
      [[runners.kubernetes.autoscaler.policy]]
        idle_count = 5
        periods = ["* * * * *", "* 8-17 * * mon-fri"]
        timezone = "UTC"
      [[runners.kubernetes.autoscaler.policy]]
        idle_count = 0
        periods = ["* * * * *"]
EOF
```

```txt
Output:

* 8-17 * * mon-fri
```


### Run inside a container

If you don't feel like installing `tq` with `go install`, you can test `tq` out
running inside a container with this command:

```sh
<<EOF docker run -i ghcr.io/mdm-code/tq:latest tq -q ".dependencies.ignore"
[dependencies]
anyhow = "1.0.75"
bstr = "1.7.0"
grep = { version = "0.3.1", path = "crates/grep" }
ignore = { version = "0.4.22", path = "crates/ignore" }
lexopt = "0.3.0"
log = "0.4.5"
serde_json = "1.0.23"
termcolor = "1.1.0"
textwrap = { version = "0.16.0", default-features = false }
EOF
```

```txt
Output:

path = 'crates/ignore'
version = '0.4.22'
```


## Development

Development and local builds require Go 1.26 or newer.

Go through the [Makefile](Makefile) to get an idea of the formatting, testing
and linting that can be used locally for development purposes.

This project uses `golangci-lint`. [Here](https://golangci-lint.run/docs/welcome/install/local/)
is how to install it on your local machine.


## License

Copyright (c) 2026 Michał Adamczyk.

This project is licensed under the [MIT license](https://opensource.org/licenses/MIT).
See [LICENSE](LICENSE) for more details.
