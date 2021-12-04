package pinduoduo

import (
	"fmt"
	"regexp"
	"encoding/json"
	"time"
	"crypto/md5"
//	"encoding/hex"
//	"unicode/utf8"
	"strings"
	"strconv"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
//	"github.com/buger/jsonparser"
)


var pinduoduo = core.NewBucket("pinduoduo")
//拼多多
var apitype ="pdd.ddk.goods.zs.unit.url.gen"
var client_id=""
var client_key=""
var pid=""
var timestamp=""

type ItemUrl struct {
	GoodsZsUnitGenerateResponse struct {
		MultiGroupMobileShortURL string `json:"multi_group_mobile_short_url"`
		MultiGroupURL            string `json:"multi_group_url"`
		MobileURL                string `json:"mobile_url"`
		MultiGroupShortURL       string `json:"multi_group_short_url"`
		MobileShortURL           string `json:"mobile_short_url"`
		MultiGroupMobileURL      string `json:"multi_group_mobile_url"`
		RequestID                string `json:"request_id"`
		URL                      string `json:"url"`
		ShortURL                 string `json:"short_url"`
	} `json:"goods_zs_unit_generate_response"`
}

func init() {

	core.Server.GET("/pinduoduo/:sku", func(c *gin.Context) {
		sku := c.Param("sku")
		c.String(200, core.OttoFuncs["pinduoduo"](sku))
	})
	//添加命令
	core.AddCommand("", []core.Function{
		{
			Rules: []string{"raw https?://mobile\\.yangkeduo\\.com/goods.?\\.html\\?goods_id=(\\d+)"},
			Handle: func(s core.Sender) interface{} {
				fmt.Println(s.GetContent())			
				return getPinduoduo(s.GetContent())
			},
		},
	})
	core.OttoFuncs["pinduoduo"] = getPinduoduo //类似于向核心组件注册
}

func getPinduoduo(info string) string{
	/*将长链接变换成短链接*/
	//从返回的数据中提取出商品id
	var source_url=""
	reg := regexp.MustCompile(`https?://mobile\.yangkeduo\.com/goods.?\.html\?goods_id=(\d+)`)
	if reg != nil {
		params := reg.FindStringSubmatch(string(info))
		source_url = params[0]
		fmt.Println("链接:"+source_url+"\n")
	}
	return getShortUrl(source_url)
}

//获取短链接
func getShortUrl(source_url string) string {
	//对公共参数和业务参数按照ASCII排序,不参加排序：app_secret和sign
 	client_id=pinduoduo.Get("client_id")
 	pid=pinduoduo.Get("pid")
	client_key=pinduoduo.Get("client_key")
	//Unix时间timestamp
	timestamp= strconv.FormatInt(time.Now().Unix(),10)
	//大写(MD5(client_secret+key1+value1+key2+value2+client_secret))
	//将排序好的参数名和参数值拼装在一起，两头加client_key
	strCon:=client_key+
			"client_id"+client_id+
			"pid"+pid+
			"source_url"+source_url+
			"timestamp"+timestamp+
			"type"+apitype+
			client_key
	//md5
	strMd5:=md5.Sum([]byte(strCon))
	upMd5:=strings.ToUpper(fmt.Sprintf("%x",strMd5))
	//将长链接变换成短链接
	req := httplib.Get("https://gw-api.pinduoduo.com/api/router?"+
					"type="+apitype+
					"&client_id="+client_id+
					"&timestamp="+timestamp+
					"&sign="+upMd5+
					"&pid="+pid+
					"&source_url="+source_url)
	data, _:=req.Bytes()
	//fmt.Println(string(data))
	res:=&ItemUrl{}
	json.Unmarshal([]byte(data),&res)
	//fmt.Println(res.GoodsZsUnitGenerateResponse.ShortUrl)
	return string(res.GoodsZsUnitGenerateResponse.ShortURL)
}

// 创建一个错误处理函数，避免过多的 if err != nil{} 出现
func dropErr(e error) {
	if e != nil {
		panic(e)
	}
}