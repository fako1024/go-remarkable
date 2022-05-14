package device

import (
	"io"
)

// Remarkable denotes a generic Remarkable device
type Remarkable interface {

	// Close closes the device
	Close() error

	// Frame retrieves a single frame
	Frame() ([]byte, error)

	// NewStream adds a new stream recipient on the provided writer
	NewStream(w io.Writer) error

	// Dimensions returns the width + height of the underlying frame(buffer)
	Dimensions() (int, int)

	// Upload uploads a file (PDF / ePUB) to the device tree
	Upload(docs ...Document) error
}
