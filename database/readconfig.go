package database

import (
	"errors"
	"fmt"
	"gitee.com/dongffish/lsknife/lsencrypt"
	"gitee.com/dongffish/lsknife/lslog"
	"github.com/go-ini/ini"
	"os"
	"path/filepath"
)

// DBConfig 数据库配置
type DBConfig struct {
	DBType          string // mysql、pgsql、mssql
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
	HostIP     string // 本机IP，放在[SYSTEM]段
}

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

	dbCfg := &DBConfig{
		MaxOpenConns:    100, //用于设置最大打开的连接数，值为0表示不限制
		MaxIdleConns:    100, //用于设置闲置的连接数,建议MaxOpenConns的值等于MaxIdleConns
		ConnMaxLifetime: 100, //通过SetConnMaxLifetime来关闭空闲连接 (秒).给db设置一个超时时间，时间小于数据库的超时时间即可
	}
	section := cfg.Section("SYSTEM")
	if section != nil {
		dbCfg.HostIP = section.Key("HostIP").String()
	}

	// 获取指定分区的 DBType 配置
	section = cfg.Section(dbSection)
	if section == nil {
		return nil, errors.New("配置分区不存在: " + dbSection)
	}

	dbCfg.DBType = section.Key("DBType").String()
	if dbCfg.DBType == "" {
		dbCfg.DBType = "mysql"
	}
	dbCfg.Charset = section.Key("Charset").String()
	dbCfg.DBHost = section.Key("DBHost").String()
	dbCfg.DBPort, err = section.Key("DBPort").Int()
	if err != nil {
		return nil, errors.New("数据库端口设置有误:DBPort")
	}
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
	dbCfg.ShowSQL, err = section.Key("ShowSQL").Int() // ShowSQL 隐藏的配置参数,默认0，即不显示SQL语句
	if err != nil {
		dbCfg.ShowSQL = 0
	}
	dbCfg.Database = section.Key("Database").String()
	dbCfg.MaxOpenConns, _ = section.Key("MaxOpenConns").Int()
	dbCfg.MaxIdleConns, _ = section.Key("MaxIdleConns").Int()
	dbCfg.ConnMaxLifetime, _ = section.Key("ConnMaxLifetime").Int()
	switch dbCfg.DBType {
	case "mysql":
		dbCfg.ConnString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", dbCfg.DBUser, dbCfg.DBPassword, dbCfg.DBHost, dbCfg.DBPort, dbCfg.Database, dbCfg.Charset)
	case "pgsql":
		dbCfg.ConnString = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", dbCfg.DBHost, dbCfg.DBUser, dbCfg.DBPassword, dbCfg.Database, dbCfg.DBPort)
	case "mssql":
		dbCfg.ConnString = fmt.Sprintf("server=%s; user id=%s; password=%s; port=%d; database=%s;encrypt=disable", dbCfg.DBHost, dbCfg.DBUser, dbCfg.DBPassword, dbCfg.DBPort, dbCfg.Database)
	case "oracle":
		//dsn := "user:password@tcp(host:port)/dbname"
		dbCfg.ConnString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbCfg.DBUser, dbCfg.DBPassword, dbCfg.DBHost, dbCfg.DBPort, dbCfg.Database)
	default:
		return nil, errors.New("暂不支持该数据库类型 DBType: " + dbCfg.DBType)
	}
	return dbCfg, nil
}
