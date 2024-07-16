package main

import (
	"OdMediaPicker/core"
	"OdMediaPicker/vars"
	_ "embed"
	"fmt"
	_ "image/gif"  // 导入gif支持
	_ "image/jpeg" // 导入jpeg支持
	_ "image/png"  // 导入png支持
	"os"
	"time"
)

func main() {
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Println("=== 获取当前路径异常", err)
		return
	}
	scanner := core.FileScanner{}
	scanner.DoScan(rootDir)
	scanner.DoFilter()
	if len(vars.GlobalImagePathList) > 0 {
		core.DoHandleImage(rootDir)
	}
	if len(vars.GlobalVideoPathList) > 0 {
		core.DoHandleVideo(rootDir)
	}
	fmt.Println("=== 5s后自动退出")
	time.Sleep(time.Second * 5)
}
