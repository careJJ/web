package controllers

import (
	"github.com/astaxie/beego"
	"regexp"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"encoding/json"
	"github.com/astaxie/beego/utils"
	"github.com/astaxie/beego/orm"
	"pyg/models"

)

type UserController struct {
	beego.Controller
}

type Message struct {
	Message string
	RequestId string
	BizId string
	Code string
}


func RespFunc(this* UserController,resp map[string]interface{}){
	//3.把容器传递给前段
	this.Data["json"] = resp
	//4.指定传递方式
	this.ServeJSON()
}

func (this*UserController)ShowRegister() {
	this.TplName="register.html"
}

func (this*UserController)HandleSendMsg() {
	phone:=this.GetString("phone")
	resp:=make(map[string]interface{})
	defer RespFunc(this,resp)

	if phone=="" {
		beego.Error("获取电话号码失败")
		resp["errno"]=1
		resp["errmsg"]="获取电话号码错误"
		this.Data["json"]=resp
		this.ServeJSON()
		return
	}
	reg,_:=regexp.Compile(`1[3-9][0-9]{9}$`)
	result:=reg.FindString(phone)
	if result=="" {
		beego.Error("电话号码格式错误")
		resp["errno"]=2
		resp["errmsg"]="电话号码格式错误"
		return
	}
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", "LTAIunREOfuKN7dT", "JfixJMggEIGMsdBdg9Y4pweWMF6Rn7" +
		"")
	if err != nil {
		beego.Error("电话号码格式错误")
		//2.给容器赋值
		resp["errno"] = 3
		resp["errmsg"] = "初始化短信错误"
		return
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["PhoneNumbers"] = phone
	request.QueryParams["SignName"] = "品优购"
	request.QueryParams["TemplateCode"] = "SMS_165107660"
	request.QueryParams["TemplateParam"] = "{\"code\":123456}"

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		beego.Error("电话号码格式错误")
		//2.给容器赋值
		resp["errno"] = 4
		resp["errmsg"] = "短信发送失败"
		return
	}
	//json数据解析
	var message Message
	json.Unmarshal(response.GetHttpContentBytes(),&message)
	if message.Message != "OK"{
		beego.Error("电话号码格式错误")
		//2.给容器赋值
		resp["errno"] = 6
		resp["errmsg"] = message.Message
		return
	}
	
	resp["errno"] = 5
	resp["errmsg"] = "发送成功"



}

func(this*UserController)HandleRegister(){
	//获取数据
	phone := this.GetString("phone")
	pwd :=this.GetString("password")
	rpwd := this.GetString("repassword")
	//校验数据
	if phone == "" || pwd == "" || rpwd == ""{
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "获取数据错误"
		this.TplName = "register.html"
		return
	}
	if pwd != rpwd{
		beego.Error("两次密码输入不一致")
		this.Data["errmsg"] = "两次密码输入不一致"
		this.TplName = "register.html"
		return
	}
	//处理数据
	//orm插入数据
	o := orm.NewOrm()
	var user models.User
	user.Name = phone
	user.Pwd = pwd
	user.Phone = phone
	o.Insert(&user)
	//激活页面
	this.Ctx.SetCookie("userName",user.Name,60 * 10)
	this.Redirect("/register-email",302)
	//返回数据
}

//展示邮箱激活
func(this*UserController)ShowEmail(){
	this.TplName = "register-email.html"
}

