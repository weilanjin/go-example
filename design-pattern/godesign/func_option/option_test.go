package option_test

import (
	option "lovec.wlj/design-patten/godesign/func_option"
	"testing"
)

func TestOption(t *testing.T) {
	client := option.NewClient("root", "admin123", "test", option.WithHost("192.168.0.1"))
	t.Log(client)
}