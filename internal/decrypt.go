package internal

import "encoding/hex"
import "github.com/gtank/cryptopasta"

var key = "101232131012312aa34353463&^^%$#312&sdfsdf"

// DbSetting 数据库连接字符串属性
var DbSetting struct {
	Dialect  string
	DbName   string
	Host     string
	User     string
	Password string
	Port     int
}

// GetKey 获取Key值
func GetKey() string {
	return key
}

// Decrypt  解密
func Decrypt(text string) (string, error) {
	k := &[32]byte{}
	copy(k[:], key)

	btext, _ := hex.DecodeString(text)
	ciphertext, err := cryptopasta.Decrypt(btext, k)
	if err != nil {
		return "", err
	}

	encodeText := string(ciphertext[:])
	return encodeText, nil
}
