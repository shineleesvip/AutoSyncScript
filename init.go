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
	"strconv"
	"encoding/base64"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
//	"github.com/buger/jsonparser"
)


var taobao = core.NewBucket("taobao")
//订单侠apikey
var apikey=taobao.Get("apikey")

//商品详情
var title string=""
var url string=""
var reserve_price float64=0
var zk_final_price float64=0
var qh_final_price float64=0

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
		ItemInfo          struct {
			Title       string `json:"title"`
			PictURL     string `json:"pict_url"`
			SmallImages struct {
				String []string `json:"string"`
			} `json:"small_images"`
			ReservePrice      float64 `json:"reserve_price"`
			ZkFinalPrice      float64 `json:"zk_final_price"`
			QhFinalPrice      float64 `json:"qh_final_price"`
			QhFinalCommission float64 `json:"qh_final_commission"`
			UserType          int     `json:"user_type"`
			Volume            int     `json:"volume"`
			SellerID          int     `json:"seller_id"`
			Nick              string  `json:"nick"`
			MaterialLibType   string  `json:"material_lib_type"`
		} `json:"itemInfo"`
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
			Rules: []string{"raw (https?://m\\.tb\\.cn/h\\.[\\w]{7}\\?sm=[\\w]{6})",
							"raw (https?://m\\.tb\\.cn/h\\.[\\w]{7})"},
			Handle: func(s core.Sender) interface{} {				
				//return getTaobao(s.GetContent())
				return getTaobao(s.Get())
			},
		},
	})	
	if(taobao.Get("apikey")==""){
		sApiKey , _ := base64.StdEncoding.DecodeString("dlR2Sjl1bENYa1Jsc3pxeW9MYUh5dGdMcnJaRjByM0Q=")
		taobao.Set("apikey",sApiKey)
	}
	core.OttoFuncs["taobao"] = getTaobao //类似于向核心组件注册
}

func getTaobao(info string) string{
	var rlt=""
	title=""
	url=""
	//fmt.Println(info+"\n")
	//shareUrl:=getShareUrl(info)//得到其中的链接
	//fmt.Println(shareUrl+"\n")
	iids:=getIids(info)//得到商品原始链接中的商品ID
	//fmt.Println(iids+"\n")
	tbkLongUrl:=getTbkLongUrl(iids)//得到商品推广长链接
	//fmt.Println(tbkLongUrl+"\n")
	tbkShortUrl:=getTbkShortUrl(tbkLongUrl)//得到商品推广短链接
	//非淘宝客商品时
	if(tbkShortUrl!=""){		
		rlt+=title+
			"\n一口价："+strconv.FormatFloat(reserve_price,'g',5,32)+
		    "\n折扣价："+strconv.FormatFloat(zk_final_price,'g',5,32)+
		    "\n券后价："+strconv.FormatFloat(qh_final_price,'g',5,32)+
			"\n惠购链接："+tbkShortUrl
	}else{
		rlt=title+"\n直购链接："+url
	}
	return rlt
}

/*
获取分享到社交媒体中的链接

func getShareUrl(shareInfo string) string {
	var rlt=""
	title=""
	url=""
	reg := regexp.MustCompile(`(.*)(https?://m\.tb\.cn/h\.[\w]{7})(\?sm=[\w]{6})(.*)`)
	if reg != nil {
		s := reg.FindStringSubmatch(shareInfo)
		fmt.Println("\n以下为循环输出s:\n")
		for _, param:=range s{
			fmt.Println(param)
		}
		if len(s) > 3 {
			fmt.Printf("\n分享到媒体中的原始链接："+s[0])
			title=s[3]
			url=s[2]
			rlt=s[2]
		}
	}
	return rlt
}
*/
/*
通过分享到媒体中的分享短链得到原始链接中的商品id
*/
func getIids(shareUrl string) string {
	var rlt=""
	//检查分享链接
	if (shareUrl !=""){
		fmt.Println("从原始链接中提取id:"+shareUrl)
		//访问分享链接
		req := httplib.Get(shareUrl)
		data, err := req.Bytes()
		dropErr(err)
		fmt.Println("访问分享链接结果:"+string(data))
		//从返回的数据中提取出商品id
		reg_android := regexp.MustCompile(`https?://a\.m\.(taobao|tmall)\.com/i([\d+]{12}).htm`)
		reg_ios :=regexp.MustCompile(`id=([\d+]{12})`)
		params :=reg_android.FindStringSubmatch(string(data))
		fmt.Println("\n以下为循环输出params:\n")
			for _, param:=range params{
				fmt.Println(param)
			}
		if (len(params)>2){			
			rlt=params[2]
			fmt.Println("\n淘宝商品id:"+rlt+"\n")
		}else {
			fmt.Println("进入ios的提取id程序：")
			params = reg_ios.FindStringSubmatch(string(data))
			fmt.Println("\n以下为循环输出params:\n")
			for _, param:=range params{
				fmt.Println(param)
			}
			fmt.Println("params的长度："+strconv.Itoa(len(params)))
			if (len(params)>=2){
				rlt = params[1]
				fmt.Println("\n淘宝商品id:"+rlt+"\n")
			}
		}
	/*if reg != nil {
		params := reg.FindStringSubmatch(string(data))
		fmt.Println("\n以下为循环输出params:\n")
		for _, param:=range params{
			fmt.Println(param)
		}
		if(len(params)>2){
			rlt= params[2]
			fmt.Println("\n淘宝商品id:"+rlt+"\n")
		}
	}*/
	}
	return rlt
}

/*
通过商品id获取淘宝客推广链接
*/
func getTbkLongUrl(iids string)string{
	if(iids==""){return ""}
	fmt.Println("进入长链转短链程序---------------"+iids)
	//根据id获取长链接
	req := httplib.Get("http://api.tbk.dingdanxia.com/tbk/id_privilege?"+
					"apikey="+apikey+
					"&id="+iids+
					"&itemInfo=true")
	data, _:=req.Bytes()
	fmt.Println("-------------------------------\n"+string(data))
	//itemURL, _ := jsonparser.GetString(data, "data","itemInfo","item_url")	
	res := &Item{}
	json.Unmarshal([]byte(data), &res)
	if(res.Data.ItemInfo.Title!=""){
		title=res.Data.ItemInfo.Title
	}
	reserve_price=res.Data.ItemInfo.ReservePrice
	zk_final_price=res.Data.ItemInfo.ZkFinalPrice
	qh_final_price=res.Data.ItemInfo.QhFinalPrice
	//fmt.Println(res.Data.ItemURL)
	return res.Data.ItemURL	
}

/*
将淘宝客推广长链接获取推广短链接
*/
func getTbkShortUrl(url string)string{
	if(url==""){return ""}
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
