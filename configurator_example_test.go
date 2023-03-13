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

// NewBar is a variadic constructor for Bar with underlying configuration storage.
func NewBar(opts ...configurator.Option[*config]) (*Bar, error) {
	storage := configurator.NewStorage(new(config))

	if err := configurator.ToStorage(storage, opts...); err != nil {
		return nil, fmt.Errorf("configurator.ToStorage: %w", err)
	}

	return &Bar{
		Storage: storage,
	}, nil
}

// Example demonstrates building extendable constructors with configurator package.
func Example() {
	// example construction of a Foo with all options specified.
	foo1, err := NewFoo(WithI(1), WithA("a"))
	if err != nil {
		panic(err)
	}

	fmt.Printf("foo1: %d, %s;\n", foo1.config.i, foo1.config.a)

	// example construction of a Foo with only one option specified.
	foo2, err := NewFoo(WithA("a"))
	if err != nil {
		panic(err)
	}

	fmt.Printf("foo2: %d, %s;\n", foo2.config.i, foo2.config.a)

	// example construction of a Foo without options.
	foo3, err := NewFoo()
	if err != nil {
		panic(err)
	}

	fmt.Printf("foo3: %d, %s;\n", foo3.config.i, foo3.config.a)

	// example of a thread-safe configuration.
	bar1, err := NewBar(WithI(1))
	if err != nil {
		panic(err)
	}

	// a setter is safe to use concurrently.
	if err = configurator.ToStorage(bar1.Storage, WithA("a")); err != nil {
		panic(err)
	}

	// same with a getter.
	i, err := configurator.FromStorage(bar1.Storage, GetI())
	if err != nil {
		panic(err)
	}

	a, err := configurator.FromStorage(bar1.Storage, GetA())
	if err != nil {
		panic(err)
	}

	fmt.Printf("bar1: %d, %s;\n", i, a)
	// Output:
	// foo1: 1, a;
	// foo2: 0, a;
	// foo3: 0, ;
	// bar1: 1, a;
}
