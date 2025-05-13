package main

import (
	"fmt"
	"gitee.com/dongffish/lsknife/database/mysql"
	"gitee.com/dongffish/lsknife/database/pgsql"
	"gitee.com/dongffish/lsknife/database/sqlite"
	"gitee.com/dongffish/lsknife/lslog"
)

func main() {
	//MySQLTest()
	//SQLiteTest()
	PGSQLTest()
}

func MySQLTest() {
	lslog.Logger.Infoln(lslog.GetLogPrefix(), "Hello MYSql !")

	db, err := mysql.GormMysql("DATABASE", "")
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

	db, err := pgsql.GormPGSql("DATABASE", "")
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
