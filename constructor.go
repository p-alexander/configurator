package configurator

import "fmt"

// Constructor applies given set of options to a config.
func Constructor[T any](config T, opts []Option[T]) error {
	for i, opt := range opts {
		if err := opt(config); err != nil {
			return fmt.Errorf("option setter #%d (%T) failed: %w", i+1, opt, err)
		}
	}

	return nil
}
