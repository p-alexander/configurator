package configurator_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/p-alexander/configurator"
)

func TestConstructorHappyCase(t *testing.T) {
	opt := func(a string) configurator.Option[*config] {
		return func(c *config) error {
			c.a = a

			return nil
		}
	}

	c := new(config)

	if err := configurator.Constructor(c, []configurator.Option[*config]{opt("a")}); err != nil {
		t.Fatal(err)
	}

	if c.a != "a" {
		t.Fatalf("%s != a", c.a)
	}
}

func TestConstructorError(t *testing.T) {
	customErr := errors.New("test")

	opt := func(a string) configurator.Option[*config] {
		return func(c *config) error {
			return customErr
		}
	}

	err := configurator.Constructor(new(config), []configurator.Option[*config]{opt("a")})
	if !errors.Is(err, customErr) {
		t.Fatalf("%v != %v", err, customErr)
	}
}

func ExampleConstructor() {
	// configuration.
	type exampleConfig struct {
		a string
		i int
	}

	// structure to be configured.
	type exampleStruct struct {
		config *exampleConfig
	}

	// int option.
	optI := func(i int) configurator.Option[*exampleConfig] {
		return func(c *exampleConfig) error {
			c.i = i

			return nil
		}
	}

	// string option.
	optA := func(a string) configurator.Option[*exampleConfig] {
		return func(c *exampleConfig) error {
			c.a = a

			return nil
		}
	}

	// variadic constructor for exampleStruct.
	exampleConstructor := func(opts ...configurator.Option[*exampleConfig]) (*exampleStruct, error) {
		c := new(exampleConfig)

		if err := configurator.Constructor(c, opts); err != nil {
			return nil, fmt.Errorf("configurator.Constructor: %w", err)
		}

		return &exampleStruct{
			config: c,
		}, nil
	}

	// example usage.
	es, err := exampleConstructor(optI(1), optA("a"))
	if err != nil {
		panic(err)
	}

	fmt.Println(es.config.i, es.config.a)
	// Output: 1 a
}
