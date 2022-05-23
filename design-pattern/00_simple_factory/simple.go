package simplefactory

// 简单工厂模式 一般用 NewXx(), 暴露出当前实例
import (
	"fmt"
	"log"
)

type API interface {
	// 语音听写
	Speech(b []byte) string
}

func NewAPI(vendor string) API {
	switch vendor {
	case "tencent":
		return &TencentAPI{}
	case "iflytek":
		return &IflytekAPI{}
	default:
		log.Println(vendor, "not mod")
		return nil
	}
}

type TencentAPI struct{}

func (*TencentAPI) Speech(b []byte) string {
	return fmt.Sprintf("[Tencent] %s", b)
}

type IflytekAPI struct{}

func (*IflytekAPI) Speech(b []byte) string {
	return fmt.Sprintf("[Iflytek] %s", b)
}
