/****************************************请求地址***********************************
请求地址
http://gw.api.taobao.com/router/rest
https://eco.taobao.com/router/rest
*******************************************************************************************/

/****************************************公共参数***********************************
method			String	是	API接口名称。
app_key			String	是	TOP分配给应用的AppKey。
target_app_key	String	否	被调用的目标AppKey，仅当被调用的API为第三方ISV提供时有效。
sign_method		String	是	签名的摘要算法，可选值为：hmac，md5。
sign			String	是	API输入参数签名结果，签名算法介绍请点击这里。
session			String	否	用户登录授权成功后，TOP颁发给应用的授权信息，详细介绍请点击这里。当此API的标签上注明：“需要授权”，则此参数必传；“不需要授权”，则此参数不需要传；“可选授权”，则此参数为可选。
timestamp		String	是	时间戳，格式为yyyy-MM-dd HH:mm:ss，时区为GMT+8，例如：2015-01-01 12:00:00。淘宝API服务端允许客户端请求最大时间误差为10分钟。
format			String	否	响应格式。默认为xml格式，可选值：xml，json。
v				String	是	API协议版本，可选值：2.0。
partner_id		String	否	合作伙伴身份标识。
simplify 		Boolean	否	是否采用精简JSON返回格式，仅当format=json时有效，默认值为：false。
**************************************************************************************/

/****************************************业务参数**************************************
fields			String	必须	num_iid,click_url		需返回的字段列表
num_iids		String	必须	123,456		商品ID串，用','分割，从taobao.tbk.item.get接口获取num_iid字段，最大40个
adzone_id		Number	必须	123		广告位ID，区分效果位置
platform		Number	可选	123	默认值：1 链接形式：1：PC，2：无线，默认：１
unid			String	可选	demo		自定义输入串，英文和数字组成，长度不能大于12个字符，区分不同的推广渠道
dx				String	可选	1		1表示商品转通用计划链接，其他值或不传表示转营销计划链接
*****************************************************************************************/

/****************************************响应参数***************************************
results		NTbkItem []		淘宝客商品
└ num_iid		Number	123商品ID
└ click_url		String	http://s.click.taobao.com/e=xxx淘客地址
****************************************************************************************/

package taobao

import (
	"fmt"
	"regexp"
	"encoding/json"
//	"time"
//	"crypto/md5"
//	"encoding/hex"
//	"unicode/utf8"
//	"strings"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
//	"github.com/buger/jsonparser"
)


var taobao = core.NewBucket("taobao")
//订单侠apikey
var apikey=taobao.Get("apikey")

//淘宝商品结构体
type Item struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		CategoryID        int    `json:"category_id"`
		CouponClickURL    string `json:"coupon_click_url"`
		CouponEndTime     string `json:"coupon_end_time"`
		CouponInfo        string `json:"coupon_info"`
		CouponRemainCount int    `json:"coupon_remain_count"`
		CouponStartTime   string `json:"coupon_start_time"`
		CouponTotalCount  int    `json:"coupon_total_count"`
		ItemID            int64  `json:"item_id"`
		ItemURL           string `json:"item_url"`
		MaxCommissionRate string `json:"max_commission_rate"`
		RewardInfo        int    `json:"reward_info"`
		Coupon            string `json:"coupon"`
	} `json:"data"`
}

//推广短链接
type ShortUrl struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Content string `json:"content"`
		ErrMsg  string `json:"err_msg"`
	} `json:"data"`
}

func init() {

	core.Server.GET("/taobao/:sku", func(c *gin.Context) {
		sku := c.Param("sku")
		c.String(200, core.OttoFuncs["taobao"](sku))
	})
	//添加命令
	core.AddCommand("", []core.Function{
		{
			Rules: []string{"raw https?://m\\.tb\\.cn/h\\.[\\w]{7}\\?sm=[\\w]{6}"},
			Handle: func(s core.Sender) interface{} {				
				return getTaobao(s.GetContent())
			},
		},
	})
	core.OttoFuncs["taobao"] = getTaobao //类似于向核心组件注册
}

func getTaobao(info string) string{
	fmt.Println(info+"\n")
	shareUrl:=getShareUrl(info)//得到其中的链接
	fmt.Println(shareUrl+"\n")
	iids:=getIids(shareUrl)//得到商品原始链接
	fmt.Println(iids+"\n")
	tbkLongUrl:=getTbkLongUrl(iids)//得到商品推广长链接
	fmt.Println(tbkLongUrl+"\n")
	tbkShortUrl:=getTbkShortUrl(tbkLongUrl)//得到商品推广短链接
	fmt.Println(tbkShortUrl+"\n")
	return tbkShortUrl
}

