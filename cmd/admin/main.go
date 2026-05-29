/*
 * @Author: Tomato
 * @Date: 2026-05-10 21:29:53
 */
package main

import (
	"fmt"
	"time"

	"github.com/compassa/tomatomq/internal/admin/config"
	"github.com/compassa/tomatomq/internal/admin/handler"
	"github.com/compassa/tomatomq/internal/admin/midware"
	"github.com/compassa/tomatomq/internal/admin/mqadmin"
	tomatocfg "github.com/compassa/tomatomq/internal/pkg/config"
	"github.com/gin-gonic/gin"

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

	// 启动gin
	if env == tomatocfg.AppBrokerEnvProd {
		gin.SetMode(gin.ReleaseMode)
	}

	h := handler.NewHandler(mqadmin.NewService(mqadmin.NewRepo(db)))
	router := gin.Default()
	router.Use(midware.RequestUUIDHandler())
	router.Use(midware.LogReqRespHandler())
	router.Use(midware.ErrorHandler())
	{
		v1 := router.Group("/v1/mqadmin")
		v1.POST("/database/register", func(g *gin.Context) {
			h.DatabaseRegister(g)
		})
		v1.GET("/database/query", func(g *gin.Context) {
			h.DatabaseQueryByGroup(g)
		})

		v1.POST("/topic/register", func(g *gin.Context) {
		})
	}

	router.Run(tomatocfg.FetchPodIp() + ":8080")
}
