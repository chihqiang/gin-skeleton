package config

// 导入必要的包
import (
	"fmt"
	"gopkg.in/yaml.v3" // YAML 解析库
	"os"
	"path/filepath" // 文件路径处理
)

// Load 加载配置文件
// 接受多个路径，按顺序尝试加载，返回第一个成功加载的配置
// 参数 paths: 配置文件路径列表
// 返回值: 加载的配置或错误
func Load(paths ...string) (*Config, error) {
	var tried []string // 记录尝试过的路径
	for _, p := range paths {
		p = expandPath(p) // 扩展路径
		tried = append(tried, p)
		// 读取文件
		data, err := os.ReadFile(p)
		if err != nil {
			continue // 文件不存在则跳过
		}
		var cfg Config
		// 解析 YAML
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse configuration (%s): %w", p, err)
		}
		return &cfg, nil
	}
	// 没有找到有效的配置文件
	return nil, fmt.Errorf("no valid configuration file was found, try path: %v", tried)
}

// expandPath 扩展路径中的波浪号 (~)
// 将以 ~/ 开头的路径扩展为用户主目录下的对应路径
// 参数 p: 输入路径
// 返回值: 扩展后的路径
func expandPath(p string) string {
	if len(p) >= 2 && p[:2] == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return p
		}
		return filepath.Join(home, p[2:])
	}
	return p
}
