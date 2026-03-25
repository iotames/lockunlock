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

// DecryptFile 解密文件
//
// 解密流程：
//  1. AES-256 解密
//  2. 可选混淆后处理：若 obfuscate 为 true，执行以下混淆恢复步骤：
//     a. 删除混淆时插入的随机字节（每隔 37 个字节删除 17 个字节）
//     b. 将字节数组首尾颠倒（恢复原始顺序）
//  3. 将结果写入输出文件
//
// 注意：本函数的 obfuscate 参数必须与加密时使用的 obfuscate 参数一致，
// 否则解密将失败或得到错误数据。
//
// 参数：
//   - inputPath  string  输入文件路径（加密文件）
//   - outputPath string  输出文件路径（解密文件）
//   - key        []byte  AES-256 密钥，长度必须为 32 字节，与加密密钥相同
//   - obfuscate  bool    是否启用混淆后处理恢复：
//     true  = 执行混淆恢复步骤
//     false = 直接输出 AES 解密结果
//
// 返回值：
//   - error 解密过程中出现的任何错误，成功时返回 nil
func DecryptFile(inputPath, outputPath string, key []byte, obfuscate bool) error {
	var encryptedData, decryptedData, dataToRestore, originalData []byte
	var err error

	// 读取加密文件
	encryptedData, err = os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	// 1. AES解密
	decryptedData, err = aesDecrypt(encryptedData, key)
	if err != nil {
		return err
	}

	// 混淆后处理恢复
	if obfuscate {
		// 2. 移除随机字节
		dataToRestore = removeRandomBytes(decryptedData)
		if err != nil {
			return err
		}
		// 3. 字节反转（恢复原始顺序）
		originalData = ReverseBytes(dataToRestore)
	} else {
		// 未启用混淆，直接使用解密数据
		originalData = decryptedData
	}

	// 写入解密文件
	return os.WriteFile(outputPath, originalData, 0644)
}
