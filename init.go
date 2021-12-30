/*******************************************
QQ交流群：418353744
QQ线报群：263723430
********************************************/
package jingdong

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/buger/jsonparser"
	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
)

var jingdong = core.NewBucket("jingdong")

func init() {

	core.Server.GET("/jingdong/:sku", func(c *gin.Context) {
		sku := c.Param("sku")
		c.String(200, core.OttoFuncs["jingdong"](sku))
	})

	//向core包中添加命令
	core.AddCommand("", []core.Function{
		{//"https:\\/\\/item.m.jd.com\\/product\\/100022622786.html
			Rules: []string{"raw https?://item\\.m\\.jd\\.[comhk]{2,3}/product/(\\d+).html",
				"raw https?:\\\\\\/\\\\\\/item\\.m\\.jd\\.[comhk]{2,3}\\\\\\/product\\\\\\/(\\d+).html",
				"raw https?://.+\\.jd\\.[comhk]{2,3}/(\\d+).html",
				"raw https?:\\\\\\/\\\\\\/.+\\.jd\\.[comhk]{2,3}\\\\\\/(\\d+).html",
				"raw https?://item\\.m\\.jd\\.[comhk]{2,3}/(\\d+).html",
				"raw https?:\\\\\\/\\\\\\/item\\.m\\.jd\\.[comhk]{2,3}\\\\\\/(\\d+).html",
				"raw https?://m\\.jingxi\\.[comhk]{2,3}/item/jxview\\?sku=(\\d+)",
				"raw https?:\\\\\\/\\\\\\/m\\.jingxi\\.[comhk]{2,3}\\\\\\/item\\\\\\/jxview\\?sku=(\\d+)",
				"raw https?://m\\.jingxi\\.[comhk]{2,3}/item/view\\?sku=(\\d+)",
				"raw https?:\\\\\\/\\\\\\/m\\.jingxi\\.[comhk]{2,3}\\\\\\/item\\\\\\/view\\?sku=(\\d+)",
				"raw https?://m\\.jingxi\\.[comhk]{2,3}.+sku=(\\d+)",
				"raw https?:\\\\\\/\\\\\\/m\\.jingxi\\.[comhk]{2,3}.+sku=(\\d+)",
				"raw https?://kpl\\.m\\.jd\\.[comhk]{2,3}/product\\?wareId=(\\d+)",
				"raw https?:\\\\\\/\\\\\\/kpl\\.m\\.jd\\.[comhk]{2,3}\\\\\\/product\\?wareId=(\\d+)",
				"raw https?://wq\\.jd\\.[comhk]{2,3}/item/view\\?sku=(\\d+)",
				"raw https?:\\\\\\/\\\\\\/wq\\.jd\\.[comhk]{2,3}\\\\\\/item\\\\\\/view\\?sku=(\\d+)",
				"raw https?://wqitem\\.jd\\.[comhk]{2,3}.+sku=(\\d+)",
				"raw https?:\\\\\\/\\\\\\/wqitem\\.jd\\.[comhk]{2,3}.+sku=(\\d+)",
				"raw https?://.+\\.jd\\.[comhk]{2,3}.+sku=(\\d+)",
				"raw https?:\\\\\\/\\\\\\/.+\\.jd\\.[comhk]{2,3}.+sku=(\\d+)",
				"raw https?://.+jd\\.[comhk]{2,3}/product/(\\d+).html",
				"raw https?:\\\\\\/\\\\\\/.+jd\\.[comhk]{2,3}\\\\\\/product/(\\d+).html",
				"raw https?://u\\.jd\\.com/(\\w{7})",
				"raw https?:\\\\\\/\\\\\\/u\\.jd\\.com\\\\\\/(\\w{7})"},
			Handle: func(s core.Sender) interface{} {
				//fmt.Println(s.Get())
				return getFanli(s.Get())
			},
		},
	})
	core.OttoFuncs["jingdong"] = getFanli //类似于向核心组件注册
}

func getFanli(url string) string {
	sku := core.Int(url) //从字符串中获取包含的int型数据
	var content string=""
	if sku != 0 {
		content=fmt.Sprintf("https://item.jd.com/%d.html", sku)
	}else{
		content="https://u.jd.com/"+url
	}
	req := httplib.Get("https://api.jingpinku.com/get_rebate_link/api?" +
		"appid=" + jingdong.Get("jingpinku_appid") +
		"&appkey=" + jingdong.Get("jingpinku_appkey") +
		"&union_id=" + jingdong.Get("jd_union_id") +
		"&content=" + content)
	data, err := req.Bytes()
	
	dropErr(err)

	fmt.Println("---------------------------------------------------------")
	fmt.Println(string(data))
	fmt.Println("---------------------------------------------------------")

	short, _ := jsonparser.GetString(data, "content")
	code, _ := jsonparser.GetInt(data, "code")
	if code != 0 {
		// msg, _ := jsonparser.GetString(data, "msg")
		return ``
	}
	official, _ := jsonparser.GetString(data, "official")
	if official == "" {
		return ``
	}
	lines := strings.Split(official, "\n")
	official = ""
	title := ""
	for i, line := range lines {
		if i == 0 {
			title = strings.Trim(regexp.MustCompile("【.*?】").ReplaceAllString(line, ""), " ")
		}
		if !strings.Contains(line, "佣金") {
			official += line + "\n"
		}
	}
	official = strings.Trim(official, "\n")
	//image, _ := jsonparser.GetString(data, "images", "[0]")
	var price string = ""
	var final string = ""
	if res := regexp.MustCompile(`京东价：(.*)\n`).FindStringSubmatch(official); len(res) > 0 {
		price = res[1]
	}
	if res := regexp.MustCompile(`促销价：(.*)\n`).FindStringSubmatch(official); len(res) > 0 {
		final = res[1]
	}
	if math.Abs(core.Float64(price)-core.Float64(final)) < 0.1 {
		final = price
	} else {
		req := httplib.Get("https://api.jingpinku.com/get_powerful_coup_link/api?" +
			"appid=" + jingdong.Get("jingpinku_appid") +
			"&appkey=" + jingdong.Get("jingpinku_appkey") +
			"&union_id=" + jingdong.Get("jd_union_id") +
			"&content=" + content)
		data, _ := req.Bytes()
		quan, _ := jsonparser.GetString(data, "content")
		if strings.Contains(quan, "https://u.jd.com") {
			short = quan
		}
	}
	/*data, _ = json.Marshal(map[string]interface{}{
		"title":    title,
		"short":    short,
		"official": official,
		"price":    price,
		"final":    final,
		"image":    image,
	})*/
    var rslt string=title+"\n京东价："+price+"\n促销价："+final+"\n惠链接："+short
	return string(rslt)
}

// 创建一个错误处理函数，避免过多的 if err != nil{} 出现
func dropErr(e error) {
	if e != nil {
		panic(e)
	}
}
