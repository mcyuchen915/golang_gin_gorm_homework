package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	DbType   string         `mapstructure:"dbType"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire string `mapstructure:"expire"`
}

func init() {
	// 设置配置文件名称（不含扩展名）
	viper.SetConfigName("config")
	// 设置配置文件类型
	viper.SetConfigType("yaml")
	// 添加配置文件搜索路径
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.app")

	// 读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	// 设置默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.mode", "debug")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Error reading config file: %v", err)
		log.Println("Using default values and environment variables")
	} else {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}
}

// LoadConfig 加载并解析配置文件，将 Viper 中的配置数据映射到 Config 结构体
// 返回解析后的配置对象指针和可能的错误
// 配置数据来源包括：配置文件、环境变量和默认值（在 init 函数中已设置）
func LoadConfig() (*Config, error) {
	var config Config

	// 使用 viper.Unmarshal 将配置数据解析到 config 结构体中
	// 如果解析失败，返回错误
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// 返回解析成功的配置对象
	return &config, nil
}
func getSqliteFilePath(filename string) string {
	// Get the db directory where SQLite files should be stored
	dbDir, err := getDBDir()
	if err != nil {
		return "test.sqlite.db"
	}

	// Ensure filename contains "sqlite"
	if filename == "" {
		filename = "test.sqlite.db"
	} else {
		// Extract base name and extension
		ext := filepath.Ext(filename)
		if ext == "" {
			ext = ".db"
		}
		base := filename[:len(filename)-len(ext)]
		if base == "" {
			base = "test"
		}
		// Check if "sqlite" is already in the filename (case-insensitive check)
		baseLower := strings.ToLower(base)
		if !strings.Contains(baseLower, "sqlite") {
			filename = base + "_sqlite" + ext
		} else {
			filename = base + ext
		}
	}

	// Database file will be stored in db directory (examples/db)
	dbPath := filepath.Join(dbDir, filename)
	return dbPath
}
func getDBDir() (string, error) {
	// Get the directory where this file (config.go) is located
	// This file is in homework/config directory
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", os.ErrInvalid
	}

	// Get the homework directory (parent of config)
	configDir := filepath.Dir(currentFile)
	homeworkDir := filepath.Dir(configDir)

	// The db directory is homework/db
	dbDir := filepath.Join(homeworkDir, "db")

	// Ensure db directory exists
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return "", err
	}

	return dbDir, nil
}

func GetConfigDB() (*gorm.DB, error) {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return nil, err
	}
	if cfg.Server.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	var db *gorm.DB
	if cfg.DbType == "sqlite" || cfg.DbType == "" {

		// 初始化数据库
		db, err = gorm.Open(sqlite.Open(getSqliteFilePath("")), &gorm.Config{
			// Silent for tests, use logger.Info for development
			Logger: logger.Default.LogMode(logger.Info),
			// NamingStrategy: Customize how GORM names tables and columns
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "",    // Prefix for all table names (e.g., "app_")
				SingularTable: false, // Use singular table names (User -> user instead of users)
				NoLowerCase:   false, // Disable automatic lowercasing
				NameReplacer:  nil,   // Custom name replacer function
			},
		})
		if err != nil {
			log.Fatalf("Failed to connect database: %v", err)
			return nil, err
		}

	} else {
		// "root:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.DBName,
		)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			// Logger: Set to logger.Info to see all SQL queries in development
			// Logger: logger.Default.LogMode(logger.Silent),
			Logger: logger.Default.LogMode(logger.Info),
			// NamingStrategy: Customize table and column naming
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "",
				SingularTable: false,
				NoLowerCase:   false,
			},
		})
		if err != nil {
			log.Fatalf("Failed to connect database: %v", err)
			return nil, err
		}
	}
	// 自动迁移
	if err := db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	// 初始化User数据
	if err := InitUserData(db); err != nil {
		log.Fatalf("Failed to init user data: %v", err)
	}
	return db, nil
}

// InitUserData 初始化User数据
func InitUserData(db *gorm.DB) error {
	User := []models.User{
		{Username: "Alice", Email: "Alice@example.com>"},
		{Username: "Sarah", Email: "Sarah@example.com>"},
		{Username: "David", Email: "David@example.com>"},
	}
	db.Unscoped().Where("username = ?", "Alice").Or("username = ?", "Sarah").Or("username = ?", "David").Delete(&models.User{})
	if err := db.Create(&User).Error; err != nil {
		return err
	}
	return nil
}
