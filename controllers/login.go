package controllers

import (
	"encoding/json"
	"fmt"

	"ClassCenter/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/tidwall/gjson"
)

type LoginController struct {
	BaseController
}

func (this *LoginController) LoginAction() {

	fmt.Println("点击登录按钮")
	var user models.User
	json.Unmarshal(this.Ctx.Input.RequestBody, &user)
	uname := user.Name
	pwd := user.Pwd
	fmt.Println("get name&password", uname, pwd)

	//	if beego.AppConfig.String("uname") == uname &&
	//		beego.AppConfig.String("pwd") == pwd {
	//存session
	//this.SetSession("islogin", 1)
	token := beego.AppConfig.String("token")
	url := beego.AppConfig.String("login_url")
	req := httplib.Post(url)
	req.Header("Content-Type", "application/json")
	req.Header("token", token)
	req.Header("system", "perm")

	login := make(map[string]interface{})
	login["username"] = uname
	login["password"] = pwd
	req.JSONBody(login)

	str, err := req.String()
	if err != nil {
		fmt.Println("err", err.Error())
	}
	fmt.Println("str", str)
	code := gjson.Get(str, "code")
	if code.Exists() {
		this.ajaxMsg("账号密码错误", MSG_ERR_Param)
	} else {
		var login_info models.Login
		err := req.ToJSON(&login_info)
		if err != nil {
			fmt.Println(err)
			this.ajaxMsg("err to json", MSG_ERR)
		}
		if login_info.ReadName == "管理员" {
			//返回jwt
			jwt, i := this.Create_token(user.Name, "ximi")
			fmt.Println("token&time", jwt, i)
			list := make(map[string]interface{})
			list["name"] = login_info.ReadName
			list["token"] = jwt
			this.ajaxList("登录成功", MSG_OK, 1, list)
		} else {
			this.ajaxMsg("你无权登录", MSG_ERR_Verified)
		}
	}

	//	} else {
	//		fmt.Println("账户密码错误")
	//		this.ajaxMsg("账户密码错误", MSG_ERR)
	//	}
}
