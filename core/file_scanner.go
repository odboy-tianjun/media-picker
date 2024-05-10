package core

import (
	"OdMediaPicker/util"
	"OdMediaPicker/vars"
	"fmt"
	"github.com/redmask-hb/GoSimplePrint/goPrint"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type FileScanner struct {
}

func (FileScanner) DoScan(rootDir string) {
	fmt.Printf("=== 开始扫描文件 \n")
	if err := filepath.Walk(rootDir, visit); err != nil {
		log.Fatal(err)
	}
	doReadFileMimeInfo()
}

func doReadFileMimeInfo() {
	total := len(vars.GlobalFilePathList)
	fmt.Printf("=== 文件总数: %d \n", total)
	// 扫描文件mime信息
	var count = 0
	bar := goPrint.NewBar(100)
	bar.SetNotice("=== 扫描文件Mime：")
	bar.SetGraph(">")
	bar.SetNoticeColor(goPrint.FontColor.Red)
	for _, currentPath := range vars.GlobalFilePathList {
		// mime
		vars.GlobalFilePath2MimeInfoMap[currentPath] = util.ReadFileMimeInfo(currentPath).String()
		count = count + 1
		bar.PrintBar(util.CalcPercentage(count, total))
	}
	bar.PrintEnd("=== Finish")
}

func (FileScanner) DoFilter() {
	total := len(vars.GlobalFilePathList)
	var count = 0
	bar := goPrint.NewBar(100)
	bar.SetNotice("=== 过滤已支持的媒体：")
	bar.SetGraph(">")
	for _, globalFilePath := range vars.GlobalFilePathList {
		fileMime := vars.GlobalFilePath2MimeInfoMap[globalFilePath]
		count = count + 1
		if strings.Contains(fileMime, "video/") { // 视频
			vars.GlobalVideoPathList = append(vars.GlobalVideoPathList, globalFilePath)
			bar.PrintBar(util.CalcPercentage(count, total))
			continue
		}
		// mime格式为application/octet-stream的视频
		ext := path.Ext(globalFilePath)
		if isSupportVideo(ext) {
			vars.GlobalVideoPathList = append(vars.GlobalVideoPathList, globalFilePath)
			bar.PrintBar(util.CalcPercentage(count, total))
			continue
		}
		if strings.Contains(fileMime, "image/") { // 图片
			vars.GlobalImagePathList = append(vars.GlobalImagePathList, globalFilePath)
			bar.PrintBar(util.CalcPercentage(count, total))
			continue
		}
		if isSupportImage(ext) {
			vars.GlobalImagePathList = append(vars.GlobalImagePathList, globalFilePath)
			bar.PrintBar(util.CalcPercentage(count, total))
			continue
		}
		// 其他的文件不处理
	}
	bar.PrintEnd("=== Finish")
}

// 定义walkFn回调函数visit
func visit(currentPath string, info os.FileInfo, err error) error {
	if err != nil {
		return err // 如果有错误，直接返回
	}
	if !info.IsDir() {
		vars.GlobalFilePathList = append(vars.GlobalFilePathList, currentPath)
		// filename, include ext
		vars.GlobalFilePath2FileNameMap[currentPath] = filepath.Base(currentPath)
		// file ext
		vars.GlobalFilePath2FileExtMap[currentPath] = path.Ext(currentPath)
	}
	return nil
}
