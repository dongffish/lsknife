package sqlite

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"

	// "gorm.io/driver/sqlite"  //　官方驱动使用了CGO实现，需要在CGO环境下才能工作，改用下面这个第三方的纯GO
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"time"
)

func GormSQLite(dbFile string) (*gorm.DB, error) {
	if dbFile == "" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return nil, err
		}
		dbFile = fmt.Sprintf("%s%c%s%c%s", dir, os.PathSeparator, "data", os.PathSeparator, "db.dat")
	}
	var gormCfg gorm.Config
	gormCfg.DisableForeignKeyConstraintWhenMigrating = true
	// 配置 logger，忽略 record not found 错误
	customLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			IgnoreRecordNotFoundError: true, // 忽略 record not found 错误
		},
	)
	customLogger.LogMode(logger.Error) // 只打印错误
	gormCfg.Logger = customLogger

	db, err := gorm.Open(sqlite.Open(dbFile), &gormCfg)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if sqlDB.Ping() == nil {

	}

	// 以下设置连接池信息
	sqlDB.SetMaxIdleConns(50)                   // 用于设置连接池中空闲连接的最在数量
	sqlDB.SetMaxOpenConns(50)                   // 用于打开数据连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Second * 100) // 设置了连接可复用的最大时间(秒)

	return db, err
}
