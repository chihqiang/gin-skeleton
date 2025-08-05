package cmd

import (
	"github.com/urfave/cli/v3"
)

const (
	FlagConfig = "config" // 配置文件标志
)

// GlobalFlags 返回全局命令行标志
func GlobalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    FlagConfig,
			Aliases: []string{"c"},
			Value:   "config.yaml",
			Usage:   "Path to the configuration file",
		},
	}
}
