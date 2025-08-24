package api

import (
	"encoding/json"
	"fmt"
	"github.com/gounits/CloudFunctions/tool"
	"net/http"

	"github.com/gounits/CloudFunctions/function/callback"
)

func Translate(w http.ResponseWriter, r *http.Request) {
	var (
		params      callback.TranslateParam
		translators []callback.ITranslator
		text        string
		err         error
		name        = r.PathValue("name")
	)

	if name != "" {
		tool.Info(fmt.Sprintf("你正在访问: %s 翻译器", name))
	}

	_ = json.NewDecoder(r.Body).Decode(&params)

	// 加入讯飞翻译
	translators = append(translators, callback.NewXFYunTranslate(params))

	if text, err = callback.Translate(params.Text, translators...); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		tool.Error(fmt.Sprintf("翻译失败原因: %s", err))
		_, _ = w.Write(fmt.Append(nil, err))
		return
	}

	tool.Success("翻译成功！")
	_, _ = w.Write([]byte(text))
}
