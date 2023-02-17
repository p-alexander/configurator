# configurator [![GoDoc](https://img.shields.io/badge/godoc-Reference-brightgreen.svg?style=flat)](https://godoc.org/github.com/p-alexander/configurator)

Package `configurator` contains helpers for creating constructors with variadic options.

## About

A commonly used pattern for constructors in Go makes use of variadics to provide extendable configuration.

For example:
```
func NewFoo(opts ...Option) (*Foo, error) { ... }
```

Using this pattern a developer can provide only a few options first:
```
foo, err := NewFoo(WithI(1))
if err != nil {
    panic(err)
}
```

And then add additional options later without breaking a function signature:
```
foo, err := NewFoo(WithI(1), WithA("a"))
if err != nil {
    panic(err)
}
```

This package provides simple utilities to help build such constructors without too much code duplication.

Full example is provided in godoc.
