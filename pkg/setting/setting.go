package setting

import (
	"fmt"
	olog "log"
	"os"

	"github.com/go-clog/clog"
	"github.com/urfave/cli"

	"github.com/rodkranz/monitor/pkg/historic"
)

var (
	// AWS is the basic information for authentication in aws
	AWS struct {
		AccessKeyID     string
		SecretAccessKey string
		Region          string
	}

	// SNS information about SNS that possibility send message
	SNS struct {
		Instance string
		ID       string
		Topic    string
	}

	// SQS is the information about SQS to listen messages
	SQS struct {
		QueueURL    string
		NumMessages int64
		WaitTime    int64
	}

	// Log is the definition of
	Log struct {
		Driver string
	}

	// WLog function default to write log.
	WLog historic.Writer = historic.DefaultWriterLog
)

func init() {
	err := clog.New(clog.CONSOLE, clog.ConsoleConfig{})
	if err != nil {
		olog.Fatalf("error to create log: %s", err)
		os.Exit(1)
	}
}

// GetQueue return the topic formated
func GetQueue() string {
	return fmt.Sprintf("%s:%s:%s", SNS.Instance, SNS.ID, SNS.Topic)
}

// Setting define default information of system
func Setting(c *cli.Context) error {
	if c.String("aws_access_key") == "" || c.String("aws_access_key") == "" || c.String("aws_access_key") == "" {
		return fmt.Errorf("missing parameters")
	}

	// AWS
	AWS.AccessKeyID = c.String("aws_access_key")
	AWS.SecretAccessKey = c.String("secret_access_key")
	AWS.Region = c.String("region")

	// Log
	if err := serviceLog(c); err != nil {
		return err
	}

	// Log
	return nil
}

func serviceLog(c *cli.Context) (err error) {
	switch c.String("log_driver") {
	case "sql":
		if c.String("log_param") == "" {
			return fmt.Errorf("missing parameter log_param")
		}

		WLog, err = historic.SQLWriter(historic.SQLConfig{Path: c.String("log_param")})
	case "file":
		if c.String("log_param") == "" {
			return fmt.Errorf("missing parameter log_param")
		}

		// Writer log
		WLog, err = historic.FileWriter(historic.FileConfig{Path: c.String("log_param")})
	default:
		err = fmt.Errorf("log must be defined the optionas are [file] or [sql]")
	}

	return err
}
