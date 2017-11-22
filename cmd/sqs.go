package cmd

import (
	"fmt"
	"os"
	
	"github.com/urfave/cli"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-clog/clog"
	
	"github.com/rodkranz/monitor/pkg/setting"
	"github.com/rodkranz/monitor/pkg/message"
	_ "github.com/rodkranz/monitor/pkg/historic"
)

var SQS = cli.Command{
	Name:        "sqs",
	Description: "Listen SQS queue",
	Usage:       "Listen messages from queue",
	Action:      runSQS,
	Before:      SQSSetting,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "queue_url", Usage: "", Value: os.Getenv("QUEUE_URL")},
		cli.IntFlag{Name: "max_message", Usage: "", Value: 10},
		cli.IntFlag{Name: "wait_time", Usage: "", Value: 20},
	},
}

func SQSSetting(c *cli.Context) error {
	if c.String("queue_url") == "" {
		return fmt.Errorf("parameter %s is missing", "queue_url")
	}
	
	setting.SQS.QueueUrl = c.String("queue_url")
	setting.SQS.NumMessages = c.Int64("max_message")
	setting.SQS.WaitTime = c.Int64("wait_time")
	
	return nil
}

func runSQS(_ *cli.Context) error {
	awsSession := sqs.New(session.Must(session.NewSession()))
	msg := make(chan *message.Message)
	out := make(chan error)
	
	go func() {
		for {
			err := receiveMessage(awsSession, msg)
			if err != nil {
				out <- err
			}
		}
		close(msg)
	}()
	
	return readMessage(msg, out)
}

func receiveMessage(s *sqs.SQS, msg chan *message.Message) error {
	result, err := s.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(setting.SQS.QueueUrl),
		MaxNumberOfMessages: aws.Int64(setting.SQS.NumMessages),
		WaitTimeSeconds:     aws.Int64(setting.SQS.WaitTime),
		MessageAttributeNames: aws.StringSlice([]string{
			"All",
		}),
	})
	
	if err != nil {
		return err
	}
	
	for _, m := range result.Messages {
		msg <- message.NewMessage(m)
	}
	
	return nil
}

func readMessage(msg chan *message.Message, out chan error) (err error) {
	for {
		select {
		case err = <-out:
			return err
		case m := <-msg:
			clog.Info("[messageID:%s][topic:%s]", *m.MessageId, m.Topic())
			
			err := setting.WLog(m)
			if err != nil {
				return err
			}
		}
	}
}
