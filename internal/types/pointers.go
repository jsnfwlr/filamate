// Package types can be imported transparently to provide the different types and functions
// import this package transparently using `. "github.com/jsnfwlr/filamate/internal/types` and then use it as
// `strPointer := PointerOf("string")`
package types

// PointerOf returns $value as a pointer of it's type
func PointerOf[T any](value T) *T {
	return &value
}
