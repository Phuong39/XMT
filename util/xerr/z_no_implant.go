//go:build !implant

package xerr

// ExtendedInfo is a compile time constant to help signal if complex string
// values should be concatinated inline.
//
// This helps prevent debugging when the "-tags implant" option is enabled.
const ExtendedInfo = true

type strErr string
type wrapErr struct {
	e error
	s string
}

// New creates a new string backed error interface and returns it.
// This error struct does not support Unwrapping.
//
// The resulting errors created will be comparable.
func New(s string) error {
	return strErr(s)
}
func (e strErr) Error() string {
	return string(e)
}
func (e wrapErr) Error() string {
	return e.s
}
func (e wrapErr) Unwrap() error {
	return e.e
}

// Sub creates a new string backed error interface and returns it.
// This error struct does not support Unwrapping.
//
// If the "-tags implant" option is selected, the second value, the error code,
// will be used instead, otherwise it's ignored.
//
// The resulting errors created will be comparable.
func Sub(s string, _ uint8) error {
	return strErr(s)
}

// Wrap creates a new error that wraps the specified error.
//
// If not nil, this function will append ": " + 'Error()' to the resulting
// string message and will keep the original error for unwrapping.
//
// If "-tags implant" is specified, this will instead return the wrapped error
// directly.
func Wrap(s string, e error) error {
	if e != nil {
		return &wrapErr{s: s + ": " + e.Error(), e: e}
	}
	return &wrapErr{s: s}
}
