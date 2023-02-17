package configurator_test

import (
	"fmt"

	"github.com/p-alexander/configurator"
)

// config is private config example.
type config struct {
	i int
	a string
}

// Foo is an exported structure.
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

// WithA is used to configure Foo through NewFoo constructor.
func WithA(a string) configurator.Option[*config] {
	return func(c *config) error {
		c.a = a

		return nil
	}
}

// NewFoo is a constructor for Foo
func NewFoo(opts ...configurator.Option[*config]) (*Foo, error) {
	c := new(config)

	// call Constructor to execute all options on given config.
	if err := configurator.Constructor[*config](c, opts); err != nil {
		return nil, fmt.Errorf("ConfigConstructor: %w", err)
	}

	return &Foo{
		config: c,
	}, nil
}

// Example demonstrates building extendable constructors with configurator package.
func Example() {
	// example construction of a Foo with all options specified.
	foo1, err := NewFoo(WithI(1), WithA("a"))
	if err != nil {
		panic(err)
	}

	if foo1.config == nil || foo1.config.i != 1 || foo1.config.a != "a" {
		panic(fmt.Sprintf("unexpected result: %+v", foo1.config))
	}

	// example construction of a Foo with only one option specified.
	foo2, err := NewFoo(WithA("a"))
	if err != nil {
		panic(err)
	}

	if foo2.config == nil || foo2.config.i != 0 || foo2.config.a != "a" {
		panic(fmt.Sprintf("unexpected result: %+v", foo2.config))
	}

	// example construction of a Foo without options.
	foo3, err := NewFoo()
	if err != nil {
		panic(err)
	}

	if foo3.config == nil || foo3.config.i != 0 || foo3.config.a != "" {
		panic(fmt.Sprintf("unexpected result: %+v", foo3.config))
	}
}
