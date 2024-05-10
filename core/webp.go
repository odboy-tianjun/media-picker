package core

import (
	"fmt"
	"golang.org/x/image/webp"
	"os"
)

// 读取webp格式的图片信息
func readWebpTypeImage(webpFilePath string) (err error, width int, height int) {
	// 打开WebP文件
	file, err := os.Open(webpFilePath)
	if err != nil {
		fmt.Printf("=== Failed to open file: %v\n", err)
		return err, 0, 0
	}
	defer file.Close()
	// 使用webp.DecodeConfig解码WebP图片配置信息（不加载完整像素数据）
	imgConfig, err := webp.DecodeConfig(file)
	if err != nil {
		fmt.Printf("=== Failed to decode WebP image config: %v\n", err)
		return err, 0, 0
	}
	return nil, imgConfig.Width, imgConfig.Height
}
