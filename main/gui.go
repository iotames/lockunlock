package main

import (
	"fmt"
	"os"

	"github.com/iotames/lockunlock"
	"github.com/sqweek/dialog"
)

// showLockUnlockDialog 显示文件锁对话框并返回用户选择
func showLockUnlockDialog(skey []byte, dir string) error {
	title := "Lockunlock文件密码锁"
	msg := fmt.Sprintf("加解密【目录(%s)】的所有文件，不包含子目录。", dir)
	if dir == "" {
		dir, _ = os.Getwd()
		msg = fmt.Sprintf("默认加解密【当前目录(%s)】的所有文件，不包含子目录。", dir)
	}

	// 显示描述信息
	// dialog.Message(msg).Title(title).Info()

	// // 1. 选择目录
	// directory, err := chooseDirectory()
	// if err != nil {
	// 	return nil, err
	// }

	// // 2. 输入密钥
	// key, err := getKeyInput()
	// if err != nil {
	// 	return nil, err
	// }

	// 3. 选择操作
	result := dialog.Message("%s\n\n点击【是(Y)】加密，【否(N)】解密", msg).Title(title).YesNo()
	if result {
		return lockunlock.LockDirFiles(skey, dir)
	}
	return lockunlock.UnlockDirFiles(skey, dir)
}