//处理邮箱激活业务
func(this*UserController)HandleEmail(){
	//获取数据
	email := this.GetString("email")
	pwd := this.GetString("password")
	rpwd := this.GetString("repassword")
	//校验数据
	if email == "" || pwd == ""|| rpwd == ""{
		beego.Error("输入数据不完整")
		this.Data["errmsg"] = "输入数据不完整"
		this.TplName = "register-email.html"
		return
	}

	//两次密码是否一直
	if pwd != rpwd{
		beego.Error("两次密码输入不一致")
		this.Data["errmsg"] = "两次密码输入不一致"
		this.TplName = "register-email.html"
		return
	}
	//校验邮箱格式
	//把字符串全部大写
	reg ,_:=regexp.Compile(`^\w[\w\.-]*@[0-9a-z][0-9a-z-]*(\.[a-z]+)*\.[a-z]{2,6}$`)
	result := reg.FindString(email)
	if result == ""{
		beego.Error("邮箱格式错误")
		this.Data["errmsg"] = "邮箱格式错误"
		this.TplName = "register-email.html"
		return
	}

	//处理数据
	//发送邮件
	//utils     全局通用接口  工具类  邮箱配置   25,587
	config := `{"username":"czbkttsx@163.com","password":"czbkpygbj3q","host":"smtp.163.com","port":25}`
	emailReg :=utils.NewEMail(config)
	//内容配置
	emailReg.Subject = "品优购用户激活"
	emailReg.From = "czbkttsx@163.com"
	emailReg.To = []string{email}
	userName := this.Ctx.GetCookie("userName")
	emailReg.HTML = `<a href="http://192.168.176.88:8080/active?userName=`+userName+`"> 点击激活该用户</a>`

	//发送
	err := emailReg.Send()
	beego.Error(err)

	//插入邮箱   更新邮箱字段
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err = o.Read(&user,"Name")
	if err != nil {
		beego.Error("错误处理")
		return
	}
	user.Email = email
	o.Update(&user)




	//返回数据
	this.Ctx.WriteString("邮件已发送，请去目标邮箱激活用户！")
}

//激活
func(this*UserController)Active(){
	//获取数据
	userName := this.GetString("userName")
	//校验数据
	if userName == "" {
		beego.Error("用户名错误")
		this.Redirect("/register-email",302)
		return
	}

	//处理数据   本质上是更新active
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user,"Name")
	if err != nil {
		beego.Error("用户名不存在")
		this.Redirect("/register-email",302)
		return
	}
	user.Active = true
	o.Update(&user,"Active")

	//返回数据
	this.Redirect("/login",302)
}

func(this*UserController)ShowLogin(){
	name := this.Ctx.GetCookie("LoginName")
	if name == ""{
		//this.Data["name"] = name
		this.Data["checked"] = ""
	}else {
		//this.Data["name"] = name
		this.Data["checked"] = "checked"
	}
	//指定视图页面
	this.Data["name"] = name
	this.TplName = "login.html"
}

//处理登录业务
func(this*UserController)HandleLogin(){
	//获取数据   注册的时候要求用户名必须为字母加数字
	name := this.GetString("name")
	pwd := this.GetString("pwd")
	//校验数据
	if name == "" || pwd == ""{
		this.Data["errmsg"] = "获取数据错误"
		this.TplName = "login.html"
		return
	}
	//处理数据
	o := orm.NewOrm()
	var user models.User
	//赋值
	reg ,_:=regexp.Compile(`^\w[\w\.-]*@[0-9a-z][0-9a-z-]*(\.[a-z]+)*\.[a-z]{2,6}$`)
	result := reg.FindString(name)
	if result != ""{
		user.Email = name
		err := o.Read(&user,"Email")
		if err != nil {
			this.Data["errmsg"] = "邮箱未注册"
			this.TplName = "login.html"
			return
		}
		if user.Pwd != pwd {
			this.Data["errmsg"] = "密码错误"
			this.TplName = "login.html"
			return
		}

	}else {
		user.Name = name
		err := o.Read(&user,"Name")
		if err != nil{
			this.Data["errmsg"] = "用户名不存在"
			this.TplName = "login.html"
			return
		}
		if user.Pwd != pwd{
			this.Data["errmsg"] = "密码错误"
			this.TplName = "login.html"
			return
		}

	}

	//校验用户是否激活
	if user.Active == false{
		this.Data["errmsg"] = "当前用户未激活，请去目标邮箱激活！"
		this.TplName = "login.html"
		return
	}
	//返回数据u  cookie不能存中文  base64   序列化
	m1 := this.GetString("m1")
	if m1 == "2"{
		this.Ctx.SetCookie("LoginName",user.Name,60*60)
	}else{
		this.Ctx.SetCookie("LoginName",user.Name,-1)
	}
	//this.SetSession("name",user.Name)
	this.Redirect("/index",302)
}


