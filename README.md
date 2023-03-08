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

## Usage

Note that full examples are provided in godoc.

### Getting a package:
```
go get -u github.com/p-alexander/configurator
```

### Using a simple variadic constructor:
```
// config is a private config example.
type config struct {
	i int
}

// Foo is an exported structure with a simple config.
type Foo struct {
	config *config
}

// WithI is used to configure Foo through NewFoo constructor.
func WithI(i int) configurator.Option[*config] {
	return func(c *config) error {
		c.i = i

		return nil
	}
}

// NewFoo is a variadic constructor for Foo.
func NewFoo(opts ...configurator.Option[*config]) (*Foo, error) {
	c := new(config)

	// call Constructor to execute all options on given config.
	if err := configurator.Constructor(c, opts); err != nil {
		return nil, fmt.Errorf("configurator.Constructor: %w", err)
	}

	return &Foo{
		config: c,
	}, nil
}
```

### Using a thread-safe configuration storage:
```
// config is a private config example.
type config struct {
	a string
}

// Bar is an exported structure with a thread-safe config storage.
type Bar struct {
	*configurator.Storage[*config]
}

// SetA is used to configure Bar.
func SetA(a string) configurator.Option[*config] {
	return func(c *config) error {
		c.a = a

		return nil
	}
}

// GetA returns a string value from a config.
func GetA() configurator.Getter[*config, string] {
	return func(c *config) (value string, err error) {
		return c.a, nil
	}
}

// NewBar is a constructor for Bar with underlying configuration storage.
func NewBar() *Bar {
	return &Bar{
		Storage: configurator.NewStorage(new(config)),
	}
}

func somewhereInYourCode() {
	// create Bar with thread-safe configuration.
	bar := NewBar()

	// a setter is safe to use concurrently.
	if err := configurator.ToStorage(bar.Storage, SetA("a")); err != nil {
		panic(err)
	}

	// same with a getter.
	if a, err := configurator.FromStorage(bar.Storage, GetA()); err != nil || a != "a" {
		panic(fmt.Sprintf("unexpected result: %v, %v", a, err))
	}
}
```
