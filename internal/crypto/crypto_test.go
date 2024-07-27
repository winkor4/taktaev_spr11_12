// Тестирование шифрования и дешифрования
package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тестирование шифровки строки и обратного получения
func TestEncrypt(t *testing.T) {

	stringToEncrypt := "123"
	masterSK := "qwertyuioasdfghj"

	encText, err := Encrypt(stringToEncrypt, masterSK)
	require.NoError(t, err)

	decText, err := Decrypt(encText, masterSK)
	require.NoError(t, err)

	assert.Equal(t, stringToEncrypt, decText)

}

// Тестирование алгоритма шифрования на сервере
func TestEncryptOnServer(t *testing.T) {

	masterSK := RandStr(16)
	data := make([]string, 3)
	data[0] = "first string"
	data[1] = "second string"
	data[2] = "third string"

	serverSK := RandStr(16)
	encryptionSK, err := Encrypt(serverSK, masterSK)
	require.NoError(t, err)

	severData := make(map[string]string, 0)
	for _, v := range data {
		dataSK := RandStr(16)
		encData, err := Encrypt(v, dataSK)
		require.NoError(t, err)
		encSK, err := Encrypt(dataSK, serverSK)
		require.NoError(t, err)
		severData[encSK] = encData
	}

	decData := make([]string, 0, 3)
	decServerSK, err := Decrypt(encryptionSK, masterSK)
	require.NoError(t, err)
	for k, v := range severData {
		decDataSK, err := Decrypt(k, decServerSK)
		require.NoError(t, err)
		decV, err := Decrypt(v, decDataSK)
		require.NoError(t, err)
		decData = append(decData, decV)
	}

	for i := 0; i < 3; i++ {
		assert.Equal(t, data[i], decData[i])
	}

}
