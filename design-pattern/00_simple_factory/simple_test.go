package simplefactory

import "testing"

func TestTencentAPI(t *testing.T) {
	api := NewAPI("tencent")
	s := api.Speech([]byte("正确率94%"))
	if s != "[Tencent] 正确率94%" {
		t.Fatal("tencent api test fail")
	}
}

func TestIflytekAPI(t *testing.T) {
	api := NewAPI("iflytek")
	s := api.Speech([]byte("正确率94%"))
	if s != "[Iflytek] 正确率94%" {
		t.Fatal("iflytek api test fail")
	}
}
