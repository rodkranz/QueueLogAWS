package historic

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/go-clog/clog"

	"github.com/rodkranz/monitor/pkg/message"
)

func getDate(layout string) string {
	return time.Now().Format(layout)
}

// FileConfig is the configuration of FileWriter depends
type FileConfig struct {
	Path string
}

// FileWriter configure and return the Writer function for log.
func FileWriter(cfg FileConfig) (Writer, error) {
	return func(msg *message.Message) error {
		// path + year + month + day + topic
		dir := path.Join(cfg.Path, getDate("2006"), getDate("01"), getDate("02"), msg.Topic())

		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}

		// hour + minutes + seconds + prefix + messageId + suffix
		name := fmt.Sprintf("%s-%s", getDate("15.04.05"), *msg.MessageId)

		f, err := os.Create(path.Join(dir, name))
		if err != nil {
			if os.IsExist(err) {
				clog.Warn("Message already registered")
			}
			return err
		}
		defer f.Close()

		fmt.Fprint(f, *msg.Body)
		return nil
	}, nil
}
