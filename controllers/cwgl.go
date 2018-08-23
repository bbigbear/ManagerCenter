package controllers

import (
	"ManagerCenter/models"
	"encoding/json"
	"fmt"
	"math"

	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/astaxie/beego/orm"
)

type CwglController struct {
	BaseController
}

func (this *CwglController) CzlxAdd() {
	fmt.Println("添加出账类型")
	//获取token
	var token models.Token
	json.Unmarshal(this.Ctx.Input.RequestBody, &token)

	if token.Token == "" {
		fmt.Println("token 为空")
		this.ajaxMsg("token is not nil", MSG_ERR_Param)
	}

	name, err := this.Token_auth(token.Token, "ximi")
	if err != nil {
		fmt.Println("token err", err)
		this.ajaxMsg("token err!", MSG_ERR_Verified)
	}
	fmt.Println("当前访问用户为:", name)

	o := orm.NewOrm()
	list := make(map[string]interface{})
	var czlx models.Czlx
	json.Unmarshal(this.Ctx.Input.RequestBody, &czlx)

	fmt.Println("czlx:", &czlx)

	//time
	//	nowtime, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	//	if err != nil {
	//		fmt.Println("time err", err.Error())
	//	}

	//insert
	_, err1 := o.Insert(&czlx)
	if err1 != nil {
		fmt.Printf("insert err", err1.Error())
		this.ajaxMsg("insert err", MSG_ERR_Resources)
	}
	list["id"] = czlx.Id
	this.ajaxList("add success", MSG_OK, 1, list)
}

func (this *CwglController) CzlxGetData() {
	fmt.Println("获取出账类型数据")
	//token
	token := this.Input().Get("token")

	if token == "" {
		fmt.Println("token 为空")
		this.ajaxMsg("token is not nil", MSG_ERR_Param)
	}

	name, err := this.Token_auth(token, "ximi")
	if err != nil {
		fmt.Println("token err", err.Error())
		this.ajaxMsg("token err!", MSG_ERR_Verified)
	}
	fmt.Println("当前访问用户为:", name)

	o := orm.NewOrm()
	var maps []orm.Params
	czlx := new(models.Czlx)
	query := o.QueryTable(czlx)
	filters := make([]interface{}, 0)
	//major
	czlx_name := this.Input().Get("name")
	if czlx_name != "" {
		filters = append(filters, "Name", czlx_name)
	}

	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}

	//index
	index, err := this.GetInt("index")
	if err != nil {
		fmt.Println("获取index下标错误")
	}

	//pagemax  一页多少
	pagemax, err := this.GetInt("pagemax")
	if err != nil {
		fmt.Println("获取每页数量为空")
	}

	//count
	count, err := query.Count()
	if err != nil {
		fmt.Println("获取数据总数为空")
		this.ajaxMsg("服务未知错误", MSG_ERR)
	}

	if pagemax != 0 {
		pagenum := int(math.Ceil(float64(count) / float64(pagemax)))

		if index > pagenum {
			//index = pagenum
			this.ajaxMsg("无法翻页了", MSG_ERR_Param)
		}
		fmt.Println("index&pagemax&pagenum", index, pagemax, pagenum)
	}
	query = query.Limit(pagemax, (index-1)*pagemax)

	//get data dB
	num, err := query.OrderBy("-Id").Values(&maps)
	if err != nil {
		fmt.Println("get czlx err", err.Error())
		this.ajaxMsg("get czlx err", MSG_ERR_Resources)
	}
	fmt.Println("get czlx reslut num:", num)
	this.ajaxList("get czlx data success", MSG_OK, count, maps)
}

func (this *CwglController) CzlxEdit() {
	fmt.Println("编辑出账类型")
	//获取token
	var token models.Token
	json.Unmarshal(this.Ctx.Input.RequestBody, &token)

	if token.Token == "" {
		fmt.Println("token 为空")
		this.ajaxMsg("token is not nil", MSG_ERR_Param)
	}

	name, err := this.Token_auth(token.Token, "ximi")
	if err != nil {
		fmt.Println("token err", err)
		this.ajaxMsg("token err!", MSG_ERR_Verified)
	}
	fmt.Println("当前访问用户为:", name)

	o := orm.NewOrm()
	var czlx models.Czlx
	json.Unmarshal(this.Ctx.Input.RequestBody, &czlx)
	fmt.Println("czlx:", &czlx)
	//updata czlx db
	_, err1 := o.Update(&czlx)
	if err1 != nil {
		fmt.Println("updata czlx err", err1.Error())
		this.ajaxMsg("updata czlx err", MSG_ERR_Resources)
	}
	this.ajaxMsg("update czxl success", MSG_OK)
}

