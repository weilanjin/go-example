package ucrypto_test

import (
	"encoding/base64"
	"lovec.wlj/pkg/ucrypto"
	"testing"
)

func TestBase64(t *testing.T) {
	raw := "hello world"
	t.Log(raw)
	ciphertext := base64.StdEncoding.EncodeToString([]byte(raw))
	t.Log(ciphertext)
	plaintext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s", plaintext)
}

func TestAES(t *testing.T) {
	key := "1234567890123456" // 16、24、32 位
	raw := "hello world"
	t.Log(raw)
	ciphertext, err := ucrypto.AESEncrypt(raw, key)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ciphertext)
	plaintext, err := ucrypto.AESDecrypt(ciphertext, key)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(plaintext)
}

func TestDES(t *testing.T) {
	key := "12345678" // 8 位
	raw := "hello world"
	t.Log(raw)
	ciphertext, err := ucrypto.DESEncrypt(raw, key)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ciphertext)
	plaintext, err := ucrypto.DESDecrypt(ciphertext, key)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(plaintext)
}

func TestMD5(t *testing.T) {
	data := "hello world"
	t.Log(data)
	md5 := ucrypto.MD5(data)
	t.Log(md5)
}

func TestSHA1(t *testing.T) {
	data := "hello world"
	t.Log(data)
	sha1 := ucrypto.SHA1(data)
	t.Log(sha1)
}

func TestSHA256(t *testing.T) {
	secret, data := "w3xeayw5smcn5ei0", "hello world"
	t.Logf("Secret: %s Data: %s\n", secret, data)
	sha := ucrypto.SHA256(secret, data)
	t.Log(sha)
}
