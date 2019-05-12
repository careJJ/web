package routers

import (
	"pyg/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
    //发送短信
    beego.Router("/sendMsg",&controllers.UserController{},"post:HandleSendMsg")
	beego.Router("/register-email",&controllers.UserController{},"get:ShowEmail;post:HandleEmail")
	//激活用户
	beego.Router("/active",&controllers.UserController{},"get:Active")

	beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")

	beego.Router("/index",&controllers.GoodsController{},"get:ShowIndex")



    }
