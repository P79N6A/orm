package fiorm

import (
	"encoding/hex"

	"git.code.oa.com/fip-team/fiorm/internal"
	"github.com/gtank/cryptopasta"
)

// HashPassword 加密密码，返回加密的数据
func HashPassword(password string) string {
	pass := []byte(password)
	hashed, err := cryptopasta.HashPassword(pass)
	hashKey := string(hashed[:])
	if err != nil {
		panic(err)
	}

	return hashKey
}

// CheckPassword 检查密码
func CheckPassword(hashKey string, password string) bool {
	pass := []byte(password)
	hash := []byte(hashKey)
	if err := cryptopasta.CheckPasswordHash(hash, pass); err != nil {
		return false
	}

	return true
}

// Encrypt 根据Key加密一个字符串，返回加密的数据
func Encrypt(text string) (string, error) {
	k := &[32]byte{}
	key := internal.GetKey()
	copy(k[:], key)

	btext := []byte(text)
	ciphertext, err := cryptopasta.Encrypt(btext, k)
	if err != nil {
		return "", err
	}

	encodeText := hex.EncodeToString(ciphertext)
	return encodeText, nil
}
