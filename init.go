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

type ItemSetBind struct {
	PIDBindResponse struct {
		Result struct {
			Msg    string `json:"msg"`
			Result bool   `json:"result"`
		} `json:"result"`
	} `json:"p_id_bind_response"`
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
			Rules: []string{"raw https?://mobile\\.yangkeduo\\.com/goods.?\\.html\\?goods_id=(\\d+)"},
			Handle: func(s core.Sender) interface{} {
				//fmt.Println(s.GetContent())
				//查询是否绑定
				if(s.IsAdmin()){
					bind=queryBind()
					fmt.Sprintf("绑定结果："+strconv.FormatBool(bind))
					if(!bind){
						if(setBind()){
							bind=true
							return "已完成绑定备案，请再次发送链接或分享"
						}else{
							bind=false
							return "自动绑定失败，请到网站自行绑定备案！"
						}
					}
				}
				pinduoduo.Set("bind",bind)
				return getPinduoduo(s.GetContent())
			},
		},
	})
	core.OttoFuncs["pinduoduo"] = getPinduoduo //类似于向核心组件注册
}

//查询是否绑定
func queryBind() bool{
	//对公共参数和业务参数按照ASCII排序,不参加排序：app_secret和sign
	apitype:="pdd.ddk.member.authority.query"
	client_id:=pinduoduo.Get("client_id")
	pid:=pinduoduo.Get("pid")
   	client_key:=pinduoduo.Get("client_key")
   	//Unix时间timestamp
   	timestamp:= strconv.FormatInt(time.Now().Unix(),10)
   	//大写(MD5(client_secret+key1+value1+key2+value2+client_secret))
   	//将排序好的参数名和参数值拼装在一起，两头加client_key
   	strCon:=client_key+
		   "client_id"+client_id+
		   "pid"+pid+
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
					"&pid="+pid)
	data, _:=req.Bytes()
	//fmt.Println(string(data))
	res:=&ItemQueryBind{}
	json.Unmarshal([]byte(data),&res)
	//fmt.Println(res.GoodsZsUnitGenerateResponse.ShortUrl)
	if(res.AuthorityQueryResponse.Bind==0){
		return false
	}else{
		return true
	}
}

func setBind() bool{
	//对公共参数和业务参数按照ASCII排序,不参加排序：app_secret和sign
	apitype:="pdd.ddk.oauth.pid.mediaid.bind"
	client_id:=pinduoduo.Get("client_id")
	pid:=pinduoduo.Get("pid")
   	client_key:=pinduoduo.Get("client_key")
   	//Unix时间timestamp
   	timestamp:= strconv.FormatInt(time.Now().Unix(),10)
   	//大写(MD5(client_secret+key1+value1+key2+value2+client_secret))
   	//将排序好的参数名和参数值拼装在一起，两头加client_key
   	strCon:=client_key+
		   "client_id"+client_id+
		   "pid"+pid+
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
					"&pid="+pid)
	data, _:=req.Bytes()
	//fmt.Println(string(data))
	res:=&ItemSetBind{}
	json.Unmarshal([]byte(data),&res)
	//fmt.Println(res.GoodsZsUnitGenerateResponse.ShortUrl)
	if(res.PIDBindResponse.Result.Result){
		return true
	}else{
		return false
	}
}


func getPinduoduo(info string) string{
	//从数据中提取title
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
	reg = regexp.MustCompile(`https?://mobile\.yangkeduo\.com/goods.?\.html\?goods_id=(\d+)`)
	if reg != nil {
		params := reg.FindStringSubmatch(string(info))
		source_url = params[0]
		fmt.Println("链接:"+source_url+"\n")
		goods_id:=params[1]		
		//通过goods_id获取goods_sign
		goods_details =getGoodsDetails(goods_id)
		//fmt.Println("\n商品goods_sign:"+goods_sign+"\n")
		
	}	
	return goodTitle + goods_details +"\n购买链接"+getShortUrl(source_url)
}

//通过goods_id获取goods_sign
func getGoodsDetails(goods_id string)string{
	//对公共参数和业务参数按照ASCII排序,不参加排序：app_secret和sign
	apitype="pdd.ddk.goods.search"
	client_id=pinduoduo.Get("client_id")
	pid=pinduoduo.Get("pid")
	client_key=pinduoduo.Get("client_key")
	//Unix时间timestamp
	timestamp= strconv.FormatInt(time.Now().Unix(),10)
	//大写(MD5(client_secret+key1+value1+key2+value2+client_secret))
	//将排序好的参数名和参数值拼装在一起，两头加client_key
	strCon:=client_key+
		"client_id"+client_id+
		"keyword"+goods_id+
		"pid"+pid+
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
				"&keyword="+goods_id)
	data, _:=req.Bytes()
	//fmt.Println(string(data))
	res:=&ItemSign{}
	json.Unmarshal([]byte(data),&res)
	//fmt.Println(res.GoodsZsUnitGenerateResponse.ShortUrl)
	rlt :=""
	
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
	//对公共参数和业务参数按照ASCII排序,不参加排序：app_secret和sign
	apitype="pdd.ddk.goods.zs.unit.url.gen"
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