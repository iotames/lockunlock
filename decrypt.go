package lockunlock

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
)

// PKCS7去除填充
func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("数据长度为零")
	}
	padding := int(data[length-1])
	if padding > length {
		return nil, fmt.Errorf("填充大小无效")
	}
	return data[:length-padding], nil
}

// 移除插入的随机字节
func removeRandomBytes(data []byte) []byte {
	blockSize := 37
	insertSize := 17
	totalBlockSize := blockSize + insertSize
	var result []byte

	i := 0
	for i < len(data) {
		end := i + blockSize
		if end > len(data) {
			end = len(data)
		}

		// 添加有效数据块
		result = append(result, data[i:end]...)

		// 跳过插入的随机字节（如果存在）
		i += totalBlockSize
		if i > len(data) {
			break
		}
	}

	return result
}

// AES解密函数
func aesDecrypt(cipherText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		return nil, fmt.Errorf("密文太短")
	}

	// 提取IV
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	// CBC解密
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	// 去除PKCS7填充
	return pkcs7Unpadding(cipherText)
}

// 解密文件
func DecryptFile(inputPath, outputPath string, key []byte) error {
	// 读取加密文件
	encryptedData, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	// 1. AES解密
	decryptedData, err := aesDecrypt(encryptedData, key)
	if err != nil {
		return err
	}

	// 2. 移除随机字节
	withoutRandomBytes := removeRandomBytes(decryptedData)

	// 3. 字节反转（恢复原始顺序）
	// 反转字节切片（与加密时相同）
	originalData := ReverseBytes(withoutRandomBytes)

	// 写入解密文件
	return os.WriteFile(outputPath, originalData, 0644)
}
