package historic

import "github.com/rodkranz/monitor/pkg/message"

type Writer func(msg *message.Message) error

func DefaultWriterLog(_ *message.Message) error {
	return nil
}
