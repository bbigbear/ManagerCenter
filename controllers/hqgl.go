package controllers

import (
	"ManagerCenter/models"
	"encoding/json"
	"fmt"
	"math"
	"time"

	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/astaxie/beego/orm"
)

type HqglController struct {
	BaseController
}

func (this *HqglController) QjqAdd() {
	fmt.Println("添加清洁区")
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
	var qjq models.Qjq
	json.Unmarshal(this.Ctx.Input.RequestBody, &qjq)

	fmt.Println("qjq:", &qjq)

	//time
	//	nowtime, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	//	if err != nil {
	//		fmt.Println("time err", err.Error())
	//	}

	//insert
	_, err1 := o.Insert(&qjq)
	if err1 != nil {
		fmt.Printf("insert err", err1.Error())
		this.ajaxMsg("insert err", MSG_ERR_Resources)
	}
	list["id"] = qjq.Id
	this.ajaxList("add success", MSG_OK, 1, list)
}

func (this *HqglController) QjqGetData() {
	fmt.Println("获取清洁区")
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
	qjq := new(models.Qjq)
	query := o.QueryTable(qjq)
	filters := make([]interface{}, 0)
	//major
	qjq_name := this.Input().Get("name")
	if qjq_name != "" {
		filters = append(filters, "Name", qjq_name)
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
		fmt.Println("get qjq err", err.Error())
		this.ajaxMsg("get qjq err", MSG_ERR_Resources)
	}
	fmt.Println("get qjq reslut num:", num)
	this.ajaxList("get qjq data success", MSG_OK, count, maps)
}

func (this *HqglController) QjqEdit() {
	fmt.Println("编辑清洁区")
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
	var qjq models.Qjq
	json.Unmarshal(this.Ctx.Input.RequestBody, &qjq)
	fmt.Println("qjq:", &qjq)
	//updata qjq db
	_, err1 := o.Update(&qjq)
	if err1 != nil {
		fmt.Println("updata qjq err", err1.Error())
		this.ajaxMsg("updata qjq err", MSG_ERR_Resources)
	}
	this.ajaxMsg("update czxl success", MSG_OK)
}

func (this *HqglController) QjqDel() {
	fmt.Println("删除清洁区")
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

	var qjq models.Qjq
	json.Unmarshal(this.Ctx.Input.RequestBody, &qjq)
	fmt.Println("qjq:", &qjq)
	id := qjq.Id
	if id == 0 {
		this.ajaxMsg("删除失败", MSG_ERR_Param)
	}
	o := orm.NewOrm()
	c := new(models.Qjq)
	num, err := o.QueryTable(c).Filter("Id", id).Delete()
	if err != nil {
		fmt.Println("del qjq err", err.Error())
		this.ajaxMsg("del qjq err", MSG_ERR_Resources)
	}
	fmt.Println("del qjq reslut num:", num)
	//list["data"] = maps
	this.ajaxMsg("del qjq success", MSG_OK)
}

//班级
func (this *HqglController) ClassAdd() {
	fmt.Println("添加班级")
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
	var class models.Class
	json.Unmarshal(this.Ctx.Input.RequestBody, &class)
	fmt.Println("class:", &class)

	//insert
	_, err1 := o.Insert(&class)
	if err1 != nil {
		fmt.Printf("insert err", err1.Error())
		this.ajaxMsg("insert err", MSG_ERR_Resources)
	}
	list["id"] = class.Id
	this.ajaxList("add success", MSG_OK, 1, list)
}

func (this *HqglController) ClassGetData() {
	fmt.Println("获取班级数据")
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
	class := new(models.Class)
	query := o.QueryTable(class)
	filters := make([]interface{}, 0)
	//major
	class_name := this.Input().Get("name")
	if class_name != "" {
		filters = append(filters, "Name", class_name)
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
		fmt.Println("get qjq err", err.Error())
		this.ajaxMsg("get qjq err", MSG_ERR_Resources)
	}
	fmt.Println("get qjq reslut num:", num)
	this.ajaxList("get qjq data success", MSG_OK, count, maps)
}

func (this *HqglController) ClassDel() {
	fmt.Println("删除班级")
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

	var class models.Class
	json.Unmarshal(this.Ctx.Input.RequestBody, &class)
	fmt.Println("class:", &class)
	id := class.Id
	if id == 0 {
		this.ajaxMsg("删除失败", MSG_ERR_Param)
	}
	o := orm.NewOrm()
	c := new(models.Class)
	num, err := o.QueryTable(c).Filter("Id", id).Delete()
	if err != nil {
		fmt.Println("del class err", err.Error())
		this.ajaxMsg("del class err", MSG_ERR_Resources)
	}
	fmt.Println("del class reslut num:", num)
	//list["data"] = maps
	this.ajaxMsg("del class success", MSG_OK)
}

//绿化
func (this *HqglController) TreeAdd() {
	fmt.Println("添加绿化")
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
	var tree models.Tree
	json.Unmarshal(this.Ctx.Input.RequestBody, &tree)

	fmt.Println("tree:", &tree)

	//time
	nowtime, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Println("time err", err.Error())
	}

	tree.Date = nowtime
	//insert
	_, err1 := o.Insert(&tree)
	if err1 != nil {
		fmt.Printf("insert err", err1.Error())
		this.ajaxMsg("insert err", MSG_ERR_Resources)
	}
	list["id"] = tree.Id
	this.ajaxList("add success", MSG_OK, 1, list)
}

func (this *HqglController) TreeGetData() {
	fmt.Println("获取绿化数据")
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
	tree := new(models.Tree)
	query := o.QueryTable(tree)
	filters := make([]interface{}, 0)
	//major
	tree_name := this.Input().Get("name")
	if tree_name != "" {
		filters = append(filters, "Name", tree_name)
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
		fmt.Println("get tree err", err.Error())
		this.ajaxMsg("get tree err", MSG_ERR_Resources)
	}
	fmt.Println("get tree reslut num:", num)
	this.ajaxList("get tree data success", MSG_OK, count, maps)
}

func (this *HqglController) TreeEdit() {
	fmt.Println("编辑绿化")
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
	var tree models.Tree
	json.Unmarshal(this.Ctx.Input.RequestBody, &tree)
	fmt.Println("tree:", &tree)
	//updata tree db
	_, err1 := o.Update(&tree)
	if err1 != nil {
		fmt.Println("updata tree err", err1.Error())
		this.ajaxMsg("updata tree err", MSG_ERR_Resources)
	}
	this.ajaxMsg("update tree success", MSG_OK)
}

func (this *HqglController) TreeDel() {
	fmt.Println("删除绿化")
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

	var tree models.Tree
	json.Unmarshal(this.Ctx.Input.RequestBody, &tree)
	fmt.Println("tree:", &tree)
	id := tree.Id
	if id == 0 {
		this.ajaxMsg("删除失败", MSG_ERR_Param)
	}
	o := orm.NewOrm()
	c := new(models.Tree)
	num, err := o.QueryTable(c).Filter("Id", id).Delete()
	if err != nil {
		fmt.Println("del tree err", err.Error())
		this.ajaxMsg("del tree err", MSG_ERR_Resources)
	}
	fmt.Println("del tree reslut num:", num)
	//list["data"] = maps
	this.ajaxMsg("del tree success", MSG_OK)
}
