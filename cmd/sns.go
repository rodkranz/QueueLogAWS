package cmd

import (
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/go-clog/clog"
	"github.com/urfave/cli"

	"github.com/rodkranz/monitor/pkg/setting"
	"github.com/rodkranz/monitor/pkg/tool"
)

// SNS is variable to export command to send message in AWS SNS
var SNS = cli.Command{
	Name:        "sns",
	Description: "Send message to SNS",
	Usage:       "Send message to queue",
	Before:      snsSetting,
	Action:      runSNS,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "instance", Usage: "Address of instance", Value: os.Getenv("INSTANCE")},
		cli.StringFlag{Name: "id", Usage: "Instance identification", Value: os.Getenv("ID")},
		cli.StringFlag{Name: "topic", Usage: "Topic to publish", Value: os.Getenv("TOPIC")},
		cli.StringFlag{Name: "data", Usage: "Data to send in topic"},
		cli.StringFlag{Name: "file", Usage: "File to send in topic"},
	},
}

var body []byte

func snsSetting(c *cli.Context) (err error) {
	if c.String("instance") == "" || c.String("id") == "" || c.String("topic") == "" {
		return ErrMissingParameters{"instance"}
	}

	setting.SNS.Instance = c.String("instance")
	setting.SNS.ID = c.String("id")
	setting.SNS.Topic = c.String("topic")

	if c.IsSet("data") {
		body = []byte(c.String("data"))
	}

	if c.IsSet("file") {
		f, err := os.Open(c.String("file"))
		if err != nil {
			return err
		}
		body, err = ioutil.ReadAll(f)
		if err != nil {
			return err
		}
	}

	if len(body) == 0 {
		return ErrBodyEmpty{}
	}

	return nil
}

func runSNS(_ *cli.Context) error {
	awsSession := session.Must(session.NewSession())
	snsMessage := &sns.PublishInput{
		Message:  aws.String(string(body)),
		TopicArn: aws.String(setting.GetQueue()),
	}

	clog.Info("Topic: [%s]", setting.GetQueue())
	clog.Info(" Body: [%s]", tool.MD5(body))

	_, err := sns.New(awsSession).Publish(snsMessage)
	if err != nil {
		clog.Error(0, err.Error())
	}

	return err
}
