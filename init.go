/*******************************************
QQ交流群：418353744
QQ线报群：263723430
********************************************/

package pinduoduo

import (
	"fmt"
	"regexp"
	"encoding/json"
	"time"
	"crypto/md5"

	"strings"
	"strconv"
	"sort"


	"github.com/buger/jsonparser"
	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
	"github.com/beego/beego/v2/core/logs"

)


var pinduoduo = core.NewBucket("pinduoduo")

var pddSite = "https://gw-api.pinduoduo.com/api/router"
var apitype = ""
var client_id=""
var client_key=""
var pid=""
var timestamp=""

var bind bool=false

var goodTitle=""
var goods_details=""

type ItemQueryBind struct {
	AuthorityQueryResponse struct {
		Bind      int    `json:"bind"`
		RequestID string `json:"request_id"`
	} `json:"authority_query_response"`
}
type UrlBind struct {
	RpPromotionURLGenerateResponse struct {
		URLList []struct {
			MobileURL string `json:"mobile_url"`
			URL       string `json:"url"`
		} `json:"url_list"`
		RequestID string `json:"request_id"`
	} `json:"rp_promotion_url_generate_response"`
}
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

type ItemSign struct {
	GoodsSearchResponse struct {
		GoodsList []struct {
			CategoryName                string        `json:"category_name"`
			CouponRemainQuantity        int           `json:"coupon_remain_quantity"`
			PromotionRate               int           `json:"promotion_rate"`
			ServiceTags                 []int         `json:"service_tags"`
			MallID                      int           `json:"mall_id"`
			MallName                    string        `json:"mall_name"`
			MallCouponEndTime           int           `json:"mall_coupon_end_time"`
			LgstTxt                     string        `json:"lgst_txt"`
			GoodsName                   string        `json:"goods_name"`
			HasMaterial                 bool          `json:"has_material"`
			GoodsID                     int64         `json:"goods_id"`
			BrandName                   string        `json:"brand_name"`
			PredictPromotionRate        int           `json:"predict_promotion_rate"`
			GoodsDesc                   string        `json:"goods_desc"`
			OptName                     string        `json:"opt_name"`
			ShareRate                   int           `json:"share_rate"`
			OptIds                      []int         `json:"opt_ids"`
			GoodsImageURL               string        `json:"goods_image_url"`
			HasMallCoupon               bool          `json:"has_mall_coupon"`
			UnifiedTags                 []interface{} `json:"unified_tags"`
			CouponStartTime             int           `json:"coupon_start_time"`
			MinGroupPrice               float64       `json:"min_group_price"`
			CouponDiscount              int           `json:"coupon_discount"`
			CouponEndTime               int           `json:"coupon_end_time"`
			ZsDuoID                     int           `json:"zs_duo_id"`
			MallCouponRemainQuantity    int           `json:"mall_coupon_remain_quantity"`
			PlanType                    int           `json:"plan_type"`
			CatIds                      []int         `json:"cat_ids"`
			CouponMinOrderAmount        int           `json:"coupon_min_order_amount"`
			CategoryID                  int           `json:"category_id"`
			MallCouponDiscountPct       int           `json:"mall_coupon_discount_pct"`
			ActivityType                int           `json:"activity_type"`
			CouponTotalQuantity         int           `json:"coupon_total_quantity"`
			MallCouponMinOrderAmount    int           `json:"mall_coupon_min_order_amount"`
			MerchantType                int           `json:"merchant_type"`
			SalesTip                    string        `json:"sales_tip"`
			OnlySceneAuth               bool          `json:"only_scene_auth"`
			DescTxt                     string        `json:"desc_txt"`
			MallCouponID                int           `json:"mall_coupon_id"`
			GoodsThumbnailURL           string        `json:"goods_thumbnail_url"`
			OptID                       int           `json:"opt_id"`
			SearchID                    string        `json:"search_id"`
			ActivityTags                []int         `json:"activity_tags"`
			HasCoupon                   bool          `json:"has_coupon"`
			MinNormalPrice              float64       `json:"min_normal_price"`
			MallCouponStartTime         int           `json:"mall_coupon_start_time"`
			ServTxt                     string        `json:"serv_txt"`
			MallCouponTotalQuantity     int           `json:"mall_coupon_total_quantity"`
			MallCouponMaxDiscountAmount int           `json:"mall_coupon_max_discount_amount"`
			MallCps                     int           `json:"mall_cps"`
			GoodsSign                   string        `json:"goods_sign"`
		} `json:"goods_list"`
		ListID     string `json:"list_id"`
		TotalCount int    `json:"total_count"`
		RequestID  string `json:"request_id"`
		SearchID   string `json:"search_id"`
	} `json:"goods_search_response"`
}

