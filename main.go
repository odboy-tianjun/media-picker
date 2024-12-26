package main

import (
	"OdMediaPicker/core"
	_ "embed"
	"fmt"
	_ "image/gif"  // 导入gif支持
	_ "image/jpeg" // 导入jpeg支持
	_ "image/png"  // 导入png支持
	"time"
)

func main() {
	rootDir := "E:\\DGD\\待整理\\hent_pic"
	//rootDir, err := os.Getwd()
	//if err != nil {
	//	fmt.Println("=== 获取当前路径异常", err)
	//	return
	//}
	scanner := core.FileScanner{}
	//scanner.DoScan(rootDir)
	//scanner.DoFilter()
	// 整理图片并分组
	//if len(vars.GlobalImagePathList) > 0 {
	//	core.DoHandleImage(rootDir)
	//}
	// 挑选文件数大于N的文件夹并转移
	scanner.DoPickerDir(rootDir, 10)
	//if len(vars.GlobalVideoPathList) > 0 {
	//	core.DoHandleVideo(rootDir)
	//}
	fmt.Println("=== 3s后自动退出")
	time.Sleep(time.Second * 3)
}
