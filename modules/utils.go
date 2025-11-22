package modules

import "fmt"

func GetArg[T any](mapAny map[string]any, key string) (T, error) {
	var zero T

	if mapAny == nil {
		return zero, fmt.Errorf("map is nil")
	}

	value, exists := mapAny[key]
	if !exists {
		return zero, fmt.Errorf("key '%s' not found", key)
	}

	typedValue, ok := value.(T)
	if !ok {
		return zero, fmt.Errorf("invalid type for key '%s': expected %T, got %T",
			key, zero, value)
	}

	return typedValue, nil
}
