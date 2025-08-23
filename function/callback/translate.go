package callback

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gounits/CloudFunctions/tool"
	"hash"
	"io"
	"net/http"
	"strings"
	"time"
)

type ITranslate interface {
	Translate(text string) (result string, err error)
}

type TranslateText struct {
	Text string `json:"text"` // 需要翻译的文本
	From string `json:"from"` // 原始语言
	To   string `json:"to"`   // 翻译成的语言
}

func Translate(tt TranslateText) (result string, err error) {
	translate := NewXFYunTranslate(tt)
	result, err = translate.Translate(tt.Text)
	return
}

type xFYunTranslate struct {
	TranslateText
	id     string // 翻译的APPID
	secret string // 翻译的密钥
	key    string // 翻译的公钥
}

func NewXFYunTranslate(translate TranslateText) ITranslate {
	xfy := tool.Conf.Translates["xfy"]
	return &xFYunTranslate{id: xfy.Appid, secret: xfy.Secret, key: xfy.ApiKey, TranslateText: translate}
}

func (x xFYunTranslate) post(data []byte) (result []byte, err error) {
	var (
		request  *http.Request
		response *http.Response
		link     = "https://itrans.xfyun.cn/v2/its"
		host     = strings.ToLower("iTrans.xfYun.cn")
	)

	client := &http.Client{}

	param := map[string]any{
		"common":   map[string]string{"app_id": x.id},
		"business": map[string]string{"from": x.From, "to": x.To},
		"data":     map[string]string{"text": base64.StdEncoding.EncodeToString(data)},
	}

	tt, _ := json.Marshal(param)

	if request, err = http.NewRequest("POST", link, bytes.NewReader(tt)); err != nil {
		return
	}

	sign := func(hash hash.Hash, data string) string {
		hash.Write([]byte(data))
		encodeData := hash.Sum(nil)
		return base64.StdEncoding.EncodeToString(encodeData)
	}

	signHmac := func(data string, secret string) string {
		mac := hmac.New(sha256.New, []byte(secret))
		return sign(mac, data)
	}

	signBody := func(data string) string {
		sha := sha256.New()
		return sign(sha, data)
	}

	//增加header选项
	{
		date := time.Now().UTC().Format(time.RFC1123)
		digest := "SHA-256=" + signBody(string(tt))

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Host", host)
		request.Header.Set("Accept", "application/json,version=1.0")
		request.Header.Set("Date", date)
		request.Header.Set("Digest", digest)

		// 根据请求头部内容，生成签名
		signature := fmt.Sprintf("host: %s\ndate: %s\nPOST /v2/its HTTP/1.1\ndigest: %s", host, date, digest)
		signature = signHmac(signature, x.secret)

		// 组装Authorization头部
		authHeader := fmt.Sprintf(`api_key="%s", algorithm="hmac-sha256", headers="host date request-line digest", signature="%s"`, x.key, signature)
		request.Header.Set("Authorization", authHeader)
	}

	//处理返回结果
	if response, err = client.Do(request); err != nil {
		return
	}

	defer response.Body.Close()

	result, err = io.ReadAll(response.Body)
	return
}

func (x xFYunTranslate) Translate(text string) (result string, err error) {
	type R struct {
		Code int `json:"code"`
		Data struct {
			Result map[string]any `json:"result"`
		} `json:"data"`
		Message string `json:"message"`
	}

	var (
		r        R
		response []byte
	)

	data := []byte(text)

	if response, err = x.post(data); err != nil {
		return
	}

	if err = json.Unmarshal(response, &r); err != nil {
		return
	}

	if r.Message != "success" {
		err = fmt.Errorf("翻译错误:%s", r.Message)
		return
	}

	dst := r.Data.Result["trans_result"].(map[string]any)["dst"]
	result = dst.(string)
	return
}