func (this *CwglController) CzlxDel() {
	fmt.Println("删除出账类型")
	//获取token
	var token models.Token
	json.Unmarshal(this.Ctx.Input.RequestBody, &token)

	if token.Token == "" {
		fmt.Println("token 为空")
		this.ajaxMsg("token is not nil", MSG_ERR_Param)
	}

	name, err := this.Token_auth(token.Token, "ximi")
	if err != nil {
		fmt.Println("token err", err)
		this.ajaxMsg("token err!", MSG_ERR_Verified)
	}
	fmt.Println("当前访问用户为:", name)

	var czlx models.Czlx
	json.Unmarshal(this.Ctx.Input.RequestBody, &czlx)
	fmt.Println("czlx:", &czlx)
	id := czlx.Id
	if id == 0 {
		this.ajaxMsg("删除失败", MSG_ERR_Param)
	}
	o := orm.NewOrm()
	c := new(models.Czlx)
	num, err := o.QueryTable(c).Filter("Id", id).Delete()
	if err != nil {
		fmt.Println("del czlx err", err.Error())
		this.ajaxMsg("del czlx err", MSG_ERR_Resources)
	}
	fmt.Println("del czlx reslut num:", num)
	//list["data"] = maps
	this.ajaxMsg("del czlx success", MSG_OK)
}

//进账登记
func (this *CwglController) JzdjAdd() {
	fmt.Println("添加进账登记")
	//获取token
	var token models.Token
	json.Unmarshal(this.Ctx.Input.RequestBody, &token)

	if token.Token == "" {
		fmt.Println("token 为空")
		this.ajaxMsg("token is not nil", MSG_ERR_Param)
	}

	name, err := this.Token_auth(token.Token, "ximi")
	if err != nil {
		fmt.Println("token err", err)
		this.ajaxMsg("token err!", MSG_ERR_Verified)
	}
	fmt.Println("当前访问用户为:", name)

	o := orm.NewOrm()
	list := make(map[string]interface{})
	var jzdj models.Jzdj
	json.Unmarshal(this.Ctx.Input.RequestBody, &jzdj)

	fmt.Println("jzdj:", &jzdj)

	//time
	//	nowtime, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	//	if err != nil {
	//		fmt.Println("time err", err.Error())
	//	}

	//insert
	_, err1 := o.Insert(&jzdj)
	if err1 != nil {
		fmt.Printf("insert err", err1.Error())
		this.ajaxMsg("insert err", MSG_ERR_Resources)
	}
	list["id"] = jzdj.Id
	list["number"] = jzdj.Number
	this.ajaxList("add success", MSG_OK, 1, list)
}

func (this *CwglController) JzdjGetData() {
	fmt.Println("获取进账登记数据")
	//token
	token := this.Input().Get("token")

	if token == "" {
		fmt.Println("token 为空")
		this.ajaxMsg("token is not nil", MSG_ERR_Param)
	}

	name, err := this.Token_auth(token, "ximi")
	if err != nil {
		fmt.Println("token err", err.Error())
		this.ajaxMsg("token err!", MSG_ERR_Verified)
	}
	fmt.Println("当前访问用户为:", name)

	o := orm.NewOrm()
	var maps []orm.Params
	jzdj := new(models.Jzdj)
	query := o.QueryTable(jzdj)
	filters := make([]interface{}, 0)
	//name
	jzdj_name := this.Input().Get("name")
	if jzdj_name != "" {
		filters = append(filters, "Name", jzdj_name)
	}

	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}

	//index
	index, err := this.GetInt("index")
	if err != nil {
		fmt.Println("获取index下标错误")
	}

	//pagemax  一页多少
	pagemax, err := this.GetInt("pagemax")
	if err != nil {
		fmt.Println("获取每页数量为空")
	}

	//count
	count, err := query.Count()
	if err != nil {
		fmt.Println("获取数据总数为空")
		this.ajaxMsg("服务未知错误", MSG_ERR)
	}

	if pagemax != 0 {
		pagenum := int(math.Ceil(float64(count) / float64(pagemax)))

		if index > pagenum {
			//index = pagenum
			this.ajaxMsg("无法翻页了", MSG_ERR_Param)
		}
		fmt.Println("index&pagemax&pagenum", index, pagemax, pagenum)
	}
	query = query.Limit(pagemax, (index-1)*pagemax)

	//get data dB
	num, err := query.OrderBy("-Id").Values(&maps)
	if err != nil {
		fmt.Println("get jzdj err", err.Error())
		this.ajaxMsg("get jzdj err", MSG_ERR_Resources)
	}
	fmt.Println("get jzdj reslut num:", num)
	this.ajaxList("get czjzdjlx data success", MSG_OK, count, maps)
}

