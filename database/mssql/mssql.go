package mssql

import (
	"errors"
	"gitee.com/dongffish/lsknife/database"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// GormMssql 初始化Mysql数据库。
// 参数PasswordDecryptionKey密码解密密钥,密钥为空时表示配置文件中的密码是明文.
// 密码加密是用 encrypt.TripleDesDecrypt 进行加密的
// 需要用一个ini类型的配置文件: dbconn.cfg
func GormMssql(iniSection, PasswordDecryptionKey string) (*gorm.DB, error) {
	if iniSection == "" {
		iniSection = "DATABASE"
	}
	dbCfg, err := database.ReadDBConfig(iniSection, PasswordDecryptionKey)
	if err != nil {
		return nil, errors.New("数据库配置文件 dbconn.cfg 读取失败!" + err.Error())
	}

	var gormCfg gorm.Config
	gormCfg.DisableForeignKeyConstraintWhenMigrating = true
	customLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			IgnoreRecordNotFoundError: true, // 忽略 record not found 错误
		},
	)
	if dbCfg.ShowSQL == 1 {
		customLogger.LogMode(logger.Info) // 打印 SQL 语句
	} else {
		customLogger.LogMode(logger.Error) // 只打印错误
	}
	gormCfg.Logger = customLogger

	db, err := gorm.Open(sqlserver.New(sqlserver.Config{
		DSN:               dbCfg.ConnString,
		DefaultStringSize: 256, // string 类型字段的默认长度
	}), &gormCfg)
	if err != nil {
		return nil, errors.New("数据库Open失败! " + err.Error())
	}
	sqlDB, err := db.DB()
	if sqlDB == nil {
		return nil, errors.New("数据库nil")
	}
	if err != nil {
		return nil, errors.New("db.DB() Error ! " + err.Error())
	}
	if sqlDB.Ping() != nil {
		return nil, errors.New("数据库连接失败! ")
	}

	// 以下设置连接池信息
	sqlDB.SetMaxIdleConns(dbCfg.MaxIdleConns)                                    // 用于设置连接池中空闲连接的最在数量
	sqlDB.SetMaxOpenConns(dbCfg.MaxOpenConns)                                    // 用于打开数据连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(dbCfg.ConnMaxLifetime)) // 设置了连接可复用的最大时间(秒)

	return db, err
}

// GormMssqlByDSN 根据DSN初始化Mysql数据库
func GormMssqlByDSN(dsn string) (*gorm.DB, error) {
	var gormCfg gorm.Config
	db, err := gorm.Open(sqlserver.New(sqlserver.Config{
		DSN:               dsn,
		DefaultStringSize: 256, // string 类型字段的默认长度
	}), &gormCfg)

	if err != nil {
		return nil, errors.New("数据库Open失败! " + err.Error())
	}
	sqlDB, err := db.DB()
	if sqlDB == nil {
		return nil, errors.New("数据库nil")
	}
	if err != nil {
		return nil, errors.New("db.DB() Error ! " + err.Error())
	}
	if sqlDB.Ping() != nil {
		return nil, errors.New("数据库连接失败! ")
	}

	// 以下设置连接池信息
	MaxIdleConns := 600
	MaxOpenConns := 600
	ConnMaxLifetime := 150
	sqlDB.SetMaxIdleConns(MaxIdleConns)                                    // 用于设置连接池中空闲连接的最在数量
	sqlDB.SetMaxOpenConns(MaxOpenConns)                                    // 用于打开数据连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(ConnMaxLifetime)) // 设置了连接可复用的最大时间(秒)

	return db, err
}
