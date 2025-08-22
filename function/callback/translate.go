package callback

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"net/http"
	"strings"
	"time"
)

type XFYunTranslate struct {
	ID     string `json:"id"`     // 翻译的APPID
	Secret string `json:"secret"` // 翻译的密钥
	Key    string `json:"key"`    // 翻译的公钥
	Text   string `json:"text"`   // 需要翻译的文本
	From   string `json:"from"`   // 原始语言
	To     string `json:"to"`     // 翻译成的语言
}

func (x XFYunTranslate) post(data []byte) (result []byte, err error) {
	var (
		request  *http.Request
		response *http.Response
		link     = "https://itrans.xfyun.cn/v2/its"
		host     = strings.ToLower("iTrans.xfYun.cn")
	)

	client := &http.Client{}

	param := map[string]any{
		"common":   map[string]string{"app_id": x.ID},
		"business": map[string]string{"from": x.From, "to": x.To},
		"data":     map[string]string{"text": base64.StdEncoding.EncodeToString(data)},
	}

	tt, _ := json.Marshal(param)

	if request, err = http.NewRequest("POST", link, bytes.NewReader(tt)); err != nil {
		return
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
		signature = signHmac(signature, x.Secret)

		// 组装Authorization头部
		authHeader := fmt.Sprintf(`api_key="%s", algorithm="hmac-sha256", headers="host date request-line digest", signature="%s"`, x.Key, signature)
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

func (x XFYunTranslate) Translate() (text string, err error) {
	type R struct {
		Code int `json:"code"`
		Data struct {
			Result map[string]any `json:"result"`
		} `json:"data"`
		Message string `json:"message"`
	}

	var (
		result   R
		response []byte
	)

	data := []byte(x.Text)

	if response, err = x.post(data); err != nil {
		return
	}

	if err = json.Unmarshal(response, &result); err != nil {
		return
	}

	if result.Message != "success" {
		err = fmt.Errorf("翻译错误:%s", result.Message)
		return
	}

	dst := result.Data.Result["trans_result"].(map[string]any)["dst"]
	text = dst.(string)
	return
}

func sign(hash hash.Hash, data string) string {
	hash.Write([]byte(data))
	encodeData := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func signHmac(data string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	return sign(mac, data)
}

func signBody(data string) string {
	sha := sha256.New()
	return sign(sha, data)
}
