// Package misc contains things which genuinely belong nowhere else and are too small to warrant a whole package!
package misc

// New returns the address of the passed parameter - similar to the behaviour of new in Go 1.26.
// NOTE: I was unable to get gopls to play nicely with the Go 1.26 feature and as of the time of writing it has only be released for a month.
func New[T any](value T) *T {
	return &value
}
