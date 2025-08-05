package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"wangzhiqiang/skeleton/pkg/app"
	"wangzhiqiang/skeleton/pkg/queue"
)

type EmailTask struct {
	To       string `json:"to"`
	Body     string `json:"body"`
	Subject  string `json:"subject"`
	Attempts int    `json:"attempts"` // ✅ 建议改为 int（见下）
}

func (e *EmailTask) Execute(ctx context.Context, queue queue.IQueue) error {
	apps, err := app.GetApps(ctx)
	if err != nil {
		confJSON, err := json.MarshalIndent(apps.Config, "", "  ")
		if err != nil {
			fmt.Println("[EmailTask] ❌ 配置序列化失败:", err)
		} else {
			fmt.Printf("[EmailTask] ✅ 当前配置:\n%s\n", confJSON)
		}
	}

	// 日志格式优化
	fmt.Printf("[EmailTask] 📧 第 %d 次尝试发送邮件\n\t👉 收件人: %s\n\t👉 标题: %s\n\t👉 内容: %s\n", e.Attempts, e.To, e.Subject, e.Body)

	// TODO: 发送邮件逻辑

	// 示例：失败重试
	/*
		attempts, _ := strconv.Atoi(e.Attempts)
		if sendFailed {
			e.Attempts = strconv.Itoa(attempts + 1)
			_ = queue.Push(e, time.Second*5)
		}
	*/
	return nil
}
