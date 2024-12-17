package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zachturing/login/common/define"
	"github.com/zachturing/util/log"
	"io"
	"net/http"
)

type BaiduUploadConvertData struct {
	Token           string `json:"token"`
	ConversionTypes []struct {
		LogidUrl        string `json:"logidUrl"`
		NewType         int    `json:"newType"`
		AttributeSource int    `json:"attributeSource"`
	} `json:"conversionTypes"`
}

type BaiduUploadConvertDataResponse struct {
	Header struct {
		Desc   string `json:"desc"`
		Errors []struct {
			Code     int    `json:"code"`
			Message  string `json:"message"`
			Position string `json:"position"`
		} `json:"errors"`
		Status int `json:"status"`
	} `json:"header"`
}

type BaiduUploadRequest struct {
	OutTradeNo string `json:"out_trade_no" binding:"required"` // 订单编号，对应支付宝商家编号
	BaiduVid   string `json:"bd_vid"`                          // 通过百度广告点击进入的，回传是需要带上百度提供的vid
}

// CallBaiduUploadConvertData 新注册的用户，调用百度推广回传数据接口
func CallBaiduUploadConvertData(bdVid string) error {
	// 如果bd_vid为空，直接返回
	if bdVid == "" {
		return fmt.Errorf("bd_vid is empty")
	}
	log.Debugf("bd_vid is:%s", bdVid)

	// 构造百度回传数据api的请求数据
	data := prepareBaiduUploadData(bdVid, define.BaiduApiToken)

	// 调用百度回传数据api，有报错则打印log，但是返回成功响应
	if err := sendBaiduConvertData(define.BaiduAPIURL, data); err != nil {
		log.Errorf("Failed to send Baidu convert data: %v", err)
		return err
	}

	return nil
}

func prepareBaiduUploadData(baiduVid, apiToken string) BaiduUploadConvertData {
	return BaiduUploadConvertData{
		Token: apiToken,
		ConversionTypes: []struct {
			LogidUrl        string `json:"logidUrl"`
			NewType         int    `json:"newType"`
			AttributeSource int    `json:"attributeSource"`
		}{
			{
				LogidUrl:        fmt.Sprintf("https://www.mixpaper.cn?bd_vid=%s", baiduVid),
				NewType:         25, // 百度定义的类型，25-注册
				AttributeSource: 0,
			},
		},
	}
}

func sendBaiduConvertData(url string, data BaiduUploadConvertData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("JSON marshal failed: %w", err)
	}

	convertReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	convertReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(convertReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var response BaiduUploadConvertDataResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if response.Header.Status != 0 {
		return fmt.Errorf("failed to upload convert data, desc: %s", response.Header.Desc)
	}

	return nil
}
