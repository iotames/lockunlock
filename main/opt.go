package main

import (
	"github.com/iotames/lockunlock"
)

func lockopt() error {
	var err error
	skey := []byte(key)
	switch opt {
	case "lock":
		err = lockunlock.LockDirFiles(skey, dir)
	case "unlock":
		err = lockunlock.UnlockDirFiles(skey, dir)
	default:
		// err = fmt.Errorf("opt错误: 未知的操作类型:%s", opt)
		err = showLockUnlockDialog(skey, dir)
	}
	return err
}
