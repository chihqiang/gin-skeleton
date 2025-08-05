package main

// 导入必要的包
import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"wangzhiqiang/skeleton/cmd"

	"github.com/urfave/cli/v3"
)

// 全局变量定义
var (
	version  = "main" // 应用版本
	commands []*cli.Command
)

// 初始化函数
func init() {
	// 添加 HTTP 命令到命令列表
	commands = append(commands, cmd.HTTPCommand())
	commands = append(commands, cmd.QueueStartCommand())
}

// 主函数
func main() {
	// 创建 CLI 应用
	app := &cli.Command{}
	app.Name = "gin-skeleton"  // 应用名称
	app.Version = version      // 应用版本
	app.Usage = "Gin Skeleton" // 应用描述

	// 自定义版本打印函数
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Printf("gin-skeleton version %s %s/%s\n", cmd.Version, runtime.GOOS, runtime.GOARCH)
	}
	app.Flags = cmd.GlobalFlags()
	app.Commands = commands
	app.Before = cmd.Before
	// 运行应用
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
