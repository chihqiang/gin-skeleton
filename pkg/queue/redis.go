package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"reflect"
	"sync"
	"time"
)

// taskWrapper 用于在 Redis 中存储任务的统一封装格式
// 主要是为了能在取出任务时，知道任务类型并反序列化成对应的结构体
type taskWrapper struct {
	Type string          `json:"type"` // 任务的唯一类型标识（包路径 + 类型名）
	Data json.RawMessage `json:"data"` // 任务的具体数据（原始 JSON）
}

// RedisQueue 基于 Redis 实现的延时任务队列
type RedisQueue struct {
	rdb      *redis.Client           // Redis 客户端
	taskKey  string                  // 普通任务队列的 Redis ZSet key
	registry map[string]reflect.Type // 已注册的任务类型映射（typeName -> 反射类型）
	mu       sync.RWMutex            // 读写锁，保护 registry 并发安全
	ctx      context.Context         // 上下文，用于 Redis 操作取消/超时
	cancel   context.CancelFunc
}

func NewRedisQueue(rdb *redis.Client, queue string) *RedisQueue {
	ctx, cancel := context.WithCancel(context.Background())
	return &RedisQueue{
		rdb:      rdb,
		taskKey:  queue,
		registry: make(map[string]reflect.Type),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Register 注册任务类型（消费端必须先注册）
// 参数 task 必须是结构体指针类型，例如：&MyTask{}
func (q *RedisQueue) Register(task ITask) error {
	typeName, err := GetTaskTypeName(task) // 获取任务类型的唯一标识
	if err != nil {
		return err
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	// 使用 Elem() 是为了获取结构体类型（而不是指针类型）
	q.registry[typeName] = reflect.TypeOf(task).Elem()
	return nil
}

func (q *RedisQueue) Start(val context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-q.ctx.Done():
			slog.Warn("[REDIS QUEUE] Stopped")
			return
		case <-ticker.C:
			task, err := q.Pop()
			if err != nil {
				continue
			}
			if task == nil {
				continue
			}
			if err := task.Execute(val, q); err != nil {
				slog.Warn("[REDIS QUEUE] Execute task error", slog.Any("err", err))
			}
		}
	}
}

// Stop 停止队列
func (q *RedisQueue) Stop() {
	q.cancel()
}

// Push 推送普通任务到队列
// delay > 0 表示延迟执行
func (q *RedisQueue) Push(task ITask, delay time.Duration) error {
	return q.enqueueTask(q.taskKey, task, delay)
}
func (q *RedisQueue) Pop() (ITask, error) {
	now := time.Now().UnixMilli()
	var popLuaScript = redis.NewScript(`
    -- KEYS[1] 队列 key
    -- ARGV[1] 当前时间（毫秒）
    local items = redis.call("ZRANGEBYSCORE", KEYS[1], 0, ARGV[1], "LIMIT", 0, 1)
    if #items == 0 then
        return nil
    end
    redis.call("ZREM", KEYS[1], items[1])
    return items[1]
`)
	// 执行 Lua 脚本
	res, err := popLuaScript.Run(q.ctx, q.rdb, []string{q.taskKey}, now).Result()
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil // 没有任务
	}
	raw, ok := res.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected type: %T", res)
	}
	// 解析任务封装
	var wrap taskWrapper
	if err := json.Unmarshal([]byte(raw), &wrap); err != nil {
		return nil, err
	}
	q.mu.RLock()
	typ, ok := q.registry[wrap.Type]
	q.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unregistered task type: %s", wrap.Type)
	}
	// 反序列化任务数据
	task := reflect.New(typ).Interface().(ITask)
	if err := json.Unmarshal(wrap.Data, task); err != nil {
		return nil, err
	}
	return task, nil
}

// enqueueTask 将任务序列化并添加到 Redis ZSet 队列中
// key   - 队列 Redis key
// delay - 延迟时间，score 为任务可执行的时间戳（毫秒）
func (q *RedisQueue) enqueueTask(key string, task ITask, delay time.Duration) error {
	typeName, err := GetTaskTypeName(task) // 获取任务类型标识
	if err != nil {
		return err
	}
	// 序列化任务数据
	b, _ := json.Marshal(task)
	// 再封装一层，包含任务类型
	wrap, _ := json.Marshal(taskWrapper{Type: typeName, Data: b})
	// 将任务放入 Redis ZSet，score 是执行时间戳
	return q.rdb.ZAdd(q.ctx, key, redis.Z{
		Score:  float64(time.Now().Add(delay).UnixMilli()), // 延迟执行时间
		Member: wrap,                                       // 序列化后的任务数据
	}).Err()
}
