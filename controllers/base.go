package controllers

import (
	"ManagerCenter/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/Go-SQL-Driver/MySQL"

	//"crypto/rand"
	"math/rand"

	"github.com/LindsayBradford/go-dbf/godbf"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/dgrijalva/jwt-go"
)

const (
	MSG_OK            = 200
	MSG_ERR_Param     = 400
	MSG_ERR_Verified  = 401
	MSG_ERR_Authority = 403
	MSG_ERR_Resources = 404
	MSG_ERR           = 500
)

type BaseController struct {
	beego.Controller
}

//ajax返回
func (this *BaseController) ajaxMsg(msg interface{}, msgno int) {
	out := make(map[string]interface{})
	out["code"] = msgno
	out["message"] = msg
	this.Data["json"] = out
	this.ServeJSON()
	this.StopRun()
}

//ajax返回 列表
func (this *BaseController) ajaxList(msg interface{}, msgno int, count int64, data interface{}) {
	out := make(map[string]interface{})
	out["code"] = msgno
	out["message"] = msg
	out["count"] = count
	out["data"] = data
	this.Data["json"] = out
	this.ServeJSON()
	this.StopRun()
}

// 通过两重循环过滤重复元素
func (this *BaseController) RemoveRepBySlice(slc []string) []string {
	result := []string{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

// 图片接口
func (this *BaseController) PutFile() {
	h, err := this.GetFiles("file")
	fmt.Println("文件名称", h[0].Filename)
	fmt.Println("文件大小", h[0].Size)
	if err != nil {
		log.Fatal("getfile err ", err)
		this.ajaxMsg(h[0].Filename+"文件上传失败", MSG_ERR_Resources)
	}
	//	defer f.Close()
	path := "static/upload/" + h[0].Filename
	this.SaveToFile("file", path) // 保存位置在 static/upload, 没有文件夹要先创建
	list := make(map[string]interface{})
	list["src"] = path
	list["name"] = h[0].Filename
	list["size"] = h[0].Size
	this.ajaxList("文件上传成功", MSG_OK, 1, list)
}

// 图片接口
func (this *BaseController) PutDbf() {
	h, err := this.GetFiles("file")
	fmt.Println("文件名称", h[0].Filename)
	fmt.Println("文件大小", h[0].Size)
	if err != nil {
		log.Fatal("getfile err ", err)
		this.ajaxMsg(h[0].Filename+"文件上传失败", MSG_ERR_Resources)
	}
	//	defer f.Close()
	path := "static/upload/" + h[0].Filename
	this.SaveToFile("file", path) // 保存位置在 static/upload, 没有文件夹要先创建
	dbfTable, err := godbf.NewFromFile(path, "UTF8")
	if err != nil {
		fmt.Println("dbf err", err.Error())
	}
	num, err := dbfTable.DecimalPlacesInField("KSH")
	if err == nil {
		fmt.Println("decimalplaces", num)
	}
	list := make(map[string]interface{})
	list["name"] = dbfTable.FieldNames()
	list["num"] = dbfTable.NumberOfRecords()
	this.ajaxList("文件上传成功", MSG_OK, 1, list)
}

//将时间化为秒
func (this *BaseController) GetSecs(ordertime string) int64 {
	var s int64
	t, err := time.ParseInLocation("2006-01-02 15:04:05", ordertime, time.Local)
	if err == nil {
		s = t.Unix()
		return s
	} else {
		return -1
	}
}

//获取相差时间
func (this *BaseController) GetMinuteDiffer(server_time, mqtime string) int64 {
	var minute int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", server_time, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", mqtime, time.Local)
	if err == nil {
		diff := t1.Unix() - t2.Unix()
		minute = diff / 60
		return minute
	} else {
		return -1
	}
}

//生成随机数
func (this *BaseController) randStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

//随机字符
func (this *BaseController) GetRandomString(l int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

type Claims struct {
	Appid string `json:"appid"`
	// recommended having
	jwt.StandardClaims
}

//jwt
func (this *BaseController) Create_token(appid string, secret string) (string, int64) {
	expireToken := time.Now().Add(time.Hour * 24).Unix()
	claims := Claims{
		appid,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    appid,
		},
	}

	// Create the token using your claims
	c_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signs the token with a secret.
	signedToken, _ := c_token.SignedString([]byte(secret))

	return signedToken, expireToken
}

func (this *BaseController) Token_auth(signedToken, secret string) (string, error) {
	token, err := jwt.ParseWithClaims(signedToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		//fmt.Printf("%v %v", claims.Username, claims.StandardClaims.ExpiresAt)
		//fmt.Println(reflect.TypeOf(claims.StandardClaims.ExpiresAt))
		//return claims.Appid, err
		return claims.Appid, err
	}
	return "", err
}

//单点登录
func (this *BaseController) SessionLogin(skey string) int {
	buf, err := base64.StdEncoding.DecodeString(skey)
	if err != nil {
		return 0
	}

	var v []string
	err1 := json.Unmarshal(buf, &v)
	if err1 != nil {
		return 0
	}

	url := beego.AppConfig.String("session_login_url")
	reqParm := httplib.Get(url)
	cookie := &http.Cookie{
		Name:  v[0],
		Value: v[1],
	}
	reqParm.SetCookie(cookie)
	//	str, err := reqParm.String()
	//	if err != nil {
	//		fmt.Println("err", err.Error())
	//	}
	//	fmt.Println(str)
	var login_info models.Login
	err2 := reqParm.ToJSON(&login_info)
	if err2 != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Println("login_info", &login_info)
	if login_info.ReadName == "管理员" {
		fmt.Println("登录成功")
		//存session
		this.SetSession("islogin", 1)
		return 1
	} else {
		fmt.Println("不是管理员，无权登录")
		return 0
	}
}
