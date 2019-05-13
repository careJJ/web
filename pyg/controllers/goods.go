package controllers

import "github.com/astaxie/beego"

type GoodsController struct {
	beego.Controller
}


func (this*GoodsController)ShowIndex() {
	name:=this.GetSession("name")
	if name!=nil{
		this.Data["name"]=name.(string)
	}
	
	this.TplName="/index"
}
