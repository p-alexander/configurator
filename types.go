package configurator

type (
	// Option is a generic option for variadic constructors and Storage.
	Option[T any] func(t T) error

	// Getter is a generic getter to retrieve typed values from Storage.
	Getter[T, O any] func(config T) (value O, err error)
)
