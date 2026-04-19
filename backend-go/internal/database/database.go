package database

import (
	"context"
	"fmt"
	"time"

	"github.com/mythicalsystems/nemesis-backend/internal/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

func Init(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DatabaseUser, cfg.DatabasePass, cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get DB object: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Initialize Redis
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPass,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := Redis.Ping(ctx).Result(); err != nil {
		// Redis failure is non-fatal — log and continue
		fmt.Printf("[WARN][DATABASE] Redis ping failed: %s\n", err)
	}

	return nil
}
