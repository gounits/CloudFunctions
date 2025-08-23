package callback_test

import (
	"fmt"
	"github.com/gounits/CloudFunctions/function/callback"
	"testing"
)

func TestXFYunTranslate_Translate(t *testing.T) {
	result, err := callback.Translate(callback.TranslateText{From: "en", To: "cn", Text: "hello world"})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result)
}