func (this *CwglController) JzdjEdit() {
	fmt.Println("编辑进账登记")
	//获取token
	var token models.Token
	json.Unmarshal(this.Ctx.Input.RequestBody, &token)

	if token.Token == "" {
		fmt.Println("token 为空")
		this.ajaxMsg("token is not nil", MSG_ERR_Param)
	}

	name, err := this.Token_auth(token.Token, "ximi")
	if err != nil {
		fmt.Println("token err", err)
		this.ajaxMsg("token err!", MSG_ERR_Verified)
	}
	fmt.Println("当前访问用户为:", name)

	o := orm.NewOrm()
	var jzdj models.Jzdj
	json.Unmarshal(this.Ctx.Input.RequestBody, &jzdj)
	fmt.Println("jzdj:", &jzdj)
	//updata jzdj db
	_, err1 := o.Update(&jzdj)
	if err1 != nil {
		fmt.Println("updata jzdj err", err1.Error())
		this.ajaxMsg("updata jzdj err", MSG_ERR_Resources)
	}
	this.ajaxMsg("update jzdj success", MSG_OK)
}

func (this *CwglController) JzdjDel() {
	fmt.Println("删除进账登记")
	//获取token
	var token models.Token
	json.Unmarshal(this.Ctx.Input.RequestBody, &token)

	if token.Token == "" {
		fmt.Println("token 为空")
		this.ajaxMsg("token is not nil", MSG_ERR_Param)
	}

	name, err := this.Token_auth(token.Token, "ximi")
	if err != nil {
		fmt.Println("token err", err)
		this.ajaxMsg("token err!", MSG_ERR_Verified)
	}
	fmt.Println("当前访问用户为:", name)

	var jzdj models.Jzdj
	json.Unmarshal(this.Ctx.Input.RequestBody, &jzdj)
	fmt.Println("jzdj:", &jzdj)
	id := jzdj.Id
	if id == 0 {
		this.ajaxMsg("删除失败", MSG_ERR_Param)
	}
	o := orm.NewOrm()
	c := new(models.Jzdj)
	num, err := o.QueryTable(c).Filter("Id", id).Delete()
	if err != nil {
		fmt.Println("del jzdj err", err.Error())
		this.ajaxMsg("del jzdj err", MSG_ERR_Resources)
	}
	fmt.Println("del jzdj reslut num:", num)
	//list["data"] = maps
	this.ajaxMsg("del jzdj success", MSG_OK)
}

//出账登记
func (this *CwglController) CzdjAdd() {
	fmt.Println("添加出账登记")
	//获取token
	var token models.Token
	json.Unmarshal(this.Ctx.Input.RequestBody, &token)

	if token.Token == "" {
		fmt.Println("token 为空")
		this.ajaxMsg("token is not nil", MSG_ERR_Param)
	}

	name, err := this.Token_auth(token.Token, "ximi")
	if err != nil {
		fmt.Println("token err", err)
		this.ajaxMsg("token err!", MSG_ERR_Verified)
	}
	fmt.Println("当前访问用户为:", name)

	o := orm.NewOrm()
	list := make(map[string]interface{})
	var czdj models.Czdj
	json.Unmarshal(this.Ctx.Input.RequestBody, &czdj)

	fmt.Println("czdj:", &czdj)
	//time
	//	nowtime, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	//	if err != nil {
	//		fmt.Println("time err", err.Error())
	//	}

	//insert
	_, err1 := o.Insert(&czdj)
	if err1 != nil {
		fmt.Printf("insert err", err1.Error())
		this.ajaxMsg("insert err", MSG_ERR_Resources)
	}
	list["id"] = czdj.Id
	list["number"] = czdj.Number
	this.ajaxList("add success", MSG_OK, 1, list)
}

