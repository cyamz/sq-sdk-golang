package sqapi

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"sort"
	"strings"
	"time"
)

type SqApi struct {
	HttpUrl      string
	ClientCode   string
	ClientSecret string
	Version      string
}

type ResponseResult struct {
	ErrorCode  string      `json:"error_code"`
	Message    string      `json:"message"`
	ReturnData interface{} `json:"return_data"`
}

func NewSqApi(ClientCode string, ClientSecret string, Version_optional ...string) SqApi {
	Version := "v1"
	if len(Version_optional) > 0 {
		Version = Version_optional[0]
	}
	sqapi := SqApi{
		HttpUrl:      "http://oms.sq-exp.com/",
		ClientCode:   ClientCode,
		ClientSecret: ClientSecret,
		Version:      Version,
	}

	return sqapi
}

func NewResponse(errorCode string, message string, returnData interface{}) ResponseResult {
	return ResponseResult{
		ErrorCode:  errorCode,
		Message:    message,
		ReturnData: returnData,
	}
}

func (api *SqApi) GetSign(data map[string]interface{}, date string) string {
	dataJson, _ := json.Marshal(data)
	json := string(dataJson)

	arr := []string{
		json,
		date,
		api.ClientSecret,
	}
	sort.Strings(arr)

	signStr := "[\"" + arr[0] + "\",\"" + arr[1] + "\"," + arr[2] + "]"
	md5Sum := md5.Sum([]byte(signStr))
	md5Str := hex.EncodeToString(md5Sum[:])
	result := strings.ToLower(string(md5Str[:]))

	return result
}

func (api *SqApi) Request(method string, params map[string]interface{}) (res ResponseResult, err error) {
	res = ResponseResult{}

	url := api.HttpUrl + api.Version + "/" + method
	date := time.Now().Format("2006-01-02 15:04:05")

	sign := api.GetSign(params, date)
	requestArr := map[string]interface{}{
		"client_code": api.ClientCode,
		"data":        params,
		"sign":        sign,
		"time":        date,
	}
	requestJson, err := json.Marshal(requestArr)
	if err != nil {
		return
	}
	requestBase64 := base64.StdEncoding.EncodeToString([]byte(string(requestJson)))

	// 创建表单
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("request_data", requestBase64)

	contentType := bodyWriter.FormDataContentType()
	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		return
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var resArr map[string]interface{}
	err = json.Unmarshal(result, &resArr)
	if _, ok := resArr["error_code"]; !ok {
		err = errors.New("返回值错误,无 error_code")
		return
	}
	if _, ok := resArr["message"]; !ok {
		err = errors.New("返回值错误,无 message")
		return
	}
	if _, ok := resArr["return_data"]; !ok {
		err = errors.New("返回值错误,无 return_data")
		return
	}

	res = NewResponse(fmt.Sprint(resArr["error_code"]), fmt.Sprint(resArr["message"]), resArr["return_data"])
	return
}
