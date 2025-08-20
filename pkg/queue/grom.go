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

type Task struct {
	ID        uint      `gorm:"primaryKey"`
	Type      string    `gorm:"size:255;index"`
	Data      string    `gorm:"type:text"`
	Status    string    `gorm:"size:32;index"` // pending / processing / done / failed
	RunAt     time.Time `gorm:"index"`
	ErrorMsg  string    `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
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
	_ = db.AutoMigrate(&Task{})
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
	model := Task{
		Type:   typeName,
		Data:   string(data),
		Status: "pending",
		RunAt:  time.Now().Add(delay),
	}
	return q.db.Create(&model).Error
}

// Pop 获取一条任务
func (q *Gorm) Pop() (ITask, *Task, error) {
	var model Task
	err := q.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("status = ? AND run_at <= ?", "pending", time.Now()).
			Order("run_at ASC").
			First(&model).Error; err != nil {
			return err
		}
		return tx.Model(&model).Update("status", "processing").Error
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	q.mu.RLock()
	typ, ok := q.registry[model.Type]
	q.mu.RUnlock()
	if !ok {
		return nil, &model, fmt.Errorf("unregistered task type: %s", model.Type)
	}

	task := reflect.New(typ).Interface().(ITask)
	if err := json.Unmarshal([]byte(model.Data), task); err != nil {
		return nil, &model, err
	}
	return task, &model, nil
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
		case <-ctx.Done():
			slog.Warn("[GORM QUEUE] Context canceled")
			return
		case <-ticker.C:
			task, model, err := q.Pop()
			if err != nil {
				slog.Warn("[GORM QUEUE] pop error", slog.Any("err", err))
				continue
			}
			if task == nil {
				continue
			}
			q.wg.Add(1)
			go func(task ITask, model *Task) {
				err := task.Execute(ctx, q)
				if err != nil {
					_ = q.db.Model(model).Updates(map[string]interface{}{
						"status":    "failed",
						"error_msg": err.Error(),
					}).Error
				} else {
					_ = q.db.Model(model).Update("status", "done").Error
				}
			}(task, model)
		}
	}
}

// Stop 停止队列
func (q *Gorm) Stop() {
	q.cancel()
	q.wg.Wait()
}
