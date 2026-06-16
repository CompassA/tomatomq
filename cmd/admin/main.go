/*
 * @Author: Tomato
 * @Date: 2026-05-10 21:29:53
 */
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/compassa/tomatomq/internal/admin/config"
	"github.com/compassa/tomatomq/internal/admin/handler"
	"github.com/compassa/tomatomq/internal/admin/meta"
	"github.com/compassa/tomatomq/internal/admin/midware"
	"github.com/compassa/tomatomq/internal/admin/mqadmin"
	tomatocfg "github.com/compassa/tomatomq/internal/pkg/config"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	ctx, cancel := context.WithCancel((context.Background()))
	defer cancel()

	// 加载配置
	env := tomatocfg.FetchEnv()
	cfg := config.LoadConfig(env)

	// 初始化日志
	rootLogger := tomatocfg.LoadLogger(env, &cfg.Log)

	// 初始化DB client
	db := initDBClient(cfg)
	rootLogger.Info("init DB client success")

	// 初始化etcd client
	etcdRepo := initEtcd(ctx, cfg)
	rootLogger.Info("init etcd client success")

	// 启动gin
	startGin(env, db, etcdRepo)
}

func initEtcd(ctx context.Context, cfg *config.Config) *meta.BrokerCacheRepo {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Etcd.Endpoints,
		DialTimeout: time.Duration(cfg.Etcd.DialTimoutSeconds) * time.Second,
	})
	if err != nil {
		panic(fmt.Errorf("failed to connect to ectd: %w", err))
	}

	metaRepo := meta.NewRepo(cli)

	err = metaRepo.StartWatch(ctx)
	if err != nil {
		panic(fmt.Errorf("BrokerCacheRepo#StartWatch: %w", err))
	}
	return metaRepo
}

func initDBClient(cfg *config.Config) *gorm.DB {
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
	return db
}

func startGin(env tomatocfg.Env, db *gorm.DB, etcdRepo *meta.BrokerCacheRepo) {
	if env == tomatocfg.AppBrokerEnvProd {
		gin.SetMode(gin.ReleaseMode)
	}

	h := handler.NewHandler(mqadmin.NewService(mqadmin.NewRepo(db), etcdRepo))
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
			h.TopicRegister(g)
		})
	}

	router.Run(tomatocfg.FetchPodIp() + ":8080")
}
