package cmd

import (
	"context"
	"wangzhiqiang/skeleton/bootstrap"

	"github.com/urfave/cli/v3"
)

// QueueStartCommand 返回一个用于启动队列处理器的 CLI 命令
func QueueStartCommand() *cli.Command {
	return &cli.Command{
		Name:  "queue:start",
		Usage: "Start the queue worker to consume and execute pending tasks",
		Flags: []cli.Flag{},
		Action: func(ctx context.Context, command *cli.Command) error {
			bootstrap.App(cfg).StartQueue()
			return nil
		},
	}
}
