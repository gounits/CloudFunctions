package callback_test

import (
	"fmt"
	"testing"

	"github.com/gounits/CloudFunctions/function/callback"
)

func TestXFYunTranslate_Translate(t *testing.T) {
	xf := callback.XFYunTranslate{
		ID:     "",
		Secret: "",
		Key:    "",
		From:   "",
		To:     "",
		Text:   "",
	}

	if result, err := xf.Translate(); err != nil {
		t.Error(err)
	} else {
		fmt.Println(result)
	}
}
