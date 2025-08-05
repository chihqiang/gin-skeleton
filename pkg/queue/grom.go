package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log/slog"
	"reflect"
	"sync"
	"time"
)

type SysTask struct {
	ID        uint      `gorm:"primaryKey"`
	Type      string    `gorm:"size:255;index"`
	Data      string    `gorm:"type:text"`
	RunAt     time.Time `gorm:"index"`
	ErrorMsg  string    `gorm:"type:text"`
	CreatedAt time.Time
}

type Gorm struct {
	db       *gorm.DB
	registry map[string]reflect.Type
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// NewGormQueue 创建队列实例
func NewGormQueue(db *gorm.DB) IQueue {
	_ = db.AutoMigrate(&SysTask{})
	ctx, cancel := context.WithCancel(context.Background())
	return &Gorm{
		db:       db,
		registry: make(map[string]reflect.Type),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Register 注册任务类型
func (q *Gorm) Register(task ITask) error {
	typeName, err := GetTaskTypeName(task)
	if err != nil {
		return err
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	q.registry[typeName] = reflect.TypeOf(task).Elem()
	return nil
}

// Push 推送任务
func (q *Gorm) Push(task ITask, delay time.Duration) error {
	typeName, err := GetTaskTypeName(task)
	if err != nil {
		return err
	}
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	model := SysTask{
		Type:  typeName,
		Data:  string(data),
		RunAt: time.Now().Add(delay),
	}
	return q.db.Create(&model).Error
}

func (q *Gorm) Pop() (ITask, error) {
	var model SysTask
	err := q.db.Transaction(func(tx *gorm.DB) error {
		// 查出一条可执行任务
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("run_at <= ?", time.Now()).
			Order("run_at ASC").
			First(&model).Error; err != nil {
			return err
		}
		// 删除任务
		if err := tx.Delete(&model).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	// 获取任务类型
	q.mu.RLock()
	typ, ok := q.registry[model.Type]
	q.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unregistered task type: %s", model.Type)
	}

	// 反序列化任务数据
	task := reflect.New(typ).Interface().(ITask)
	if err := json.Unmarshal([]byte(model.Data), task); err != nil {
		return nil, err
	}

	return task, nil
}

// Start 启动队列
func (q *Gorm) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-q.ctx.Done():
			slog.Warn("[GORM QUEUE] Stopped")
			return
		case <-ticker.C:
			task, err := q.Pop()
			if err != nil {
				slog.Warn("[GORM QUEUE] pop error", slog.Any("err", err))
				continue
			}
			if task == nil {
				continue
			}
			q.wg.Add(1)
			go func(task ITask) {
				err := task.Execute(ctx, q)
				if err != nil {
					slog.Warn("[GORM QUEUE] Execute task error", slog.Any("err", err))
				}
			}(task)
		}
	}
}

// Stop 停止队列
func (q *Gorm) Stop() {
	q.cancel()
	q.wg.Wait()
}
