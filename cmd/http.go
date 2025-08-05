package cmd

// 导入必要的包
import (
	"context"
	"github.com/urfave/cli/v3"
	"wangzhiqiang/skeleton/bootstrap"
)

// HTTPCommand 返回一个用于启动 HTTP 服务器的 CLI 命令
func HTTPCommand() *cli.Command {
	return &cli.Command{
		Name:  "http",
		Usage: "start http server",
		Flags: []cli.Flag{},
		Action: func(ctx context.Context, command *cli.Command) error {
			bootstrap.App(cfg).StartHTTP()
			return nil
		},
	}
}
