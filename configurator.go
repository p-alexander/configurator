// Package configurator contains helpers for creating constructors with variadic options. See usage in the example.
package configurator

import "fmt"

// Option is a generic option for variadic constructors.
type Option[T any] func(t T) error

// Constructor applies given set of options to a config.
func Constructor[T any](config T, opts []Option[T]) error {
	for i, opt := range opts {
		if err := opt(config); err != nil {
			return fmt.Errorf("option #%d failed: %w", i+1, err)
		}
	}

	return nil
}