/*
获取分享到社交媒体中的链接
*/
func getShareUrl(shareInfo string) string {
		reg := regexp.MustCompile(`https?://m\.tb\.cn/h\.[\w]{7}\?sm=[\w]{6}`)
		if reg != nil {
			s := reg.FindStringSubmatch(shareInfo)
			if len(s) > 0 {
			//	fmt.Printf(s[0])
				return s[0]
			}
		}
		return ""
}

/*
通过分享到媒体中的分享短链得到原始链接中的商品id
*/
func getIids(shareUrl string) string {
	req := httplib.Get(shareUrl)
	data, err := req.Bytes()
	if err != nil {
		return `{}`
	}
	//从返回的数据中提取出商品id
	reg := regexp.MustCompile(`https?://a\.m\.taobao\.com/i([\d+]{12}).htm`)
	if reg != nil {
		params := reg.FindStringSubmatch(string(data))
		iids := params[1]
		//fmt.Println("商品id:"+num_iids+"\n")
		return iids
	}
	return ""
}

/*
通过商品id获取淘宝客推广链接
*/
func getTbkLongUrl(iids string)string{
	//根据id获取长链接
	req := httplib.Get("http://api.tbk.dingdanxia.com/tbk/id_privilege?"+
					"apikey="+apikey+
					"&id="+iids)	
	data, _:=req.Bytes()
	fmt.Println(string(data))
	res := &Item{}
	json.Unmarshal([]byte(data), &res)
	fmt.Println(res.Data.ItemURL)
	return res.Data.ItemURL	
}

/*
将淘宝客推广长链接获取推广短链接
*/
func getTbkShortUrl(url string)string{
	//将长链接变换成短链接
	req := httplib.Get("http://api.tbk.dingdanxia.com/tbk/spread_get?"+
					"apikey="+apikey+
					"&url="+url)
	data, _:=req.Bytes()
	//fmt.Println(string(data))
	res:=&ShortUrl{}
	json.Unmarshal([]byte(data),&res)
	return res.Data.Content
}

// 创建一个错误处理函数，避免过多的 if err != nil{} 出现
func dropErr(e error) {
	if e != nil {
		panic(e)
	}
}

/*
//获取md5签名
func getMd5Sign() string {
	//对公共参数和业务参数按照ASCII排序
    //不参加排序：app_secret和sign
	//1.adzone_id
	//2.app_key
	3.fields
	4.method
	5.num_iids
	7.sign_method   
	8.timestamp  
	9.v    

	//将排序好的参数名和参数值拼装在一起
	strCon:="adzone_id"+string(adzone_id)+
			"app_key"+app_key+
			"fields"+fields+
			"method"+method+
			"num_iids"+num_iids+
			"sign_method"+sign_method+
			"timestamp"+timestamp+
			"v"+v
	fmt.Println("拼接的字符串："+strCon)
	//把拼装好的字符串采用utf-8编码，使用签名算法对编码后的字节流进行摘要。如果使用MD5算法，则需要在拼装的字符串前后加上app的secret后，再进行摘要，如：md5(secret+bar2foo1foo_bar3foobar4+secret)；如果使用HMAC_MD5算法，则需要用app的secret初始化摘要算法后，再进行摘要，如：hmac_md5(bar2foo1foo_bar3foobar4)。
	strCon = app_secret + strCon + app_secret
	fmt.Println("字符串两头加app_secret:"+strCon)
	//将摘要得到的字节流结果使用十六进制表示，如：hex("helloworld".getBytes("utf-8")) = "68656C6C6F776F726C64"
	str:=mahonia.ConvertString(strCon)
	md5Rlt := Md5String(str)
	return md5Rlt
}

func Md5String(data string) string{
	md5:=md5.New()
	md5.Write([]byte(data))
	md5Data:=md5.Sum([]byte(nil))
	return hex.EncodeToString(md5Data)
}
//
func getConvert() string{
	url:="http://gw.api.taobao.com/router/rest?"+
	"method=taobao.tbk.item.convert"+
	"&app_key="+app_key+
	"&sign_method="+sign_method+
	"&sign="+sign+
	"&timestamp="+timestamp+
	"&v="+v+
	"&fields="+fields+
	"&num_iids="+num_iids+
	"&adzone_id="+adzone_id
	fmt.Println("进入链接转换步骤,准备访问公共接口\n")
	fmt.Println(url+"\n")
	req := httplib.Get(url)	
	data, err := req.Bytes()
	if err != nil {
		return `{}`
	}
	return string(data)
}

*/