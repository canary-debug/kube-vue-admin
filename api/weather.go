package api

import (
	"encoding/json"
	"fmt"
	config "github.com/canary-debug/kube-vue-admin/configs"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"log"
)

// 定义百度天气API返回数据的结构体（包含温度和穿衣指数）
type BaiduWeatherResponse struct {
	Status int `json:"status"`
	Result struct {
		Now struct {
			Temp int `json:"temp"` // 实时温度
		} `json:"now"`
		Indexes []struct { // 生活指数列表
			Name   string `json:"name"`   // 指数名称（如"穿衣指数"）
			Brief  string `json:"brief"`  // 简要建议（如"较冷"）
			Detail string `json:"detail"` // 详细说明（如"建议着厚外套加毛衣等服装"）
		} `json:"indexes"`
	} `json:"result"`
}

func Weather(c *gin.Context) {
	// 封装 API 地址
	url := fmt.Sprintf("https://api.map.baidu.com/weather/v1/?district_id=%s&data_type=all&ak=%s", config.C.Weather.DistrictID, config.C.Weather.ApiKey)
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "获取天气信息失败",
			"error": err.Error(),
		})
		log.Println("获取天气失败: ", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "读取天气数据失败",
			"error": err.Error(),
		})
		log.Println("读取响应体失败: ", err)
		return
	}

	// 解析JSON数据到结构体
	var weatherResp BaiduWeatherResponse
	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "解析天气数据失败",
			"error": err.Error(),
		})
		log.Println("解析天气JSON失败: ", err)
		return
	}

	// 从指数列表中筛选出「穿衣指数」
	var dressingIndex struct {
		Name   string `json:"name"`
		Brief  string `json:"brief"`
		Detail string `json:"detail"`
	}
	for _, index := range weatherResp.Result.Indexes {
		if index.Name == "穿衣指数" {
			dressingIndex = index
			break
		}
	}

	// 判断时间并返回 早安、午安、晚安
	var timeOfDay string
	switch {
	case 6 <= time.Now().Hour() && time.Now().Hour() < 12:
		timeOfDay = "早安"
	case 12 <= time.Now().Hour() && time.Now().Hour() < 18:
		timeOfDay = "中午好"
	default:
		timeOfDay = "晚上好"
	}

	// 返回温度 + 穿衣指数
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取天气信息成功",
		"data": gin.H{
			"time_of_day":     timeOfDay,                   // 时间段
			"temp":            weatherResp.Result.Now.Temp, // 实时温度
			"dressing":        dressingIndex.Brief,         // 穿衣建议（简要）
			"dressing_detail": dressingIndex.Detail,        // 穿衣建议（详细，可选删除）
		},
	})
}
