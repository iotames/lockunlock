package lockunlock

import (
	"fmt"
	"os"
	"path/filepath"
)

func LockDirFiles(skey []byte, dir string) error {
	return diropt("加密", dir, func(fname string) error { return EncryptFile(fname, fname, skey) })
}

func UnlockDirFiles(skey []byte, dir string) error {
	return diropt("解密", dir, func(fname string) error { return DecryptFile(fname, fname, skey) })
}

func diropt(optname, dir string, eachopt func(fname string) error) error {
	// 获取dir目录下的所有文件（不包括子目录）
	if dir == "" {
		dir = "."
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("读取目录失败：%v", err)
	}

	// 获取当前正在运行的可执行文件名
	executable, err := os.Executable()
	if err != nil {
		fmt.Printf(" %v\n", err)
		return fmt.Errorf("获取可执行文件名失败: %v", err)
	}

	execName := ""
	if fileInfo, err := os.Stat(executable); err == nil {
		execName = fileInfo.Name()
	}
	fmt.Println("当前文件名：", execName, "os.Args[0]", os.Args[0])
	count := 0 // 记录处理的文件数量

	// 遍历所有条目
	for _, entry := range entries {
		// 只处理文件，跳过目录
		if entry.IsDir() {
			continue
		}

		// 获取文件名
		filename := entry.Name()

		// 跳过正在运行的可执行文件（避免加密自身）
		// filename == execName 替代 filename == os.Args[0]
		if filename == execName || filename == "." || filename == ".." {
			continue
		}
		if dir != "." {
			filename = filepath.Join(dir, filename)
		}
		if err := eachopt(filename); err != nil {
			// 加密文件 encrypt.exe 失败: open encrypt.exe: The process cannot access the file because it is being used by another process.
			fmt.Printf("%s文件 %s 失败: %v\n", optname, filename, err)
			continue
		}

		count++
		fmt.Printf("文件%s成功: %s\n", optname, filename)
	}
	// 显示弹框提示
	ShowMessage(fmt.Sprintf("%s完成", optname), fmt.Sprintf("总共%s了 %d 个文件", optname, count))
	return nil
}
