// Функции шифрование и дешифрования информации
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	mrand "math/rand"
)

// Преобразование байт в строку
func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// Преобразование строки в байты
func Decode(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Шифрует строку
func Encrypt(text, MySecret string) (string, error) {
	plainText := []byte(text)

	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	return Encode(cipherText), nil
}

// Расшифровывает строку
func Decrypt(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	cipherText, err := Decode(text)
	if err != nil {
		return "", err
	}
	if len(cipherText) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

// Генерирует случайную строку с длинной n символов
func RandStr(n int) string {
	var lower = []byte("abcdefghijklmnopqrstuvwxyz")
	var upper = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var number = []byte("0123456789")
	var special = []byte("~=+%^*/()[]{}/!@#$?|")
	allChar := append(lower, upper...)
	allChar = append(allChar, number...)
	allChar = append(allChar, special...)

	b := make([]byte, n)
	// select 1 upper, 1 lower, 1 number and 1 special
	b[0] = lower[mrand.Intn(len(lower))]
	b[1] = upper[mrand.Intn(len(upper))]
	b[2] = number[mrand.Intn(len(number))]
	b[3] = special[mrand.Intn(len(special))]
	for i := 4; i < n; i++ {
		// randomly select 1 character from given charset
		b[i] = allChar[mrand.Intn(len(allChar))]
	}

	//shuffle character
	mrand.Shuffle(len(b), func(i, j int) {
		b[i], b[j] = b[j], b[i]
	})

	return string(b)
}
