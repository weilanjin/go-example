// Package ucrypto
// 随机生成16位字符 https://link.fobshanghai.com/tools/xpassword.php
//
// 1.base64  任意二进制 -> 文本的编码
//
// a.编码需要 64个字符表 base64.StdEncoding [+/] base64.URLEncoding [-_]
// b.大小写+数字+[+/ | -_] = 26 + 26 + 10 + 2 = 64个字符表
// c.每三个字节共24位作为一个处理单元,再分为四组,每组6位,查表
//
//	use
//	    base64.StdEncoding.EncodeToString([]byte(plaintext))
//		base64.StdEncoding.DecodeString(ciphertext)
//
// 常用于: URL,Cookie,网页中传输少量二进制数据.
package ucrypto