func (this *CwglController) CzdjGetData() {
	fmt.Println("获取进账登记数据")
	//token
	token := this.Input().Get("token")

	if token == "" {
		fmt.Println("token 为空")
		this.ajaxMsg("token is not nil", MSG_ERR_Param)
	}

	name, err := this.Token_auth(token, "ximi")
	if err != nil {
		fmt.Println("token err", err.Error())
		this.ajaxMsg("token err!", MSG_ERR_Verified)
	}
	fmt.Println("当前访问用户为:", name)

	o := orm.NewOrm()
	var maps []orm.Params
	czdj := new(models.Czdj)
	query := o.QueryTable(czdj)
	filters := make([]interface{}, 0)
	//name
	czdj_name := this.Input().Get("name")
	if czdj_name != "" {
		filters = append(filters, "Name", czdj_name)
	}

	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}

	//index
	index, err := this.GetInt("index")
	if err != nil {
		fmt.Println("获取index下标错误")
	}

	//pagemax  一页多少
	pagemax, err := this.GetInt("pagemax")
	if err != nil {
		fmt.Println("获取每页数量为空")
	}

	//count
	count, err := query.Count()
	if err != nil {
		fmt.Println("获取数据总数为空")
		this.ajaxMsg("服务未知错误", MSG_ERR)
	}

	if pagemax != 0 {
		pagenum := int(math.Ceil(float64(count) / float64(pagemax)))

		if index > pagenum {
			//index = pagenum
			this.ajaxMsg("无法翻页了", MSG_ERR_Param)
		}
		fmt.Println("index&pagemax&pagenum", index, pagemax, pagenum)
	}
	query = query.Limit(pagemax, (index-1)*pagemax)

	//get data dB
	num, err := query.OrderBy("-Id").Values(&maps)
	if err != nil {
		fmt.Println("get jzdj err", err.Error())
		this.ajaxMsg("get jzdj err", MSG_ERR_Resources)
	}
	fmt.Println("get jzdj reslut num:", num)
	this.ajaxList("get czjzdjlx data success", MSG_OK, count, maps)
}

func (this *CwglController) CzdjEdit() {
	fmt.Println("编辑进账登记")
	//获取token
	var token models.Token
	json.Unmarshal(this.Ctx.Input.RequestBody, &token)

	if token.Token == "" {
		fmt.Println("token 为空")
		this.ajaxMsg("token is not nil", MSG_ERR_Param)
	}

	name, err := this.Token_auth(token.Token, "ximi")
	if err != nil {
		fmt.Println("token err", err)
		this.ajaxMsg("token err!", MSG_ERR_Verified)
	}
	fmt.Println("当前访问用户为:", name)

	o := orm.NewOrm()
	var czdj models.Czdj
	json.Unmarshal(this.Ctx.Input.RequestBody, &czdj)
	fmt.Println("czdj:", &czdj)
	//updata czdj db
	_, err1 := o.Update(&czdj)
	if err1 != nil {
		fmt.Println("updata czdj err", err1.Error())
		this.ajaxMsg("updata czdj err", MSG_ERR_Resources)
	}
	this.ajaxMsg("update czdj success", MSG_OK)
}

func (this *CwglController) CzdjDel() {
	fmt.Println("删除进账登记")
	//获取token
	var token models.Token
	json.Unmarshal(this.Ctx.Input.RequestBody, &token)

	if token.Token == "" {
		fmt.Println("token 为空")
		this.ajaxMsg("token is not nil", MSG_ERR_Param)
	}

	name, err := this.Token_auth(token.Token, "ximi")
	if err != nil {
		fmt.Println("token err", err)
		this.ajaxMsg("token err!", MSG_ERR_Verified)
	}
	fmt.Println("当前访问用户为:", name)

	var czdj models.Czdj
	json.Unmarshal(this.Ctx.Input.RequestBody, &czdj)
	fmt.Println("czdj:", &czdj)
	id := czdj.Id
	if id == 0 {
		this.ajaxMsg("删除失败", MSG_ERR_Param)
	}
	o := orm.NewOrm()
	c := new(models.Czdj)
	num, err := o.QueryTable(c).Filter("Id", id).Delete()
	if err != nil {
		fmt.Println("del czdj err", err.Error())
		this.ajaxMsg("del czdj err", MSG_ERR_Resources)
	}
	fmt.Println("del czdj reslut num:", num)
	//list["data"] = maps
	this.ajaxMsg("del czdj success", MSG_OK)
}
