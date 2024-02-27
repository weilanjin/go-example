package ucrypto_test

import (
	"lovec.wlj/pkg/ucrypto"
	"testing"
)

func TestAes(t *testing.T) {
	key := "1234567890123456"
	raw := "hello world"
	t.Log(raw)
	ciphertext, err := ucrypto.AesEncrypt(raw, key)
	if err != nil {
		t.Error(err)
	}
	t.Log(ciphertext)
	plaintext, err := ucrypto.AesDecrypt(ciphertext, key)
	if err != nil {
		t.Error(err)
	}
	t.Log(plaintext)
}
