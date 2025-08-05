package cmd

import (
	"context"
	"wangzhiqiang/skeleton/config"

	"github.com/urfave/cli/v3"
	_ "wangzhiqiang/skeleton/routes"
)

var (
	cfg *config.Config
)

// Before 在执行 CLI 命令之前加载配置
// 该函数会在 CLI 命令执行前被调用，用于加载配置文件
func Before(ctx context.Context, cli *cli.Command) (context.Context, error) {
	var err error
	filename := cli.String(FlagConfig)
	cfg, err = config.Load(filename)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}
