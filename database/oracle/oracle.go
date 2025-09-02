package oracle

import (
	"errors"
	"gitee.com/dongffish/lsknife/database"
	oracle "github.com/seelly/gorm-oracle"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// GormOracle 初始化数据库。
// 参数PasswordDecryptionKey密码解密密钥,密钥为空时表示配置文件中的密码是明文.
// 密码加密是用 encrypt.TripleDesDecrypt 进行加密的
// 需要用一个toml类型的配置文件: config.toml
func GormOracle(iniSection, PasswordDecryptionKey string) (*gorm.DB, error) {
	if iniSection == "" {
		iniSection = "DATABASE"
	}
	dbCfg, err := database.ReadDBConfig(iniSection, PasswordDecryptionKey)
	if err != nil {
		return nil, errors.New("数据库配置文件config.toml读取失败!" + err.Error())
	}

	var gormCfg gorm.Config
	var customLogger logger.Interface
	gormCfg.DisableForeignKeyConstraintWhenMigrating = true

	// fmt.Println("dbCfg.ShowSQL =", dbCfg.ShowSQL)
	if dbCfg.ShowSQL == 1 { //显示所有SQL
		customLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				IgnoreRecordNotFoundError: true, // 忽略 record not found 错误
				LogLevel:                  logger.Info,
			},
		)
	} else {
		customLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				IgnoreRecordNotFoundError: true, // 忽略 record not found 错误
				LogLevel:                  logger.Error,
			},

		)
	}
	gormCfg.Logger = customLogger

	/*databaseURL := oracle.BuildUrl(server, port, service, user, password, nil)
	db, err := gorm.Open(oracle.Open(databaseURL),&gorm.Config{})*/

	db, err := gorm.Open(oracle.New(oracle.Config{
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

// GormOracleByDSN 根据DSN初始化数据库
func GormOracleByDSN(dsn string) (*gorm.DB, error) {
	var gormCfg gorm.Config
	db, err := gorm.Open(oracle.New(oracle.Config{
		DSN:               dsn, //dbCfg.ConnString,
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
