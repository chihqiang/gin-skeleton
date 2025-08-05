package database

// 导入必要的包
import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

// Init 根据配置初始化 GORM 数据库连接
// 参数 cfg: 数据库配置
// 返回值: 数据库连接实例或错误
func Init(cfg *Config) (*gorm.DB, error) {
	var (
		db  *gorm.DB // 数据库连接实例
		err error    // 错误变量
	)

	// 公共 GORM 配置
	gormCfg := &gorm.Config{
		QueryFields: true,
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名使用单数形式
		},
		Logger: logger.Default,
	}

	// 根据 Driver 字段选择不同数据库驱动
	switch cfg.Driver {
	case "mysql":
		// 构建 MySQL 数据源名称 (DSN)
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DBName,
		)
		// 打开 MySQL 连接
		db, err = gorm.Open(mysql.Open(dsn), gormCfg)
	case "sqlite":
		// SQLite 使用 DBName 作为文件路径
		db, err = gorm.Open(sqlite.Open(cfg.DBName), gormCfg)

	default:
		// 不支持的数据库驱动
		return nil, fmt.Errorf("unsupported driver: %s", cfg.Driver)
	}

	// 检查连接是否成功
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 配置连接池（仅对非 SQLite 有意义）
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed get DB: %w", err)
	}

	// 非 SQLite 数据库才需要配置连接池
	if cfg.Driver != "sqlite" {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)                                    // 设置最大空闲连接数
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)                                    // 设置最大打开连接数
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second) // 设置连接最大生存时间

		// 测试连接
		if err := sqlDB.Ping(); err != nil {
			return nil, fmt.Errorf("failed to ping database: %w", err)
		}
	}
	return db, nil
}
