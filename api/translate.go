package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gounits/CloudFunctions/function/callback"
)

func Translate(w http.ResponseWriter, r *http.Request) {
	var (
		translate callback.TranslateText
		text      string
		err       error
	)

	_ = json.NewDecoder(r.Body).Decode(&translate)

	if text, err = callback.Translate(translate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(fmt.Append(nil, err))
		return
	}

	_, _ = w.Write([]byte(text))
}
