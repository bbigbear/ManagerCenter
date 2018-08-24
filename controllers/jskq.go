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

type JskqController struct {
	BaseController
}

func (this *JskqController) JskqAdd() {
	fmt.Println("添加教师考勤")
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
	var jskq models.Jskq
	json.Unmarshal(this.Ctx.Input.RequestBody, &jskq)

	fmt.Println("jskq:", &jskq)

	//time
	nowtime, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Println("time err", err.Error())
	}
	jskq.CreateDate = nowtime

	//insert
	_, err1 := o.Insert(&jskq)
	if err1 != nil {
		fmt.Printf("insert err", err1.Error())
		this.ajaxMsg("insert err", MSG_ERR_Resources)
	}
	list["id"] = jskq.Id
	this.ajaxList("add success", MSG_OK, 1, list)
}

func (this *JskqController) JskqGetData() {
	fmt.Println("获取教师考勤")
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
	//
	teachName := this.Input().Get("name")
	if teachName != "" {
		filters = append(filters, "TeacherName", teachName)
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
		fmt.Println("get jskq err", err.Error())
		this.ajaxMsg("get jskq err", MSG_ERR_Resources)
	}
	fmt.Println("get tax reslut num:", num)
	this.ajaxList("get tax data success", MSG_OK, count, maps)
}

func (this *JskqController) JskqEdit() {
	fmt.Println("编辑教师考勤")
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
	var jskq models.Jskq
	json.Unmarshal(this.Ctx.Input.RequestBody, &jskq)
	fmt.Println("jskq:", &jskq)
	//updata jskq db
	_, err1 := o.Update(&jskq)
	if err1 != nil {
		fmt.Println("updata jskq err", err1.Error())
		this.ajaxMsg("updata jskq err", MSG_ERR_Resources)
	}
	this.ajaxMsg("update jskq success", MSG_OK)
}

func (this *JskqController) TaxDel() {
	fmt.Println("删除教师考勤")
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

	var jskq models.Jskq
	json.Unmarshal(this.Ctx.Input.RequestBody, &jskq)
	fmt.Println("jskq:", &jskq)
	id := jskq.Id
	if id == 0 {
		this.ajaxMsg("删除失败", MSG_ERR_Param)
	}
	o := orm.NewOrm()
	c := new(models.Tax)
	num, err := o.QueryTable(c).Filter("Id", id).Delete()
	if err != nil {
		fmt.Println("del jskq err", err.Error())
		this.ajaxMsg("del jskq err", MSG_ERR_Resources)
	}
	fmt.Println("del jskq reslut num:", num)
	//list["data"] = maps
	this.ajaxMsg("del jskq success", MSG_OK)
}
