package option

import "testing"

func TestOption(t *testing.T) {
	client := NewClient("root", "admin123", "test", WithHost("192.168.0.1"))
	t.Log(client)
}