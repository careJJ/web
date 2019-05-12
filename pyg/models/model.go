package models

import ("github.com/astaxie/beego/orm"
_ "github.com/go-sql-driver/mysql"
)
type User struct {
	Id int
	Name string	`orm:"size(40)"`
	Pwd string	`orm:"size(40)"`
	Phone string	`orm:"size(11)"`
	Email string	`orm:"null"`
	Active bool		`orm:"default(false)"`
	Addresses []*Address 	`orm:"reverse(many)"`
}

type Address struct {
	Id int
	Receiver string	`orm:"size(40)"`
	Addr string	`orm:"size(100)"`
	PostCode string
	Phone string	`orm:"size(11)"`
	IsDefault bool `orm:"default(false)"`
	User *User 	`orm:"rel(fk)"`
}

func init(){
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/pyg")
	orm.RegisterModel(new(User),new(Address))
	orm.RunSyncdb("default",false,true)

}