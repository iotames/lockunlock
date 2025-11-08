package lockunlock

import (
	"fmt"
	"os/exec"
	"runtime"
)

// ShowMessage 显示系统弹框消息
func ShowMessage(title, message string) {
	switch runtime.GOOS {
	case "windows":
		// Windows系统使用PowerShell显示弹框
		script := fmt.Sprintf(`Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.MessageBox]::Show('%s', '%s')`, message, title)
		exec.Command("powershell", "-Command", script).Run()
	case "darwin":
		// macOS系统使用osascript显示弹框
		script := fmt.Sprintf(`display dialog "%s" with title "%s" buttons {"OK"}`, message, title)
		exec.Command("osascript", "-e", script).Run()
	default:
		// Linux系统使用zenity显示弹框（如果已安装）
		_, err := exec.LookPath("zenity")
		if err != nil {
			// 如果没有安装zenity，回退到终端输出
			fmt.Printf("%s: %s\n", title, message)
			return
		}
		exec.Command("zenity", "--info", "--title", title, "--text", message).Run()
	}
}

func DevDialog() {
	// TODO
	// 帮我写一个跨平台的dialog对话框组件，业务场景是对指定目录的文件进行加密或解密。
	// 功能如下：
	// 	标题：lockunlock文件密码锁
	// 	1个密码输入框：label=密钥，placeholder="选填。输入32位密钥"。默认空字符串。警告：丢失密钥会无法解锁。
	// 	1个目录选择器：label=选择目录，可选。不选择默认使用当前目录。
	// 	描述: 默认加密当前目录所有文件，不包括子目录。
	// 	最下方两个按钮：左边是加密，右边是解密。
	// 要求：
	// 1. 尽可能做轻量级组件。调用操作系统底层API的，不引用Go语言第三方库，第三方依赖越少越好。
	// 2. 代码解耦，方便元素组装，函数拆分为多个文件，文件之间平级，在一个目录下。
	// 3. 由于调用系统API，故不需要定制过于花哨的GUI，只要实现基本功能即可。
	// 4. 只要做工具函数，不要做加解密的业务逻辑。
	ShowMessage("TODO 功能开发中", "TODO 功能开发中......")
}
