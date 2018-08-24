package controllers

import (
	"ManagerCenter/models"
	"encoding/json"
	"fmt"
	"math"
	//	"time"

	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/astaxie/beego/orm"
)

type GsglController struct {
	BaseController
}

func (this *GsglController) TaxAdd() {
	fmt.Println("添加个税项目")
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
	var tax models.Tax
	json.Unmarshal(this.Ctx.Input.RequestBody, &tax)

	fmt.Println("tax:", &tax)

	//time
	//	nowtime, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	//	if err != nil {
	//		fmt.Println("time err", err.Error())
	//	}

	//insert
	_, err1 := o.Insert(&tax)
	if err1 != nil {
		fmt.Printf("insert err", err1.Error())
		this.ajaxMsg("insert err", MSG_ERR_Resources)
	}
	list["id"] = tax.Id
	this.ajaxList("add success", MSG_OK, 1, list)
}

func (this *GsglController) TaxGetData() {
	fmt.Println("获取个税项目数据")
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
	tax := new(models.Tax)
	query := o.QueryTable(tax)
	filters := make([]interface{}, 0)
	//major
	tax_name := this.Input().Get("name")
	if tax_name != "" {
		filters = append(filters, "Name", tax_name)
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
		fmt.Println("get tax err", err.Error())
		this.ajaxMsg("get tax err", MSG_ERR_Resources)
	}
	fmt.Println("get tax reslut num:", num)
	this.ajaxList("get tax data success", MSG_OK, count, maps)
}

func (this *GsglController) TaxEdit() {
	fmt.Println("编辑个税项目")
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
	var tax models.Tax
	json.Unmarshal(this.Ctx.Input.RequestBody, &tax)
	fmt.Println("tax:", &tax)
	//updata tax db
	_, err1 := o.Update(&tax)
	if err1 != nil {
		fmt.Println("updata tax err", err1.Error())
		this.ajaxMsg("updata tax err", MSG_ERR_Resources)
	}
	this.ajaxMsg("update tax success", MSG_OK)
}

func (this *GsglController) TaxDel() {
	fmt.Println("删除个税项目")
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

	var tax models.Tax
	json.Unmarshal(this.Ctx.Input.RequestBody, &tax)
	fmt.Println("tax:", &tax)
	id := tax.Id
	if id == 0 {
		this.ajaxMsg("删除失败", MSG_ERR_Param)
	}
	o := orm.NewOrm()
	c := new(models.Tax)
	num, err := o.QueryTable(c).Filter("Id", id).Delete()
	if err != nil {
		fmt.Println("del tax err", err.Error())
		this.ajaxMsg("del tax err", MSG_ERR_Resources)
	}
	fmt.Println("del tax reslut num:", num)
	//list["data"] = maps
	this.ajaxMsg("del tax success", MSG_OK)
}

//个税月份管理

func (this *GsglController) TaxMonthAdd() {
	fmt.Println("添加个税项目")
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
	var tax models.TaxMonth
	json.Unmarshal(this.Ctx.Input.RequestBody, &tax)

	fmt.Println("tax:", &tax)

	//time
	//	nowtime, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	//	if err != nil {
	//		fmt.Println("time err", err.Error())
	//	}

	//insert
	_, err1 := o.Insert(&tax)
	if err1 != nil {
		fmt.Printf("insert err", err1.Error())
		this.ajaxMsg("insert err", MSG_ERR_Resources)
	}
	list["id"] = tax.Id
	this.ajaxList("add success", MSG_OK, 1, list)
}

func (this *GsglController) TaxMonthGetData() {
	fmt.Println("获取个税项目数据")
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
	tax := new(models.TaxMonth)
	query := o.QueryTable(tax)
	filters := make([]interface{}, 0)
	//major
	tax_name := this.Input().Get("name")
	if tax_name != "" {
		filters = append(filters, "Name", tax_name)
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
	num, err := query.OrderBy("Sort").Values(&maps)
	if err != nil {
		fmt.Println("get taxMonth err", err.Error())
		this.ajaxMsg("get taxMonth err", MSG_ERR_Resources)
	}
	fmt.Println("get taxMonth reslut num:", num)
	this.ajaxList("get taxMonth data success", MSG_OK, count, maps)
}

func (this *GsglController) TaxMonthEdit() {
	fmt.Println("编辑个税项目")
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
	var tax models.TaxMonth
	json.Unmarshal(this.Ctx.Input.RequestBody, &tax)
	fmt.Println("tax:", &tax)
	//updata tax db
	_, err1 := o.Update(&tax)
	if err1 != nil {
		fmt.Println("updata taxMonth err", err1.Error())
		this.ajaxMsg("updata taxMonth err", MSG_ERR_Resources)
	}
	this.ajaxMsg("update taxMonth success", MSG_OK)
}

func (this *GsglController) TaxMonthDel() {
	fmt.Println("删除个税项目")
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

	var tax models.TaxMonth
	json.Unmarshal(this.Ctx.Input.RequestBody, &tax)
	fmt.Println("tax:", &tax)
	id := tax.Id
	if id == 0 {
		this.ajaxMsg("删除失败", MSG_ERR_Param)
	}
	o := orm.NewOrm()
	c := new(models.Tax)
	num, err := o.QueryTable(c).Filter("Id", id).Delete()
	if err != nil {
		fmt.Println("del taxMonth err", err.Error())
		this.ajaxMsg("del taxMonth err", MSG_ERR_Resources)
	}
	fmt.Println("del taxMonth reslut num:", num)
	//list["data"] = maps
	this.ajaxMsg("del taxMonth success", MSG_OK)
}
