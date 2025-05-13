package mysql

import (
	"errors"
	"fmt"
	"gitee.com/dongffish/lsknife/lsencrypt"
	"gitee.com/dongffish/lsknife/lslog"
	"github.com/go-ini/ini"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path/filepath"
	"time"
)

// GormMysql 初始化Mysql数据库。
// 参数PasswordDecryptionKey密码解密密钥,密钥为空时表示配置文件中的密码是明文.
// 密码加密是用 encrypt.TripleDesDecrypt 进行加密的
// 需要用一个ini类型的配置文件: dbconn.cfg
func GormMysql(iniSection, PasswordDecryptionKey string) (*gorm.DB, error) {
	if iniSection == "" {
		iniSection = "DATABASE"
	}
	dbCfg, err := ReadDBConfig(iniSection, PasswordDecryptionKey)
	if err != nil {
		return nil, errors.New("数据库配置文件 dbconn.cfg 读取失败!" + err.Error())
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
	if dbCfg.ShowSQL == 1 {
		customLogger.LogMode(logger.Info) // 打印 SQL 语句
	} else {
		customLogger.LogMode(logger.Error) // 只打印错误
	}
	gormCfg.Logger = customLogger

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dbCfg.ConnString,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置

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

// GormMysqlByDSN 根据DSN初始化Mysql数据库
func GormMysqlByDSN(dsn string) (*gorm.DB, error) {
	var gormCfg gorm.Config
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   //dbCfg.ConnString,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
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

// DBConfig 数据库配置
type DBConfig struct {
	DBType          string //`xml:"DBType,attr"`
	DBHost          string // 数据库服务器地址
	DBPort          int    // 数据库端口
	DBUser          string // 数据库用户名
	DBPassword      string // 数据库密码
	Database        string // 数据库名
	Charset         string // 数据库编码
	MaxOpenConns    int    // 50  用于设置最大打开的连接数，默认值为0表示不限制。
	MaxIdleConns    int    // 50  用于设置闲置的连接数
	ConnMaxLifetime int    // 秒	(100 * time.Second) , 给db设置一个超时时间，时间小于数据库的超时时间即可
	ShowSQL         int    //执行时是否显示SQL, 1-是  0-否
	//DiscardLog      bool   //`xml:"DiscardLog"`
	ConnString string // 数据库连接字符串
}

// ReadDBConfig 读数据库配置文件
func ReadDBConfig(dbSection, PasswordDecryptionKey string) (*DBConfig, error) {
	var err error
	// 获取可执行文件所在目录绝对路径
	exeDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, fmt.Errorf("获取执行目录失败: %v", err)
	}

	// 拼接配置文件完整路径
	cfgPath := filepath.Join(exeDir, "dbconn.cfg")

	// 加载 INI 配置文件
	cfg, err := ini.Load(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 获取指定分区的 DBType 配置
	section := cfg.Section(dbSection)
	if section == nil {
		return nil, errors.New("配置分区不存在: " + dbSection)
	}

	dbCfg := &DBConfig{
		DBType:  "mysql",
		ShowSQL: 0, //执行时是否显示SQL, 1-是  0-否
		//	DiscardLog:      true,
		MaxOpenConns:    100, //用于设置最大打开的连接数，值为0表示不限制
		MaxIdleConns:    100, //用于设置闲置的连接数,建议MaxOpenConns的值等于MaxIdleConns
		ConnMaxLifetime: 100, //通过SetConnMaxLifetime来关闭空闲连接 (秒).给db设置一个超时时间，时间小于数据库的超时时间即可
	}
	dbCfg.DBType = section.Key("DBType").String()
	dbCfg.Charset = section.Key("Charset").String()
	dbCfg.DBHost = section.Key("DBHost").String()
	dbCfg.DBPort, _ = section.Key("DBPort").Int()
	dbCfg.DBUser = section.Key("DBUser").String()
	dbCfg.DBPassword = section.Key("DBPassword").String()

	if PasswordDecryptionKey != "" {
		var b []byte
		if b, err = lsencrypt.TripleDesDecrypt([]byte(dbCfg.DBPassword), []byte(PasswordDecryptionKey)); err != nil {
			lslog.Logger.Errorln(lslog.GetLogPrefix(), "密码转换异常！")
			return nil, err
		}
		dbCfg.DBPassword = string(b)
	}

	dbCfg.Database = section.Key("Database").String()
	dbCfg.MaxOpenConns, _ = section.Key("MaxOpenConns").Int()
	dbCfg.MaxIdleConns, _ = section.Key("MaxIdleConns").Int()
	dbCfg.ConnMaxLifetime, _ = section.Key("ConnMaxLifetime").Int()
	dbCfg.ConnString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", dbCfg.DBUser, dbCfg.DBPassword, dbCfg.DBHost, dbCfg.DBPort, dbCfg.Database, dbCfg.Charset)
	dbCfg.ShowSQL, _ = section.Key("ShowSQL").Int() // ShowSQL 隐藏的配置参数,默认false，即不显示SQL语句
	return dbCfg, nil
}
