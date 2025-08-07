package queue

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testTask struct {
	ID  int    `json:"id"`
	Msg string `json:"msg"`
}

// Execute 模拟执行任务并通知外部（通过 context 中的 channel）
func (t *testTask) Execute(ctx context.Context, q IQueue) {
	execTime := time.Now().UnixMilli()
	fmt.Printf("[TestTask Execute] Msg: %s ExecutedAt: %d\n", t.Msg, execTime)

	// 通知外部任务执行
	if ch, ok := ctx.Value("notify").(chan int64); ok {
		ch <- execTime
	}
}
func setupRedisQueue(key string) (*RedisQueue, *redis.Client) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	rdb.FlushAll(context.Background())
	queue := NewRedisQueue(rdb, key)
	_ = queue.Register(&testTask{})
	return queue, rdb
}

func Test_ImmediateTaskExecution(t *testing.T) {
	queue, _ := setupRedisQueue("test:immediate")

	task := &testTask{ID: 1, Msg: "Immediate"}

	// 推送立即执行的任务
	assert.NoError(t, queue.Push(task, 0))

	notify := make(chan int64, 1)
	ctx, cancel := context.WithTimeout(context.WithValue(context.Background(), "notify", notify), 2*time.Second)
	defer cancel()

	go queue.Start(ctx, 50*time.Millisecond)

	select {
	case execTime := <-notify:
		assert.NotZero(t, execTime, "任务应该被立即执行")
	case <-time.After(1 * time.Second):
		t.Fatal("任务未在预期时间内执行")
	}
}

func Test_DelayedTaskExecution(t *testing.T) {
	queue, _ := setupRedisQueue("test:delayed")
	task := &testTask{ID: 2, Msg: "Delayed"}
	delay := 1 * time.Second
	start := time.Now().UnixMilli()
	assert.NoError(t, queue.Push(task, delay))

	notify := make(chan int64, 1)
	ctx, cancel := context.WithTimeout(context.WithValue(context.Background(), "notify", notify), 3*time.Second)
	defer cancel()

	go queue.Start(ctx, 50*time.Millisecond)

	select {
	case execTime := <-notify:
		cost := execTime - start
		assert.GreaterOrEqual(t, cost, int64(delay/time.Millisecond)-50, "任务应该延迟执行")
	case <-time.After(3 * time.Second):
		t.Fatal("延迟任务未在预期时间内执行")
	}
}

func BenchmarkQueueThroughput(b *testing.B) {
	queue, _ := setupRedisQueue("perf:queue")
	// 启动队列消费协程
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go queue.Start(ctx, 10*time.Millisecond)

	// 记录起始时间
	start := time.Now()

	// 模拟大量任务推送
	for i := 0; i < b.N; i++ {
		task := &testTask{ID: i}
		err := queue.Push(task, 0)
		if err != nil {
			b.Fatalf("Push failed: %v", err)
		}
	}

	// 等待所有任务执行完毕，最多等5秒
	b.Run("Wait", func(b *testing.B) {
		timeout := time.After(5 * time.Second)
		for {
			len, err := queue.rdb.ZCard(ctx, queue.taskKey).Result()
			if err != nil {
				b.Fatalf("ZCard failed: %v", err)
			}
			if len == 0 {
				break
			}
			select {
			case <-timeout:
				b.Fatalf("Timeout waiting for tasks to finish")
			default:
				time.Sleep(50 * time.Millisecond)
			}
		}
	})

	// 打印性能指标
	duration := time.Since(start)
	b.Logf("Processed %d tasks in %v (%.2f tasks/sec)", b.N, duration, float64(b.N)/duration.Seconds())
}
