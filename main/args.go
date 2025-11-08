package main

import (
	"flag"
)

var version, dev bool
var opt, key, dir string

func parseCmdArgs() {
	flag.BoolVar(&version, "version", false, "显示版本信息")
	flag.StringVar(&opt, "opt", "", "操作类型: lock, unlock")
	flag.StringVar(&key, "key", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "密钥")
	flag.StringVar(&dir, "dir", "", "指定目录")
	flag.BoolVar(&dev, "dev", false, "开发调试")
	flag.Parse()
}
