package configurator_test

import (
	"errors"
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

	if err := configurator.Constructor[*config](c, []configurator.Option[*config]{opt("a")}); err != nil {
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

	c := new(config)

	err := configurator.Constructor[*config](c, []configurator.Option[*config]{opt("a")})
	if !errors.Is(err, customErr) {
		t.Fatalf("%v != %v", err, customErr)
	}
}
