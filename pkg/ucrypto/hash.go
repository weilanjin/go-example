package ucrypto

import (
	"golang.org/x/crypto/bcrypt"
	"lovec.wlj/pkg/uid"
)

var bcryptSalt = uid.RandSeqID(8)

// BcryptHash 明文加密
func BcryptHash(passwd string) (hash, salt string) {
	salt = bcryptSalt()
	bytes, _ := bcrypt.GenerateFromPassword([]byte(passwd+salt), bcrypt.DefaultCost)
	return string(bytes), salt
}

// BcryptVerify 校验密文和明文
func BcryptVerify(hash, salt, passwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd+salt))
	return err == nil
}