func init() {

	core.Server.GET("/pinduoduo/:sku", func(c *gin.Context) {
		sku := c.Param("sku")
		c.String(200, core.OttoFuncs["pinduoduo"](sku))
	})
	core.AddCommand("", []core.Function{
		{
			//Rules: []string{"raw https?://mobile\\.yangkeduo\\.com/goods.?\\.html\\?goods_id=(\\d+)"},
			Rules: []string{"raw mobile\\.yangkeduo\\.com"},
			Handle: func(s core.Sender) interface{} {
				var resMessage=""
				
				client_id = pinduoduo.Get("client_id")
				client_key = pinduoduo.Get("client_key")
				pid = pinduoduo.Get("pid")
				bind = (pinduoduo.Get("bind")=="true")
				
				if(s.IsAdmin()){
					if(bind){
						if(client_id!="" && client_key!="" && pid!=""){
							resMessage=getPinduoduo(s.GetContent())
						}else{
							resMessage="请设置client_id、client_key、pid必要信息"
						}
					}else {
						bind=queryBind()
						pinduoduo.Set("bind",strconv.FormatBool(bind))
						fmt.Sprintf("绑定结果："+strconv.FormatBool(bind))
						if (!bind && client_id!="" && client_key!="" && pid!=""){
							resMessage= "点击链接授权备案:\n"+setBind()
						}else if (bind && client_id!="" && client_key!="" && pid!=""){
							resMessage=getPinduoduo(s.GetContent())
						}else{
							resMessage="请设置client_id、client_key、pid必要信息"
						}
					}
				}else{
					if(client_key!=""&&client_id!=""&&pid!=""&&bind){
						resMessage=getPinduoduo(s.GetContent())
					}else{
						resMessage=""
					}
				}
				return resMessage
			},
		},
	})
	core.OttoFuncs["pinduoduo"] = getPinduoduo 
	logs.Info("拼多多佣金短链启动：关注QQ群418353744获取更多消息")
}


func queryBind() bool{
	params:=map[string]string{
		"type":"pdd.ddk.member.authority.query",
		"client_id":client_id,
		"pid":pid,
   		"timestamp": strconv.FormatInt(time.Now().Unix(),10),
	}
   	sign:=getMd5(client_key,params)
	data:=accessApi(pddSite,params,sign)
	
	bind , _ :=jsonparser.GetInt([]byte(data),"authority_query_response","bind")
	fmt.Println("检测是否完成授权备案："+strconv.FormatInt(bind,10))
	return bind==1	
}

func setBind() string{
	if(client_id!=""&&client_key!=""&&pid!=""){
		params:=map[string]string{
			"type": 		"pdd.ddk.rp.prom.url.generate",
			"client_id": 	client_id,
			"p_id_list": 	"[\""+pid+"\"]",
			//"p_id_list": 	"%5B%22"+pinduoduo.Get("pid")+"%22%5D",
			"data_type":    "JSON",
			"channel_type": "10",
			"timestamp": strconv.FormatInt(time.Now().Unix(),10),
		}
		sign:=getMd5(client_key,params)
		fmt.Println("------------------------------------------------------------")
		fmt.Println("sign:"+sign)
		data:=accessApi(pddSite,params,sign)
		res := &UrlBind{}
		json.Unmarshal(data, &res)
		fmt.Println("------------------------------------------------------------")
		if ( len(res.RpPromotionURLGenerateResponse.URLList)>=1 ){
			return string(res.RpPromotionURLGenerateResponse.URLList[0].MobileURL)
		}else{
			return "自动授权备案失败，请手动授权备案https://open.pinduoduo.com/application/document/apiTools?scopeName=pdd.ddk.rp.prom.url.generate"
		}
	}else {
		return "缺少必要的设置client_id、client_key、pid"
	}
}

