// Тестирование шифрования и дешифрования
package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncrypt(t *testing.T) {

	stringToEncrypt := "123"
	masterSK := "qwertyuioasdfghj"

	encText, err := Encrypt(stringToEncrypt, masterSK)
	require.NoError(t, err)

	decText, err := Decrypt(encText, masterSK)
	require.NoError(t, err)

	assert.Equal(t, stringToEncrypt, decText)

}
