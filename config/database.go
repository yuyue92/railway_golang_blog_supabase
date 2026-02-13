package config

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	dsn := strings.TrimSpace(os.Getenv("SUPABASE_DB_URL"))
	if dsn == "" {
		dsn = strings.TrimSpace(os.Getenv("DATABASE_URL"))
	}
	if dsn == "" {
		return fmt.Errorf("SUPABASE_DB_URL / DATABASE_URL 未设置")
	}

	dsn = ensureSSLModeRequire(dsn)

	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,

		// ✅ 对 Supabase pooler/pgbouncer 更稳
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("gorm.Open 失败: %w", err)
	}

	// ✅ 启动时就验证 DB 可用，不要等到请求时才爆炸
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取底层 sql.DB 失败: %w", err)
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("数据库 Ping 失败(很可能是连接串/网络/SSL问题): %w", err)
	}

	log.Println("✅ PostgreSQL 连接 & Ping 成功")
	return nil
}

func GetDB() *gorm.DB { return DB }

func ensureSSLModeRequire(dsn string) string {
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		u, err := url.Parse(dsn)
		if err != nil {
			return dsn
		}
		q := u.Query()
		if q.Get("sslmode") == "" {
			q.Set("sslmode", "require")
			u.RawQuery = q.Encode()
		}
		return u.String()
	}
	if !strings.Contains(dsn, "sslmode=") {
		return dsn + " sslmode=require"
	}
	return dsn
}
