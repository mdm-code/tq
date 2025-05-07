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
    <a href="https://github.com/mdm-code/tq/actions?query=workflow%3ACI">
        <img alt="Build status" src="https://github.com/mdm-code/tq/workflows/CI/badge.svg">
    </a>
    <a href="https://app.codecov.io/gh/mdm-code/tq">
        <img alt="Code coverage" src="https://codecov.io/gh/mdm-code/tq/branch/main/graphs/badge.svg?branch=main">
    </a>
    <a href="https://opensource.org/licenses/MIT" rel="nofollow">
        <img alt="MIT license" src="https://img.shields.io/github/license/mdm-code/tq">
    </a>
    <a href="https://goreportcard.com/report/github.com/mdm-code/tq/v2">
        <img alt="Go report card" src="https://goreportcard.com/badge/github.com/mdm-code/tq/v2">
    </a>
    <a href="https://pkg.go.dev/github.com/mdm-code/tq/v2">
        <img alt="Go package docs" src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white">
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
with the command. Here is table with the supported filter expressions and some
examples to get you going on how to use `tq` in your workflow.

Some effort has been made to make queries less clunky to type out on the
command line and the syntax for queries more aligned with the TOML syntax and
semantics. It's been decided to drop the requirement for square brackets for
selectors and quotation marks for bare strings. Queries can now span across
multiple lines so that they are still legible as their complexity increases.
Longer queries run in a shell script might benefit for it. As for quoted
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

# Other relevant notations like date, time, date-time, with and without the
# offset, inf, nan, negative numbers, stay the way they're written in the
# input file.
```


### Multiline query with bare strings

Here is a dummy configuration file in TOML found on the web for Gitlab
connected to a Kubernetes. The file attempts to configure some Gitlab runners.
The file is (1) queried with the key `runners` to access the table that is then
(2) converted to an iterator with `[]`. Then (3) the query goes for
`kubernetes`, `volumes`, and `host_path` in this order to (4) turn the last one
to an iterator with `[]`, and then (5) query each element of the iterator for
`"host path"`. Mind the quoted string with the space.


```sh
tq -q '
    .runners[]
        .kubernetes
        .volumes
        .host_path[]
            ."host path"
' << EOF
[session_server]
  session_timeout = 1800

[[runners]]
  name = ""
  url = ""
  token = ""
  executor = "kubernetes"
  cache_dir = "/tmp/gitlab/cache"
  [runners.kubernetes]
    host = ""
    bearer_token_overwrite_allowed = false
    image = ""
    namespace = ""
    namespace_overwrite_allowed = ""
    privileged = false
    memory_limit = "1Gi"
    service_account_overwrite_allowed = ""
    pod_annotations_overwrite_allowed = ""
    [runners.kubernetes.node_selector]
      gitlab = "true"
    [runners.kubernetes.volumes]
      [[runners.kubernetes.volumes.host_path]]
        name = "gitlab-cache"
        mount_path = "/tmp/gitlab/cache"
        "host path" = "/home/core/data/gitlab-runner/data"

[[runners]]
  name = "runner-gitlab-runner-xxx-xxx"
  url = "https://gitlab.com/"
  token = "<my-token>"
  executor = "kubernetes"
  [runners.cache]
    [runners.cache.s3]
    [runners.cache.gcs]
  [runners.kubernetes]
    host = ""
    bearer_token_overwrite_allowed = false
    image = "ubuntu:16.04"
    namespace = "gitlab-managed-apps"
    namespace_overwrite_allowed = ""
    privileged = true
    service_account_overwrite_allowed = ""
    pod_annotations_overwrite_allowed = ""
    [runners.kubernetes.volumes]
EOF

Output:

/home/core/data/gitlab-runner/data
```


### Retrieve IPs from a table of server tables

In the example below, the TOML input file is (1) queried with the key
`servers`, then (2) the retrieved table is converted to an iterator of objects
with `[]`, and then (3) the IP address is recovered from each of the objects
with the quoted key `"ip"`.

```sh
tq -q '.servers[]."ip"' <<EOF
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
```


### Retrieve selected ports from a list of databases

This example uses the older syntax and queries the TOML input for the for the
all ports aside from the first one assigned to the first database record on the
list.

```sh
tq -q '.["databases"][0]["ports"][1:][]' <<EOF
databases = [ {enabled = true, ports = [ 5432, 5433, 5434 ]} ]
EOF

Output:

5433
5434
```


### Run inside of a container

If you don't feel like installing `tq` with `go install`, you can test `tq` out
running inside of a container with this command:

```sh
docker run -i ghcr.io/mdm-code/tq:latest tq -q ".dependencies.ignore" <<EOF
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


## Development

Go through the [Makefile](Makefile) to get an idea of the formatting, testing
and linting that can be used locally for development purposes.


## License

Copyright (c) 2025 Michał Adamczyk.

This project is licensed under the [MIT license](https://opensource.org/licenses/MIT).
See [LICENSE](LICENSE) for more details.
