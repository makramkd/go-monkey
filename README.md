# go-monkey

[Monkey](https://monkeylang.org/) is the programming language outlined in Thorsten Ball's [Writing an Interpreter in Go](https://interpreterbook.com/).

It's a small language that is meant to serve as an educational vehicle for teaching the basics of compilation and interpretation of high level code.

This is a Go implementation that closely follows the implementation in the book but has a few extensions:

* A minimal standard library and module system,
* For-each loops for both arrays and dictionaries,
* Possibly more, depending on what I come up with :)

## Building, Running the REPL and Tests

To build the top level interpreter, run the following after checking out the project:

```bash
go build ./cmd/monkeyc/
# or: go build .\cmd\monkeyc\ on Windows machines.
```

Run unit tests:

```bash
go test -cover ./...
```

Run the REPL:

```bash
./monkeyc
# or: .\monkeyc.exe on Windows
```
