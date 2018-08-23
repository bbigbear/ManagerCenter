package routers

import (
	"ManagerCenter/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/v1/put_file", &controllers.BaseController{}, "*:PutFile")
	beego.Router("/v1/login", &controllers.LoginController{}, "post:LoginAction")
}
