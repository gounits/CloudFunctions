package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gounits/CloudFunctions/function/callback"
	"github.com/gounits/CloudFunctions/tool"
)

func Translate(w http.ResponseWriter, r *http.Request) {
	var (
		params        callback.TranslateParam
		translators   []callback.ITranslator
		oneTranslator callback.ITranslator
		text          string
		err           error
		name          = r.PathValue("name")
	)

	_ = json.NewDecoder(r.Body).Decode(&params)

	// 加入讯飞翻译
	translators = append(translators, callback.NewXFYunTranslate())

	if name != "" {
		tool.Info(fmt.Sprintf("你正在访问: %s 翻译器", name))
		for _, translator := range translators {
			if translator.Name() == name {
				oneTranslator = translator
				break
			}
		}

		if oneTranslator == nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(tool.Error(fmt.Sprintf("选择的翻译器不存在: %s", name)))
			return
		}

		translators = []callback.ITranslator{oneTranslator}
	} else {
		tool.Warn("没有指定具体的翻译器，将执行所有的翻译器！")
	}

	if text, err = callback.Translate(params, translators...); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(tool.Error(fmt.Sprintf("翻译失败原因: %s", err)))
		return
	}

	tool.Success("翻译成功！")
	_, _ = w.Write([]byte(text))
}
