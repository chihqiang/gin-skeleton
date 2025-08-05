package database

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 公共模型字段
type BaseModel struct {
	ID        uint           `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	CreatedAt time.Time      `gorm:"not null;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null;comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间(软删除)" json:"-"`
}

type PageRequest struct {
	Page int `json:"page" form:"page" param:"page" uri:"page" query:"page"`
	Size int `json:"size" form:"size" param:"size" uri:"size" query:"size"`
}

type PageResponse[T any] struct {
	//当前分页数
	CurrentPage int `json:"current_page" xml:"CurrentPage"`
	//当前拉去多少条
	Size int `json:"size" xml:"Size"`
	//总数
	Total int64 `json:"total" xml:"Total"`
	//列表
	Items []T `json:"items" xml:"Items"`
}

// Paginate
// https://gorm.io/zh_CN/docs/scopes.html#%E5%88%86%E9%A1%B5
// db := db.Model(Model{})
// paginate, _ := Paginate[*Model](db, context.Request)
func Paginate[T any](db *gorm.DB, req PageRequest) (*PageResponse[T], error) {
	resp := &PageResponse[T]{}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	resp.CurrentPage = req.Page
	resp.Size = req.Size
	var total int64
	_ = db.Count(&total).Error
	resp.Total = total
	var items []T
	err := db.Scopes(func(db *gorm.DB) *gorm.DB { return db.Offset((req.Page - 1) * req.Size).Limit(req.Size) }).Find(&items).Error
	resp.Items = items
	return resp, err
}

// IConversion 用于约束可转换为字符串的类型
type IConversion interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		string
}
