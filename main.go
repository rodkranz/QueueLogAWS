package main

import (
	"os"

	"github.com/go-clog/clog"
	"github.com/urfave/cli"

	"github.com/rodkranz/monitor/cmd"
	"github.com/rodkranz/monitor/pkg/setting"
)

func main() {
	app := cli.App{
		Name:   "Monitor",
		Before: setting.Setting,
		Commands: []cli.Command{
			cmd.SNS,
			cmd.SQS,
		},
		Flags: []cli.Flag{
			cli.StringFlag{Name: "aws_access_key", Usage: "AWS awsAccessKey", Value: os.Getenv("AWS_ACCESS_KEY_ID")},
			cli.StringFlag{Name: "secret_access_key", Usage: "AWS secretAccessKey", Value: os.Getenv("AWS_SECRET_ACCESS_KEY")},
			cli.StringFlag{Name: "region", Usage: "AWS region", Value: os.Getenv("AWS_REGION")},

			cli.StringFlag{Name: "log_driver", Usage: "log driver available are [sql, file]", Value: os.Getenv("LOG_DRIVER")},
			cli.StringFlag{Name: "log_param", Usage: "extra information maybe can help log's driver", Value: os.Getenv("LOG_PARAM")},
		},
	}

	app.Flags = append(app.Flags, []cli.Flag{}...)
	if err := app.Run(os.Args); err != nil {
		clog.Fatal(0, "exiting error: %v", err)
	}

	os.Exit(0)
}
