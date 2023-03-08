package configurator_test

import (
	"fmt"

	"github.com/p-alexander/configurator"
)

// config is a private config example.
type config struct {
	i int
	a string
}

// Foo is an exported structure with a simple config.
type Foo struct {
	config *config
}

// Bar is an exported structure with a thread-safe config storage.
type Bar struct {
	*configurator.Storage[*config]
}

// WithI is used to modify a config, integer setter.
func WithI(i int) configurator.Option[*config] {
	return func(c *config) error {
		c.i = i

		return nil
	}
}

// WithA is used to modify a config, string setter.
func WithA(a string) configurator.Option[*config] {
	return func(c *config) error {
		c.a = a

		return nil
	}
}

// GetI returns an integer value from a config.
func GetI() configurator.Getter[*config, int] {
	return func(c *config) (value int, err error) {
		return c.i, nil
	}
}

// GetA returns a string value from a config.
func GetA() configurator.Getter[*config, string] {
	return func(c *config) (value string, err error) {
		return c.a, nil
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

// NewBar is a constructor for Bar with underlying configuration storage.
func NewBar() *Bar {
	return &Bar{
		Storage: configurator.NewStorage(new(config)),
	}
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

	// example of a thread-safe configuration.
	bar := NewBar()

	// a setter is safe to use concurrently.
	if err = configurator.ToStorage(bar.Storage, WithI(1), WithA("a")); err != nil {
		panic(err)
	}

	// same with a getter.
	if i, err := configurator.FromStorage(bar.Storage, GetI()); err != nil || i != 1 {
		panic(fmt.Sprintf("unexpected result: %v, %v", i, err))
	}

	if a, err := configurator.FromStorage(bar.Storage, GetA()); err != nil || a != "a" {
		panic(fmt.Sprintf("unexpected result: %v, %v", a, err))
	}
}
