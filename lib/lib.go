package lib

import (
	"crypto/aes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
)

// generateKey
// @param key
// @date 2022-09-21 16:58:57
func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

// AesEncryptECB
// @param src
// @param key
// @date 2022-09-21 16:58:56
func AesEncryptECB(src []byte, key []byte) ([]byte, error) {
	key, err := AesSha1prng(key, 128) // 比示例一多出这一步
	if err != nil {
		return nil, err
	}

	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(src) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, src)
	pad := byte(len(plain) - len(src))
	for i := len(src); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted := make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(src); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted, nil
}

// AesSha1prng
// @param keyBytes
// @param encryptLength
// @date 2022-09-21 16:58:55
func AesSha1prng(keyBytes []byte, encryptLength int) ([]byte, error) {
	hashs := Sha1(Sha1(keyBytes))
	maxLen := len(hashs)
	realLen := encryptLength / 8
	if realLen > maxLen {
		return nil, errors.New("invalid length!")
	}

	return hashs[0:realLen], nil
}

// Sha1
// @param data
// @date 2022-09-21 16:58:54
func Sha1(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

// AesDecryptECB
// @param encrypted
// @param key
// @date 2022-09-21 16:58:53
func AesDecryptECB(encrypted []byte, key []byte) ([]byte, error) {
	key, err := AesSha1prng(key, 128) // 比示例一多出这一步
	if err != nil {
		return nil, err
	}

	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted := make([]byte, len(encrypted))
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim], nil
}

// AesDecryptECBFromHex
// @param encrypted
// @param key
// @date 2022-09-21 16:58:52
func AesDecryptECBFromHex(encrypted string, key []byte) ([]byte, error) {
	param, err := hex.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

	return AesDecryptECB(param, key)
}
