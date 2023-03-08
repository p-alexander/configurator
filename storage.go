package configurator

import (
	"fmt"
	"sync"
)

// Getter is a typed getter function for Storage.
type Getter[T, O any] func(config T) (value O, err error)

// Storage is a thread-safe config holder.
type Storage[T any] struct {
	mtx sync.RWMutex

	config T
}

// NewStorage returns a Storage for given config.
func NewStorage[T any](config T) *Storage[T] {
	return &Storage[T]{
		config: config,
	}
}

// ToStorage sets a typed value to given config storage according to specified Getter.
func ToStorage[T any](storage *Storage[T], setters ...Option[T]) error {
	if err := storage.set(setters...); err != nil {
		return fmt.Errorf("storage.set: %w", err)
	}

	return nil
}

// FromStorage returns a typed value from given config storage according to specified Getter.
func FromStorage[T, O any](storage *Storage[T], getter Getter[T, O]) (value O, err error) {
	v, err := storage.get(
		func(config T) (value any, err error) {
			return getter(config)
		},
	)
	if err != nil {
		return value, fmt.Errorf("storage.get: %w", err)
	}

	return v.(O), nil
}

// get returns a value specified by given getter from a stored config in a thread-safe way.
func (m *Storage[T]) get(getter func(T) (value any, err error)) (any, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	opt, err := getter(m.config)
	if err != nil {
		return nil, fmt.Errorf("option getter %T failed: %w", getter, err)
	}

	return opt, nil
}

// set applies given options to a stored config.
func (m *Storage[T]) set(opts ...Option[T]) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for i, opt := range opts {
		if err := opt(m.config); err != nil {
			return fmt.Errorf("option #%d (%T) failed: %w", i+1, opt, err)
		}
	}

	return nil
}
