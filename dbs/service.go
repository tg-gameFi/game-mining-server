package dbs

import (
	"fmt"
	"game-mining-server/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
	"time"
)

type Service struct {
	DBInstance *gorm.DB
}

func initSqlIfNeed(cfg *configs.DatabaseConfig, dbName string) {
	_, e0 := os.Stat(cfg.InitPath)
	if os.IsNotExist(e0) {
		fmt.Printf("Init DB sql file not exists: %s, no need init sql", cfg.InitPath)
		return
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True", cfg.User, cfg.Pass, cfg.Host, cfg.Port)
	log.Printf("Init DB dsn %s\n", dsn)
	gormDB, e1 := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if e1 != nil {
		log.Panicf("Init DB create connection failed: %s", e1)
	}

	sqlDB, e2 := gormDB.DB()
	if e2 != nil {
		log.Panicf("Init DB get db instance failed: %s", e2)
	}

	raw, _ := os.ReadFile(cfg.InitPath)
	sqlList := strings.Split(string(raw), ";")
	for i, sqlStr := range sqlList {
		sql := strings.TrimSpace(sqlStr)
		if sql == "" {
			continue
		}
		_, e := sqlDB.Exec(strings.Replace(sql, "{db_name}", dbName, -1))
		if e != nil {
			log.Panicf("Init DB exec sql line %d file failed: %s", i, e)
		}
	}
	log.Printf("Init DB exec %d sqls\n", len(sqlList))
}

func createDBInstance(cfg *configs.DatabaseConfig, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", cfg.User, cfg.Pass, cfg.Host, cfg.Port, dbname)
	gormDB, e0 := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if e0 != nil {
		return nil, e0
	}

	sqlDB, e1 := gormDB.DB()
	if e1 != nil {
		return nil, e1
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return gormDB, nil
}

func CreateDBService(cfg *configs.DatabaseConfig, env string) *Service {
	dbName := cfg.Dbname + "_" + env // database name: {dbName}_dev or {dbName}_prod
	initSqlIfNeed(cfg, dbName)
	dbInstance, e2 := createDBInstance(cfg, dbName)
	if e2 != nil {
		log.Panicf("create DB connect failed: %s, host: %s", e2, cfg.Host)
	}
	return &Service{DBInstance: dbInstance}
}
