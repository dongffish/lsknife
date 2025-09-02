package main

import (
	"fmt"
	"gitee.com/dongffish/lsknife/database"
	"gitee.com/dongffish/lsknife/database/mssql"
	"gitee.com/dongffish/lsknife/database/mysql"
	"gitee.com/dongffish/lsknife/database/pgsql"
	"gitee.com/dongffish/lsknife/database/sqlite"
	"gitee.com/dongffish/lsknife/lslog"
	"gitee.com/dongffish/lsknife/test/models"
)

var PasswordKeyDefault = "YyHz*sH_Yh!Hxy#528866838@qq.com;$China"

func main() {
	dbcfg, err := database.ReadDBConfig("DATABASE", PasswordKeyDefault)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ShowSQL =", dbcfg.ShowSQL)
	if dbcfg.DBType == "mysql" {
		MySQLTest()
	}
	if dbcfg.DBType == "pgsql" {
		PGSQLTest()
	}
	if dbcfg.DBType == "mssql" {
		MssqlTest()
	}
	if dbcfg.DBType == "oracle" {
		OracleTest()
	}
}

func MySQLTest() {
	lslog.Logger.Infoln(lslog.GetLogPrefix(), "Hello MYSql !")
	db, err := mysql.GormMysql("DATABASE", PasswordKeyDefault)
	if err != nil {
		lslog.Logger.Errorln(err)
		return
	}
	fmt.Println("db open success !", db)
	var tb []models.Linkseeusers

	db.Model(&models.Linkseeusers{}).Find(&tb)
	for _, v := range tb {
		fmt.Println(v)
	}

	defer func() {
		db1, _ := db.DB() // 程序结束前关闭数据库链接
		_ = db1.Close()
		fmt.Println("db closed")
	}()

}
func OracleTest() {
	lslog.Logger.Infoln(lslog.GetLogPrefix(), "Hello Oracle Database !")
	db, err := mysql.GormMysql("DATABASE", PasswordKeyDefault)
	if err != nil {
		lslog.Logger.Errorln(err)
		return
	}
	fmt.Println("db open success !", db)
	defer func() {
		db1, _ := db.DB() // 程序结束前关闭数据库链接
		_ = db1.Close()
		fmt.Println("db closed")
	}()

}

func MssqlTest() {
	lslog.Logger.Infoln(lslog.GetLogPrefix(), "Hello MicroSoft SQL Server !")
	db, err := mssql.GormMssql("DATABASE", PasswordKeyDefault)
	if err != nil {
		lslog.Logger.Errorln(err)
		return
	}
	fmt.Println("db open success !", db)
	defer func() {
		db1, _ := db.DB() // 程序结束前关闭数据库链接
		_ = db1.Close()
		fmt.Println("db closed")
	}()

}

type DahuaVideo struct {
	ID        int    `gorm:"autoIncrement:true;primaryKey;column:id;type:int;not null" json:"id"`
	Hostip    string `gorm:"index:dahua_video_idx1;column:hostip;type:varchar(15);not null" json:"hostip"`
	Username  string `gorm:"column:username;type:varchar(20);not null" json:"username"`
	Password  string `gorm:"column:password;type:varchar(32);not null" json:"password"`
	Result    string `gorm:"column:result;type:varchar(8);not null" json:"result"`
	Msg       string `gorm:"column:msg;type:varchar(128);default:null" json:"msg"`
	RetCode   string `gorm:"column:ret_code;type:varchar(6);default:'0'" json:"retCode"`
	LoginMode string `gorm:"column:login_mode;type:varchar(2);default:'0'" json:"loginMode"`
	Lasttime  string `gorm:"column:lasttime;type:varchar(20);default:null" json:"lasttime"`
}

func SQLiteTest() {
	lslog.Logger.Infoln(lslog.GetLogPrefix(), "Hello SQLite!")

	db, err := sqlite.GormSQLite("")
	if err != nil {
		lslog.Logger.Errorln(err)
		return
	}
	var dh DahuaVideo

	err = db.AutoMigrate(&dh)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("db open success !", db)
	defer func() {
		db1, _ := db.DB() // 程序结束前关闭数据库链接
		_ = db1.Close()
		fmt.Println("db closed")
	}()

}

func PGSQLTest() {
	lslog.Logger.Infoln(lslog.GetLogPrefix(), "Hello PQSQL!")

	db, err := pgsql.GormPGSql("DATABASE", PasswordKeyDefault)
	if err != nil {
		lslog.Logger.Errorln(err)
		return
	}
	var dh DahuaVideo

	err = db.AutoMigrate(&dh)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("db open success !", db)
	defer func() {
		db1, _ := db.DB() // 程序结束前关闭数据库链接
		_ = db1.Close()
		fmt.Println("db closed")
	}()

}
