package historic

import "github.com/rodkranz/monitor/pkg/message"

// Writer is the function create log of message
type Writer func(msg *message.Message) error

// DefaultWriterLog function that does nothing
func DefaultWriterLog(_ *message.Message) error {
	return nil
}
