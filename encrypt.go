package lockunlock

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"io"
	"os"
)

// PKCS7填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// 插入随机字节：每隔37个字节插入17个随机字节
func insertRandomBytes(data []byte) ([]byte, error) {
	blockSize := 37
	insertSize := 17
	var result []byte

	for i := 0; i < len(data); i += blockSize {
		end := i + blockSize
		if end > len(data) {
			end = len(data)
		}

		// 添加当前块
		result = append(result, data[i:end]...)

		// 如果不是最后一块且块大小足够，插入随机字节
		if end < len(data) && end-i == blockSize {
			randomBytes := make([]byte, insertSize)
			if _, err := rand.Read(randomBytes); err != nil {
				return nil, err
			}
			result = append(result, randomBytes...)
		}
	}

	return result, nil
}

// AES加密函数
func aesEncrypt(plainText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// PKCS7填充
	plainText = pkcs7Padding(plainText, block.BlockSize())

	// 生成随机IV
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// CBC加密
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	return cipherText, nil
}

// 加密文件
func EncryptFile(inputPath, outputPath string, key []byte) error {
	// 读取原始文件
	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	// 1. 字节反转
	reversed := ReverseBytes(inputData)

	// 2. 插入随机字节
	withRandomBytes, err := insertRandomBytes(reversed)
	if err != nil {
		return err
	}

	// 3. AES加密
	encryptedData, err := aesEncrypt(withRandomBytes, key)
	if err != nil {
		return err
	}

	// 写入加密文件
	return os.WriteFile(outputPath, encryptedData, 0644)
}
