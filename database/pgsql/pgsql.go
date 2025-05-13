package pgsql

import (
	"errors"
	"fmt"
	"gitee.com/dongffish/lsknife/lsencrypt"
	"gitee.com/dongffish/lsknife/lslog"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path/filepath"
	"time"
)

// GormPGSql 初始化GormPGSql数据库。
// 参数PasswordDecryptionKey密码解密密钥,密钥为空时表示配置文件中的密码是明文.
// 密码加密是用 encrypt.TripleDesDecrypt 进行加密的
// 需要用一个ini类型的配置文件: dbconn.cfg
func GormPGSql(iniSection, PasswordDecryptionKey string) (*gorm.DB, error) {
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

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dbCfg.ConnString,
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

// GormPGSqlByDSN 根据DSN初始化数据库
func GormPGSqlByDSN(dsn string) (*gorm.DB, error) {
	var gormCfg gorm.Config
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
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
	ConfigPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	v := viper.New()
	v.AddConfigPath(ConfigPath)   // 路径
	v.SetConfigName("dbconn.cfg") // 名称
	v.SetConfigType("ini")        // 类型

	err = v.ReadInConfig() // 读配置
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return nil, errors.New("找不到配置文件")
		}
		return nil, err
	}

	dbCfg := &DBConfig{
		DBType: "pgsql",
		// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
		// user:password@tcp(localhost:5555)/dbname  // username@tcp(localhost)/my_db
		ShowSQL: 0, //执行时是否显示SQL, 1-是  0-否
		//	DiscardLog:      true,
		MaxOpenConns:    100, //用于设置最大打开的连接数，值为0表示不限制
		MaxIdleConns:    100, //用于设置闲置的连接数,建议MaxOpenConns的值等于MaxIdleConns
		ConnMaxLifetime: 100, //通过SetConnMaxLifetime来关闭空闲连接 (秒).给db设置一个超时时间，时间小于数据库的超时时间即可
	}
	dbCfg.DBType = v.GetString(dbSection + ".DBType")
	dbCfg.Charset = v.GetString(dbSection + ".Charset")
	dbCfg.DBHost = v.GetString(dbSection + ".DBHost")
	dbCfg.DBPort = v.GetInt(dbSection + ".DBPort")
	dbCfg.DBUser = v.GetString(dbSection + ".DBUser")
	dbCfg.DBPassword = v.GetString(dbSection + ".DBPassword")

	if PasswordDecryptionKey != "" {
		var b []byte
		if b, err = lsencrypt.TripleDesDecrypt([]byte(dbCfg.DBPassword), []byte(PasswordDecryptionKey)); err != nil {
			lslog.Logger.Errorln(lslog.GetLogPrefix(), "密码转换异常！")
			return nil, err
		}
		dbCfg.DBPassword = string(b)
	}

	dbCfg.Database = v.GetString(dbSection + ".Database")
	dbCfg.MaxOpenConns = v.GetInt(dbSection + ".MaxOpenConns")
	dbCfg.MaxIdleConns = v.GetInt(dbSection + ".MaxIdleConns")
	dbCfg.ConnMaxLifetime = v.GetInt(dbSection + ".ConnMaxLifetime")
	//dsn := "host=127.0.0.1 user=yourusername password=yourpassword dbname=userDB port=5432 sslmode=disable"
	dbCfg.ConnString = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", dbCfg.DBHost, dbCfg.DBUser, dbCfg.DBPassword, dbCfg.Database, dbCfg.DBPort)

	// ShowSQL 隐藏的配置参数,默认false，即不显示SQL语句
	dbCfg.ShowSQL = v.GetInt(dbSection + ".ShowSQL")
	return dbCfg, nil
}
