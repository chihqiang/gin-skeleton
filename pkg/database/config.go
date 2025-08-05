package database

// Config 数据库配置结构体
// 包含数据库连接所需的各种配置参数
type Config struct {
	Driver          string `yaml:"driver" json:"driver,omitempty"`                       // 数据库驱动（mysql/sqlite）
	Host            string `yaml:"host" json:"host,omitempty"`                           // 数据库主机地址
	Port            int    `yaml:"port" json:"port,omitempty"`                           // 数据库端口
	Username        string `yaml:"username" json:"username,omitempty"`                   // 数据库用户名
	Password        string `yaml:"password" json:"password,omitempty"`                   // 数据库密码
	DBName          string `yaml:"dbname" json:"dbname,omitempty"`                       // 数据库名称
	MaxIdleConns    int    `yaml:"max_idle_conns" json:"max_idle_conns,omitempty"`       // 最大空闲连接数
	MaxOpenConns    int    `yaml:"max_open_conns" json:"max_open_conns,omitempty"`       // 最大打开连接数
	ConnMaxLifetime int    `yaml:"conn_max_lifetime" json:"conn_max_lifetime,omitempty"` // 连接最大生存时间（秒）
}
