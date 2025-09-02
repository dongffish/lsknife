package models

// Linkseeusers [...]
type Linkseeusers struct {
	ID           int    `gorm:"autoIncrement:true;primaryKey;column:id;type:int;not null" json:"id"`
	LoginName    string `gorm:"unique;column:login_name;type:varchar(20);default:null" json:"loginName"`
	UserName     string `gorm:"column:user_name;type:varchar(20);default:null" json:"userName"`
	UserPassword string `gorm:"column:user_password;type:varchar(64);default:null" json:"userPassword"`
	Roleid       int    `gorm:"column:roleid;type:int;default:null;default:0" json:"roleid"`
	Email1       string `gorm:"column:email1;type:varchar(40);default:null" json:"email1"`
	Email2       string `gorm:"column:email2;type:varchar(40);default:null" json:"email2"`
	Phone        string `gorm:"column:phone;type:varchar(11);default:null" json:"phone"`
	Telephone    string `gorm:"column:telephone;type:varchar(15);default:null" json:"telephone"`
	Weixin       string `gorm:"column:weixin;type:varchar(30);default:null" json:"weixin"`
	Isdomain     string `gorm:"column:isdomain;type:char(1);default:null;default:0" json:"isdomain"`
	Qq           string `gorm:"column:qq;type:varchar(15);default:null" json:"qq"`
	Lasttime     string `gorm:"column:lasttime;type:varchar(20);default:null" json:"lasttime"`
	Lastip       string `gorm:"column:lastip;type:varchar(15);default:null" json:"lastip"`
}

// TableName get sql table name.获取数据库表名
func (m *Linkseeusers) TableName() string {
	return "linkseeusers"
}

// LinkseeusersColumns get sql column name.获取数据库列名
var LinkseeusersColumns = struct {
	ID           string
	LoginName    string
	UserName     string
	UserPassword string
	Roleid       string
	Email1       string
	Email2       string
	Phone        string
	Telephone    string
	Weixin       string
	Isdomain     string
	Qq           string
	Lasttime     string
	Lastip       string
}{
	ID:           "id",
	LoginName:    "login_name",
	UserName:     "user_name",
	UserPassword: "user_password",
	Roleid:       "roleid",
	Email1:       "email1",
	Email2:       "email2",
	Phone:        "phone",
	Telephone:    "telephone",
	Weixin:       "weixin",
	Isdomain:     "isdomain",
	Qq:           "qq",
	Lasttime:     "lasttime",
	Lastip:       "lastip",
}
