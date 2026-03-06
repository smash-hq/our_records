package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Minio    MinioConfig    `mapstructure:"minio"`
}

type ServerConfig struct {
	Port      string `mapstructure:"port"`
	Mode      string `mapstructure:"mode"`
	JWTSecret string `mapstructure:"jwt_secret"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type UploadConfig struct {
	Path    string `mapstructure:"path"`
	MaxSize int64  `mapstructure:"max_size"`
}

type MinioConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	Bucket          string `mapstructure:"bucket"`
	UseSSL          bool   `mapstructure:"use_ssl"`
}

var AppConfig *Config

func Init(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// 自动绑定环境变量（可选）
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 解析配置
	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		return err
	}

	log.Printf("配置加载成功：服务器端口=%s, 数据库=%s, MinIO=%s",
		AppConfig.Server.Port, AppConfig.Database.DBName, AppConfig.Minio.Endpoint)
	return nil
}