func getPinduoduo(info string) string{
	reg := regexp.MustCompile(`<title>(.*)</title>`)
	if reg != nil {
		params := reg.FindStringSubmatch(string(info))
		if(len(params)>1){
			goodTitle = params[1]
			fmt.Println("商品名称:"+goodTitle+"\n")
		}
	}
	var source_url=""
	reg = regexp.MustCompile(`goods_id=(\d+)`)
	if reg != nil {
			params := reg.FindStringSubmatch(string(info))
			if(len(params)>=2){
			source_url = "https://mobile.yangkeduo.com/goods.html?goods_id="+params[1]
			fmt.Println("链接:"+source_url+"\n")
			goods_id:=params[1]		
			goods_details =getGoodsDetails(goods_id)		
		}
	}
	short_url:=""
	if (goods_details!=""){
		short_url = goods_details +"\n惠购链接："+getShortUrl(source_url)
	}else{
		short_url = goodTitle + "\n直购链接："+source_url
	}
	return short_url
}

func getGoodsDetails(goods_id string)string{
	params :=map[string]string {
		"type":"pdd.ddk.goods.search",
		"client_id": client_id,
		"pid": pid,
		"timestamp":strconv.FormatInt(time.Now().Unix(),10),
		"keyword": goods_id,
	}

	upMd5:=getMd5(client_key,params)

	data := accessApi(pddSite,params,upMd5)
	fmt.Println("商品详情："+string(data))
	goodName, _ := jsonparser.GetString(data, "goods_name")
	fmt.Println("商品名称："+ goodName)
	res:=&ItemSign{}
	json.Unmarshal([]byte(data),&res)
	//fmt.Println(res.GoodsZsUnitGenerateResponse.ShortUrl)
	if ( len(res.GoodsSearchResponse.GoodsList)==0){
		return ""
	}
	rlt :=""
	if(res.GoodsSearchResponse.GoodsList[0].GoodsName!=""){
		rlt+=(res.GoodsSearchResponse.GoodsList[0].GoodsName)			
	}	
	if(res.GoodsSearchResponse.GoodsList[0].MinGroupPrice!=0){
		rlt+="\n最小拼团价："+strconv.FormatFloat((res.GoodsSearchResponse.GoodsList[0].MinGroupPrice/100),'g',5,32)			
	}
	if(res.GoodsSearchResponse.GoodsList[0].MinGroupPrice!=0){
		rlt+="\n最小单买价："+strconv.FormatFloat((res.GoodsSearchResponse.GoodsList[0].MinNormalPrice/100),'g',5,32)
	}
	return rlt
}

func getShortUrl(source_url string) string {
	var rlt=""
	if(source_url!=""){
		params :=map[string]string {
			"type":		"pdd.ddk.goods.zs.unit.url.gen",
			"client_id": 	client_id,
			"pid":			pid,
			"timestamp":	strconv.FormatInt(time.Now().Unix(),10),
			"source_url":	source_url,
		}

		upMd5:=getMd5(client_key,params)

		data:=accessApi(pddSite,params,upMd5)
		fmt.Println("getShortUrl中得到的接口返回值："+string(data))
		res:=&ItemUrl{}
		json.Unmarshal([]byte(data),&res)
		rlt = string(res.GoodsZsUnitGenerateResponse.ShortURL)
	} 
	return rlt
}


func accessApi(site string,params map[string]string,sign string) []byte{
	var dataParams string=""
	for key :=range params{
		dataParams += key+"="+params[key]+"&"
	}
	dataParams +="sign="+getMd5(client_key,params)
	urlstr := site+"?"+dataParams
	fmt.Println(urlstr)
	req := httplib.Get(urlReplace(urlstr))
	rlt,err := req.Bytes()
	dropErr(err)
	return rlt
}

func getMd5(client_key string,params map[string]string) string{
	var dataParams string
	var keys []string
	for k :=range params{
		keys=append(keys,k)
	}
	sort.Strings(keys)
	for _, k:=range keys{
		fmt.Println("key:",k,"value:",params[k])
		dataParams +=k+params[k]
	}
	dataParams=client_key+dataParams+client_key
	fmt.Println("MD5函数拼接的字符串："+dataParams)
	strMd5:=md5.Sum([]byte(dataParams))
	upMd5:=strings.ToUpper(fmt.Sprintf("%x",strMd5))
	fmt.Println("MD5函数获取的MD5值："+upMd5)
	return upMd5
}

func urlReplace(url string) string{
	urlstr:=strings.Replace(url,"[","%5B",-1)
	urlstr=strings.Replace(urlstr,"]","%5D",-1)
	urlstr=strings.Replace(urlstr,"\"","%22",-1)
	return urlstr
}

func dropErr(e error) {
	if e != nil {
		panic(e)
	}
}
