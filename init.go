/*******************************************
QQ交流群：418353744
QQ线报群：263723430
********************************************/

package taobao

import (
	"fmt"
	"regexp"
	"encoding/json"

	"strconv"
	"encoding/base64"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"

)


var taobao = core.NewBucket("taobao")

var apikey=taobao.Get("apikey")


var title string=""
var url string=""
var reserve_price float64=0
var zk_final_price float64=0
var qh_final_price float64=0
var coupon string=""
var coupon_tpwd string=""
var item_tpwd string=""

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
		CouponTpwd			string 		`json:"coupon_tpwd"`
		ItemTpwd			string 		`json:"item_tpwd"`
	} `json:"data"`
}

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
	core.OttoFuncs["taobao"] = getTaobao 
}

func getTaobao(info string) string{
	var rlt=""
	title=""
	url=""
	iids:=getIids(info)//得到商品原始链接中的商品ID
	tbkLongUrl:=getTbkLongUrl(iids)//得到商品推广长链接
	if(tbkLongUrl!=""){		
		rlt+=title+
			"\n一口价："+strconv.FormatFloat(reserve_price,'g',5,32)+
		    "\n折扣价："+strconv.FormatFloat(zk_final_price,'g',5,32)+
		    "\n券后价："+strconv.FormatFloat(qh_final_price,'g',5,32)+
			//"\n券口令："+coupon_tpwd+
			"\n淘口令："+item_tpwd+
			"\n优惠券："+coupon+
			"\n惠链接："+tbkLongUrl
	}
	return rlt
}

func getIids(shareUrl string) string {
	var rlt=""
	if (shareUrl !=""){
		fmt.Println("从原始链接中提取id:"+shareUrl)
		req := httplib.Get(shareUrl)
		data, err := req.Bytes()
		dropErr(err)
		fmt.Println("访问分享链接结果:"+string(data))
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
	}
	return rlt
}
func getTbkLongUrl(iids string)string{
	if(iids==""){return ""}
	fmt.Println("进入长链转短链程序---------------"+iids)
	req := httplib.Get("http://api.tbk.dingdanxia.com/tbk/id_privilege?"+
					"apikey="+apikey+
					"&id="+iids+
					"&itemInfo=true"+
					"&shorturl=true"+
					"&tpwd=true")
	data, _:=req.Bytes()
	fmt.Println("-------------------------------\n"+string(data))
	res := &Item{}
	json.Unmarshal([]byte(data), &res)
	if(res.Data.ItemInfo.Title!=""){
		title=res.Data.ItemInfo.Title
	}
	reserve_price=res.Data.ItemInfo.ReservePrice
	zk_final_price=res.Data.ItemInfo.ZkFinalPrice
	qh_final_price=res.Data.ItemInfo.QhFinalPrice
	coupon=res.Data.CouponClickURL
	coupon_tpwd=res.Data.CouponTpwd
	item_tpwd=res.Data.ItemTpwd
	//fmt.Println(res.Data.ItemURL)
	return res.Data.ItemURL	
}

func dropErr(e error) {
	if e != nil {
		panic(e)
	}
}
