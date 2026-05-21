/*
 * @Author: Tomato
 * @Date: 2026-05-10 21:29:53
 */
package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/compassa/tomatomq/internal/admin/config"
	"github.com/compassa/tomatomq/internal/admin/mqadmin"
	tomatocfg "github.com/compassa/tomatomq/internal/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	env := tomatocfg.FetchEnv()
	cfg := config.LoadConfig(env)
	defer config.SyncLogger()

	// 初始化DB client
	db, err := gorm.Open(mysql.Open(cfg.Database.Dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("open mysql failed: %w", err))
	}
	sqldb, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("fetch sqldb failed: %w", err))
	}
	sqldb.SetConnMaxIdleTime(time.Duration(cfg.Database.ConnMaxLifeMinutes) * time.Minute)
	sqldb.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqldb.SetMaxOpenConns(cfg.Database.MaxOpenConns)

	adminRepo := mqadmin.NewRepo(db)

	res := adminRepo.QueryById(1)

	for _, r := range res {
		v, _ := json.Marshal(r)
		config.AppLogger.Info("test info", slog.Any("row", string(v)))
	}
}
