package database

import (
	"errors"
	"fmt"
	"gitee.com/dongffish/lsknife/lsencrypt"
	"gitee.com/dongffish/lsknife/lslog"
	"github.com/go-ini/ini"
	"github.com/spf13/viper"
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
	MaxOpenConns    int    // 100  这是最重要的参数，用于防止应用压垮数据库。切忌设置为0（无限制）。建议从 CPU核心数 × 2~3 起步，并确保不超过MySQL服务器 max_connections 的 80%
	MaxIdleConns    int    // 20 设置原则：核心目标是保持适量的“热连接”以减少建连开销。通常建议设置为 MaxOpenConns 的 25% ~ 50%。如果设置过高，空闲连接会长期占用数据库资源。
	ConnMaxLifetime int    // 1800秒，	(100 * time.Second) , 设置原则：主动轮换连接，避免使用因数据库或中间件（如ProxySQL、RDS）主动断开的“僵尸连接”。必须设置，且应小于数据库的 wait_timeout
	ConnMaxIdleTime int    // 300秒，连接最大空闲时间	(100 * time.Second) ,设置原则：这是对 MaxIdleConns 的精细化补充，用于及时释放闲置过久的连接，避免它们占用池子却因超时失效。建议设置为 5~10分钟
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
	// 设置Viper配置
	viper.SetConfigName("config") // 配置文件名（不带扩展名）
	viper.SetConfigType("toml")   // 配置文件类型
	// viper.AddConfigPath(".")      // 在目录查找配置文件
	viper.AddConfigPath(exeDir) // 在目录查找配置文件

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}
	// 读取特定配置项
	dbCfg := &DBConfig{
		MaxOpenConns:    100, //用于设置最大打开的连接数，值为0表示不限制
		MaxIdleConns:    100, //用于设置闲置的连接数,建议MaxOpenConns的值等于MaxIdleConns
		ConnMaxLifetime: 100, //通过SetConnMaxLifetime来关闭空闲连接 (秒).给db设置一个超时时间，时间小于数据库的超时时间即可
	}

	dbCfg.HostIP = viper.GetString("SYSTEM.HostIP")
	dbCfg.DBType = viper.GetString(dbSection + ".DBType")
	dbCfg.Charset = viper.GetString(dbSection + ".Charset")
	dbCfg.DBHost = viper.GetString(dbSection + ".DBHost")
	dbCfg.DBPort = viper.GetInt(dbSection + ".DBPort")
	dbCfg.DBUser = viper.GetString(dbSection + ".DBUser")
	dbCfg.DBPassword = viper.GetString(dbSection + ".DBPassword")
	dbCfg.ShowSQL = viper.GetInt(dbSection + ".ShowSQL") // ShowSQL 隐藏的配置参数,默认0，即不显示SQL语句
	dbCfg.Database = viper.GetString(dbSection + ".Database")
	dbCfg.MaxOpenConns = viper.GetInt(dbSection + ".MaxOpenConns")
	dbCfg.MaxIdleConns = viper.GetInt(dbSection + ".MaxIdleConns")
	dbCfg.ConnMaxLifetime = viper.GetInt(dbSection + ".ConnMaxLifetime")
	dbCfg.ConnMaxIdleTime = viper.GetInt(dbSection + ".ConnMaxIdleTime")
	
	if dbCfg.MaxOpenConns == 0 {
		dbCfg.MaxOpenConns = 100
	}
	if dbCfg.MaxIdleConns == 0 {
		dbCfg.MaxIdleConns = 20
	}
	if dbCfg.ConnMaxLifetime == 0 {
		dbCfg.ConnMaxLifetime = 1800
	}
	if dbCfg.ConnMaxIdleTime == 0 {
		dbCfg.ConnMaxIdleTime = 300
	}
	if PasswordDecryptionKey != "" {
		var b []byte
		if b, err = lsencrypt.TripleDesDecrypt([]byte(dbCfg.DBPassword), []byte(PasswordDecryptionKey)); err != nil {
			lslog.Logger.Errorln(lslog.GetLogPrefix(), "密码转换异常！")
			return nil, err
		}
		dbCfg.DBPassword = string(b)
	}

	switch dbCfg.DBType {
	case "mysql":
		dbCfg.ConnString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", dbCfg.DBUser, dbCfg.DBPassword, dbCfg.DBHost, dbCfg.DBPort, dbCfg.Database, dbCfg.Charset)
	case "pgsql":
		dbCfg.ConnString = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", dbCfg.DBHost, dbCfg.DBUser, dbCfg.DBPassword, dbCfg.Database, dbCfg.DBPort)
	case "mssql":
		dbCfg.ConnString = fmt.Sprintf("server=%s; user id=%s; password=%s; port=%d; database=%s;encrypt=disable", dbCfg.DBHost, dbCfg.DBUser, dbCfg.DBPassword, dbCfg.DBPort, dbCfg.Database)
	case "oracle":
		dbCfg.ConnString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbCfg.DBUser, dbCfg.DBPassword, dbCfg.DBHost, dbCfg.DBPort, dbCfg.Database)
	default:
		return nil, errors.New("暂不支持该数据库类型 DBType: " + dbCfg.DBType)
	}
	return dbCfg, nil
}

func ReadDBIniConfig(dbSection, PasswordDecryptionKey string) (*DBConfig, error) {
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
