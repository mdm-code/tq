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

<h4 align="center">Query TOML configuration files with the `tq` terminal utility</h4>

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
    <a href="https://goreportcard.com/report/github.com/mdm-code/tq">
        <img alt="Go report card" src="https://goreportcard.com/badge/github.com/mdm-code/tq">
    </a>
    <a href="https://pkg.go.dev/github.com/mdm-code/tq">
        <img alt="Go package docs" src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white">
    </a>
</p>
</div>

The `tq` program lets you query TOML configuration files with a sequence of
intuitive filters. It works as a regular Unix filter program reading input data
from the standard input and producing results to the standard output. Consult the
[package documentation](https://pkg.go.dev/github.com/mdm-code/tq) or check the
[Usage](#usage) section to see how you can use `tq`.


## Installation

Install the program to use the command-line `tq` to filter TOML files on the
terminal.

```sh
go install github.com/mdm-code/tq/cmd/tq@latest
```

Here is how you can get the whole Go package downloaded to fiddle with, but it
does not expose any public interfaces in code per se.

```sh
go get github.com/mdm-code/tq
```




## License

Copyright (c) 2024 Michał Adamczyk.

This project is licensed under the [MIT license](https://opensource.org/licenses/MIT).
See [LICENSE](LICENSE) for more details.
