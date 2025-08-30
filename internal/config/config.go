// package config 专门用于处理应用的配置。
// 将配置相关的逻辑集中在一个包中，可以使代码结构更清晰，
// 其他包需要配置信息时，只需导入本包即可。
package config

import (
	"fmt"

	"github.com/spf13/viper" // 引入 viper 库，它是一个功能强大的配置解决方案
)

// Conf 是一个全局变量，用于存储从配置文件中加载的所有配置信息。
// 使用 new(Config) 初始化，使其成为一个指向 Config 结构体实例的指针。
// 这样，在其他包中就可以通过 config.Conf 直接访问配置项。
var Conf = new(Config)

// Config 结构体定义了应用的所有配置项。
// 它与 config.yaml 文件的结构一一对应。
// `mapstructure:"server"` 这种标签(tag)是给 viper 用的，
// 告诉 viper 在解析 YAML 文件时，如何将键(key)映射到结构体的字段(field)。
type Config struct {
	Server `mapstructure:"server"`
	MySQL  `mapstructure:"mysql"`
	Log    `mapstructure:"log"`
}

// Server 结构体定义了服务相关的配置。
type Server struct {
	Port      int    `mapstructure:"port"`       // 服务器监听的端口
	Mode      string `mapstructure:"mode"`       // Gin 框架的运行模式 (debug, test, release)
	JWTSecret string `mapstructure:"jwt_secret"` // JWT 签发密钥
}

// MySQL 结构体定义了数据库连接相关的配置。
type MySQL struct {
	Host     string `mapstructure:"host"`     // 数据库主机地址
	Port     int    `mapstructure:"port"`     // 数据库端口
	User     string `mapstructure:"user"`     // 用户名
	Password string `mapstructure:"password"` // 密码
	DBName   string `mapstructure:"dbname"`   // 数据库名称
	Charset  string `mapstructure:"charset"`  // 字符集
}

// Log 结构体定义了日志系统相关的配置。
type Log struct {
	Level      string `mapstructure:"level"`       // 日志级别 (debug, info, warn, error)
	Format     string `mapstructure:"format"`      // 日志格式 (console, json)
	Filename   string `mapstructure:"filename"`    // 日志文件名
	MaxSize    int    `mapstructure:"max_size"`    // 日志文件大小上限 (MB)
	MaxBackups int    `mapstructure:"max_backups"` // 最大备份数量
	MaxAge     int    `mapstructure:"max_age"`     // 最大保留天数
	Compress   bool   `mapstructure:"compress"`    // 是否压缩旧日志
}

// Init 函数负责初始化配置。它会在程序启动时被调用。
func Init() error {
	// 设置配置文件的名称（不带扩展名）
	viper.SetConfigName("config")
	// 设置配置文件的类型
	viper.SetConfigType("yaml")
	// 添加配置文件的搜索路径。可以多次调用以添加多个路径。
	viper.AddConfigPath("./configs")

	// 读取配置文件。如果找不到或格式错误，会返回 error。
	if err := viper.ReadInConfig(); err != nil {
		// 使用 fmt.Errorf 包装错误，提供更详细的上下文信息。
		return fmt.Errorf("read config failed: %w", err)
	}

	// 将读取到的配置信息反序列化（unmarshal）到全局的 Conf 变量中。
	// 这一步会将 YAML 中的数据填充到我们定义的 Config 结构体里。
	if err := viper.Unmarshal(Conf); err != nil {
		return fmt.Errorf("unmarshal config failed: %w", err)
	}

	// 如果一切顺利，返回 nil 表示没有错误。
	return nil
}
