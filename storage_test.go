package configurator_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/p-alexander/configurator"
)

func TestToStorage(t *testing.T) {
	opt1 := func(i int) configurator.Option[*config] {
		return func(c *config) error {
			c.i = i

			return nil
		}
	}

	opt2 := func(a string) configurator.Option[*config] {
		return func(c *config) error {
			c.a = a

			return nil
		}
	}

	c := new(config)
	s := configurator.NewStorage(c)
	wg := new(sync.WaitGroup)
	wg.Add(10000)

	for i := 0; i < 10000; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			if err := configurator.ToStorage(s, opt1(1), opt2("a")); err != nil {
				t.Error(err)

				return
			}
		}(wg)
	}

	wg.Wait()

	if c.a != "a" || c.i != 1 {
		t.Fatal("unexpected:", c.a, c.i)
	}
}

func TestFromStorage(t *testing.T) {
	getter1 := func(c *config) (option any, err error) {
		return c.i, nil
	}

	getter2 := func(c *config) (option any, err error) {
		return c.a, nil
	}

	c := &config{1, "a"}
	s := configurator.NewStorage(c)
	wg := new(sync.WaitGroup)
	wg.Add(10000)

	for i := 0; i < 10000; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			if n, _ := configurator.FromStorage(s, getter1); n != 1 {
				t.Error("unexpected:", n)
				return
			}

			if a, _ := configurator.FromStorage(s, getter2); a != "a" {
				t.Error("unexpected:", a)
				return
			}
		}(wg)
	}

	wg.Wait()
}

func ExampleFromStorage() {
	// configuration.
	type exampleConfig struct {
		i int
	}

	getter := func(config *exampleConfig) (value int, err error) {
		return config.i, nil
	}

	s := configurator.NewStorage(&exampleConfig{i: 1})
	i, _ := configurator.FromStorage(s, getter)
	fmt.Println(i)
	// Output: 1
}

func ExampleToStorage() {
	// configuration.
	type exampleConfig struct {
		i int
	}

	// int option.
	setter := func(i int) configurator.Option[*exampleConfig] {
		return func(c *exampleConfig) error {
			c.i = i

			return nil
		}
	}

	s := configurator.NewStorage(&exampleConfig{i: 1})
	err := configurator.ToStorage(s, setter(1))
	fmt.Println(err)
	// Output: <nil>
}

func ExampleStorage() {
	// configuration.
	type exampleConfig struct {
		i int
	}

	// int option.
	setter := func(i int) configurator.Option[*exampleConfig] {
		return func(c *exampleConfig) error {
			c.i = i

			return nil
		}
	}

	// getter to be executed behind a lock.
	getter := func(config *exampleConfig) (value int, err error) {
		return config.i, nil
	}

	s := configurator.NewStorage(new(exampleConfig))
	wg := new(sync.WaitGroup)
	wg.Add(10000)

	for i := 0; i < 10000; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			if err := configurator.ToStorage(s, setter(1)); err != nil {
				panic(err)
			}

			if n, _ := configurator.FromStorage(s, getter); n != 1 {
				panic(n)
			}
		}(wg)
	}

	wg.Wait()

	n, _ := configurator.FromStorage(s, getter)
	fmt.Println(n)
	// Output: 1
}
