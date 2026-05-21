package utils

import (
	"reflect"
)

func UpdateStruct(existing interface{}, updates interface{}) {
	existingVal := reflect.ValueOf(existing).Elem()
	updatesVal := reflect.ValueOf(updates)

	for i := range updatesVal.NumField() {
		field := updatesVal.Field(i)

		// Skip zero values (default values)
		if !field.IsZero() {
			existingField := existingVal.Field(i)

			// Ensure the field can be set (ignores unexported fields)
			if existingField.CanSet() {
				existingField.Set(field)
			}
		}
	}
}

func FilterSlice[T any](items []T, filterFunc func(T) bool) []T {
	result := make([]T, 0, len(items)) // Preallocate memory

	for _, item := range items {
		if filterFunc(item) {
			result = append(result, item)
		}
	}
	return result
}

func FindFirst[T any](items []T, predicate func(T) bool) (T, bool) {
	var zeroValue T
	
	for _, item := range items {
		if predicate(item) {
			return item, true
		}
	}
	return zeroValue, false
}

func MapSlice[T any, R any](source []T, transformer func(T) R) []R {
	if source == nil {
		return nil
	}
	if len(source) == 0 {
		return []R{}
	}

	result := make([]R, len(source))

	for i, item := range source {
		result[i] = transformer(item)
	}

	return result
}