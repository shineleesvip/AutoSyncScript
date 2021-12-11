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
	"sort"
//	"net/url"

	"github.com/buger/jsonparser"
	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
//	"github.com/buger/jsonparser"
)


var pinduoduo = core.NewBucket("pinduoduo")
//拼多多
var pddSite = "https://gw-api.pinduoduo.com/api/router"
var apitype = ""
var client_id=""
var client_key=""
var pid=""
var timestamp=""
//是否完成推广位的媒体id的绑定
var bind bool=false
//商品
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
	//添加命令
	core.AddCommand("", []core.Function{
		{
			//Rules: []string{"raw https?://mobile\\.yangkeduo\\.com/goods.?\\.html\\?goods_id=(\\d+)"},
			Rules: []string{"raw mobile\\.yangkeduo\\.com"},
			Handle: func(s core.Sender) interface{} {
				var resMessage=""
				//获取必要信息
				client_id = pinduoduo.Get("client_id")
				client_key = pinduoduo.Get("client_key")
				pid = pinduoduo.Get("pid")
				bind = (pinduoduo.Get("bind")=="true")
				//发送信息的是管理员
				if(s.IsAdmin()){
					if(bind){//完成授权备案的
						if(client_id!="" && client_key!="" && pid!=""){//设置了必要参数的
							resMessage=getPinduoduo(s.GetContent())
						}else{//没设置必要参数
							resMessage="请设置client_id、client_key、pid必要信息"
						}
					}else {//未完成授权备案的
						bind=queryBind()
						//bind=false
						pinduoduo.Set("bind",strconv.FormatBool(bind))
						fmt.Sprintf("绑定结果："+strconv.FormatBool(bind))
						if(!bind){
							resMessage= "点击链接授权备案:\n"+setBind()
						}else{
							//fmt.
							resMessage=getPinduoduo(s.GetContent())
						}
					}
				}else{//发送信息的不是管理员
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
	core.OttoFuncs["pinduoduo"] = getPinduoduo //类似于向核心组件注册
}
//授权备案：https://jinbao.pinduoduo.com/qa-system?questionId=218
//查询是否绑定
func queryBind() bool{
	//对公共参数和业务参数按照ASCII排序,不参加排序：app_secret和sign
	params:=map[string]string{
		"type":"pdd.ddk.member.authority.query",
		"client_id":client_id,
		"pid":pid,
   		"timestamp": strconv.FormatInt(time.Now().Unix(),10),
	}
	//client_key:=pinduoduo.Get("client_key")
   	//大写(MD5(client_secret+key1+value1+key2+value2+client_secret))
   	//将排序好的参数名和参数值拼装在一起，两头加client_key
   	sign:=getMd5(client_key,params)
	data:=accessApi(pddSite,params,sign)
	
	bind , _ :=jsonparser.GetInt([]byte(data),"authority_query_response","bind")
	fmt.Println("检测是否完成授权备案："+strconv.FormatInt(bind,10))
	return bind==1
	
}

func setBind() string{
	//对公共参数和业务参数按照ASCII排序,不参加排序：app_secret和sign
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
		//client_key:=pinduoduo.Get("client_key")
		//大写(MD5(client_secret+key1+value1+key2+value2+client_secret))
		sign:=getMd5(client_key,params)
		fmt.Println("------------------------------------------------------------")
		fmt.Println("sign:"+sign)
		//params["p_id_list"]="%5B%22"+pinduoduo.Get("pid")+"%22%5D"
		data:=accessApi(pddSite,params,sign)
		//fmt.Println("绑定返回值："+string(data))
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
	//从分享到媒体中的信息提取title
	reg := regexp.MustCompile(`<title>(.*)</title>`)
	if reg != nil {
		params := reg.FindStringSubmatch(string(info))
		if(len(params)>1){
			goodTitle = params[1]
			fmt.Println("商品名称:"+goodTitle+"\n")
		}
	}
	//从返回的数据中提取出商品id
	var source_url=""
	//reg = regexp.MustCompile(`https?://mobile\.yangkeduo\.com/goods.?\.html\?goods_id=(\d+)`)
	reg = regexp.MustCompile(`goods_id=(\d+)`)
	if reg != nil {
			params := reg.FindStringSubmatch(string(info))
			if(len(params)>=2){
			source_url = "https://mobile.yangkeduo.com/goods.html?goods_id="+params[1]
			fmt.Println("链接:"+source_url+"\n")
			goods_id:=params[1]		
			//通过goods_id获取goods_sign
			goods_details =getGoodsDetails(goods_id)
			//fmt.Println("\n商品goods_sign:"+goods_sign+"\n")		
		}
	}
	short_url:=""
	if (goods_details!=""){
		short_url = goods_details +"\n惠购链接"+getShortUrl(source_url)
	}else{
		short_url = goodTitle + "\n惠购链接"+getShortUrl(source_url)
	}
	return short_url
}

//通过goods_id获取商品详情及goods_sign
func getGoodsDetails(goods_id string)string{
	params :=map[string]string {
		"type":"pdd.ddk.goods.search",
		"client_id": client_id,
		"pid": pid,
		"timestamp":strconv.FormatInt(time.Now().Unix(),10),
		"keyword": goods_id,
	}
	//client_key=pinduoduo.Get("client_key")
	upMd5:=getMd5(client_key,params)
	//将长链接变换成短链接
	data := accessApi(pddSite,params,upMd5)
	fmt.Println("商品详情："+string(data))
	goodName, _ := jsonparser.GetString(data, "goods_name")
	fmt.Println("商品名称："+ goodName)
	res:=&ItemSign{}
	json.Unmarshal([]byte(data),&res)
	//fmt.Println(res.GoodsZsUnitGenerateResponse.ShortUrl)
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

//获取短链接
func getShortUrl(source_url string) string {
	var rlt=""
	if(source_url!=""){
		//对公共参数和业务参数按照ASCII排序,不参加排序：app_secret和sign
		params :=map[string]string {
			"type":		"pdd.ddk.goods.zs.unit.url.gen",
			"client_id": 	client_id,
			"pid":			pid,
			"timestamp":	strconv.FormatInt(time.Now().Unix(),10),
			"source_url":	source_url,
		}
		//client_key=pinduoduo.Get("client_key")
		//MD5
		upMd5:=getMd5(client_key,params)
		//将长链接变换成短链接
		data:=accessApi(pddSite,params,upMd5)
		fmt.Println("getShortUrl中得到的接口返回值："+string(data))
		res:=&ItemUrl{}
		json.Unmarshal([]byte(data),&res)
		//fmt.Println(res.GoodsZsUnitGenerateResponse.ShortUrl)
		rlt = string(res.GoodsZsUnitGenerateResponse.ShortURL)
	} 
	return rlt
}

//访问接口
func accessApi(site string,params map[string]string,sign string) []byte{
	var dataParams string=""
	//准备参数
	for key :=range params{
		dataParams += key+"="+params[key]+"&"
	}
	dataParams +="sign="+getMd5(client_key,params)
	//真实访问
	urlstr := site+"?"+dataParams
	//escapeUrl := url.QueryEscape(urlstr)
	fmt.Println(urlstr)
	req := httplib.Get(urlReplace(urlstr))
	rlt,err := req.Bytes()
	//str,err:=req.String()
	dropErr(err)
	return rlt
}

//获取md5的值
func getMd5(client_key string,params map[string]string) string{
	var dataParams string
	var keys []string
	//从map中提取所有的key
	for k :=range params{
		keys=append(keys,k)
	}
	//对keys进行排序
	sort.Strings(keys)
	//拼接
	for _, k:=range keys{
		fmt.Println("key:",k,"value:",params[k])
		dataParams +=k+params[k]
	}
	dataParams=client_key+dataParams+client_key
	fmt.Println("MD5函数拼接的字符串："+dataParams)
	//md5
	strMd5:=md5.Sum([]byte(dataParams))
	upMd5:=strings.ToUpper(fmt.Sprintf("%x",strMd5))
	fmt.Println("MD5函数获取的MD5值："+upMd5)
	return upMd5
}

//替换url中的特殊字符
func urlReplace(url string) string{
	urlstr:=strings.Replace(url,"[","%5B",-1)
	urlstr=strings.Replace(urlstr,"]","%5D",-1)
	urlstr=strings.Replace(urlstr,"\"","%22",-1)
	return urlstr
}

// 创建一个错误处理函数，避免过多的 if err != nil{} 出现
func dropErr(e error) {
	if e != nil {
		panic(e)
	}
}