package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// @Descifrado
// @Descifrado de datos recibidos
// @Tags Descifrado
// @Param encryptedData dato encriptado
// @param key  clave de encriptacion La misma que usaste en PHP u otro language
// @Success 201 {object} models.ingresoResp "OK"
// @Failure 400 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// decryptAES descifra datos cifrados con AES en modo CBC

func DecryptAES(encryptedData, key string) ([]byte, error) {
	// Decodificar el dato cifrado desde base64
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	// Extraer el IV (los primeros 16 bytes)
	iv := data[:aes.BlockSize]
	encrypted := data[aes.BlockSize:]

	// Crear el bloque de cifrado
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	// Crear el stream de descifrado
	mode := cipher.NewCBCDecrypter(block, iv)

	// Descifrar los datos
	mode.CryptBlocks(encrypted, encrypted)

	// Eliminar el padding (si es necesario)
	decrypted := pkcs7Unpad(encrypted)

	return decrypted, nil
}

// pkcs7Unpad elimina el padding PKCS7
func pkcs7Unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

// @cifrado
// @cifrado de datos enviados
// @Tags cifrado
// @Param data dato que quiere cifrar
// @param key  clave de encriptacion La misma que usaste en PHP u otro language
// @Success 201 {object} models.ingresoResp "OK"
// @Failure 400 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// decryptAES descifra datos cifrados con AES en modo CBC

// encryptAES cifra los datos utilizando AES-256-CBC
func EncryptAES(data string, key string) (string, error) {
	// Crear el bloque de cifrado
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Generar un IV aleatorio
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Cifrar los datos
	encrypted := make([]byte, len(data))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted, []byte(data))

	// Concatenar el IV con los datos cifrados y codificar en base64
	encryptedData := base64.StdEncoding.EncodeToString(append(iv, encrypted...))

	return encryptedData, nil
}
