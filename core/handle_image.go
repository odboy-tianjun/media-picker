package core

import (
	"OdMediaPicker/util"
	"OdMediaPicker/vars"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"image"
	_ "image/gif"  // 导入gif支持
	_ "image/jpeg" // 导入jpeg支持
	_ "image/png"  // 导入png支持
	"os"
	"strings"
)

var ignoreImagePathList []string                       // 忽略的文件路径
var readErrorImagePathList []string                    // 读取信息异常的路径
var imagePath2WidthHeightMap = make(map[string]string) // 图片路径和宽高比
var supportImageTypes = []string{
	".bmp",
	".gif",
	".jpg",
	".jpeg",
	".jpe",
	".png",
	".webp",
}

// 水平图片
var horizontalImageList = make(map[string][]string)
var horizontalGifImageList []string

// 垂直图片
var verticalImageList = make(map[string][]string)
var verticalGifImageList []string

// 等比图片
var squareImageList = make(map[string][]string)
var squareGifImageList []string

// psd
var psdImageList []string

func DoHandleImage(rootDir string) {
	total := len(vars.GlobalImagePathList) // 总数
	successCount := 0                      // 成功数
	errorCount := 0                        // 失败数
	ignoreCount := 0                       // 忽略数
	for _, imageFilePath := range vars.GlobalImagePathList {
		suffix := vars.GlobalFilePath2FileExtMap[imageFilePath]
		if isSupportImage(suffix) {
			err, width, height := readImageInfo(imageFilePath)
			if err == nil {
				successCount = successCount + 1
				imagePath2WidthHeightMap[imageFilePath] = fmt.Sprintf("%d-%d", width, height)
				fmt.Printf("=== Image总数: %d, 已读取Info: %d, 成功数: %d, 失败数: %d \n", total, successCount+errorCount+ignoreCount, successCount, errorCount)
			} else {
				errorCount = errorCount + 1
				readErrorImagePathList = append(readErrorImagePathList, imageFilePath)
				fmt.Printf("=== 异常图片: %s \n", imageFilePath)
			}
			continue
		}
		if strings.EqualFold(suffix, ".webp") { // 特殊文件处理, webp为网络常用图片格式
			webpErr, webpWidth, webpHeight := readWebpTypeImage(imageFilePath)
			if webpErr == nil {
				imagePath2WidthHeightMap[imageFilePath] = fmt.Sprintf("%d-%d", webpWidth, webpHeight)
				successCount = successCount + 1
			} else {
				errorCount = errorCount + 1
				fmt.Printf("=== 异常图片: %s \n", imageFilePath)
			}
			continue
		}
		if strings.EqualFold(suffix, ".bmp") { // 特殊文件处理
			bpmErr, bmpWidth, bmpHeight := readBmpInfo(imageFilePath)
			if bpmErr == nil {
				imagePath2WidthHeightMap[imageFilePath] = fmt.Sprintf("%d-%d", bmpWidth, bmpHeight)
				successCount = successCount + 1
			} else {
				errorCount = errorCount + 1
				fmt.Printf("=== 异常图片: %s \n", imageFilePath)
			}
			continue
		}
		if strings.EqualFold(suffix, ".psd") { // 特殊文件处理
			psdImageList = append(psdImageList, imageFilePath)
			successCount = successCount + 1
			continue
		}
		// 其他的直接先忽略吧, 爱改改, 不改拉倒
		ignoreCount = ignoreCount + 1
		ignoreImagePathList = append(ignoreImagePathList, imageFilePath)
	}
	uid := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	if len(psdImageList) > 0 {
		psdImagePath := rootDir + string(os.PathSeparator) + uid + "_图片_PSD"
		if util.CreateDir(psdImagePath) {
			doMoveFileToDir(psdImageList, psdImagePath)
		}
	}
	if len(readErrorImagePathList) > 0 {
		readInfoErrorPath := rootDir + string(os.PathSeparator) + uid + "_图片_读取异常"
		if util.CreateDir(readInfoErrorPath) {
			doMoveFileToDir(readErrorImagePathList, readInfoErrorPath)
		}
	}
	if len(ignoreImagePathList) > 0 {
		ignorePath := rootDir + string(os.PathSeparator) + uid + "_图片_已忽略"
		if util.CreateDir(ignorePath) {
			doMoveFileToDir(ignoreImagePathList, ignorePath)
		}
	}
	doPickImageFile(uid, rootDir, imagePath2WidthHeightMap)
	fmt.Printf("=== 图片处理完毕(UID): %s \n\n", uid)
}

// 条件图片并分组存放
func doPickImageFile(uid string, rootDir string, imagePath2WidthHeightMap map[string]string) {
	if len(imagePath2WidthHeightMap) == 0 {
		fmt.Printf("=== 当前目录下没有扫描到图片文件, %s \n", rootDir)
		return
	}
	for currentImagePath, infoStr := range imagePath2WidthHeightMap {
		width2Height := strings.Split(infoStr, "-")
		width := util.String2int(width2Height[0])
		height := util.String2int(width2Height[1])
		suffix := vars.GlobalFilePath2FileExtMap[currentImagePath]
		if width > height {
			handleHorizontalImage(currentImagePath, width, height, suffix)
			continue
		}
		if width < height {
			handleVerticalImage(currentImagePath, height, suffix)
			continue
		}
		handleSquareImage(currentImagePath, width, height, suffix)
	}
	moveHorizontalImage(rootDir, uid)
	moveVerticalImage(rootDir, uid)
	moveSquareImage(rootDir, uid)
}

func moveSquareImage(rootDir string, uid string) {
	pathSeparator := string(os.PathSeparator)
	squareGifImagePath := rootDir + pathSeparator + uid + "_图片_等比_GIF"
	if len(squareGifImageList) > 0 {
		util.CreateDir(squareGifImagePath)
		doMoveFileToDir(squareGifImageList, squareGifImagePath)
	}
	for widthStr, imagePaths := range horizontalImageList {
		squareImagePath := rootDir + pathSeparator + uid + "_图片_等比_" + widthStr
		util.CreateDir(squareImagePath)
		doMoveFileToDir(imagePaths, squareImagePath)
	}
}

func handleSquareImage(currentImagePath string, width int, height int, suffix string) {
	if strings.EqualFold(suffix, ".gif") {
		squareGifImageList = append(squareGifImageList, currentImagePath)
		return
	}
	widthStr := fmt.Sprintf("%d", width)
	if squareImageList[widthStr] == nil {
		squareImageList[widthStr] = make([]string, 0)
	}
	squareImageList[widthStr] = append(squareImageList[widthStr], currentImagePath)
}

// 移动垂直图片
func moveVerticalImage(rootDir string, uid string) {
	pathSeparator := string(os.PathSeparator)
	verticalGifImagePath := rootDir + pathSeparator + uid + "_图片_竖屏_GIF"
	if len(verticalGifImageList) > 0 {
		util.CreateDir(verticalGifImagePath)
		doMoveFileToDir(verticalGifImageList, verticalGifImagePath)
	}
	for heightStr, imagePaths := range verticalImageList {
		verticalImagePath := rootDir + pathSeparator + uid + "_图片_竖屏_" + heightStr
		util.CreateDir(verticalImagePath)
		doMoveFileToDir(imagePaths, verticalImagePath)
	}
}

// 移动水平图片
func moveHorizontalImage(rootDir string, uid string) {
	pathSeparator := string(os.PathSeparator)
	horizontalGifImagePath := rootDir + pathSeparator + uid + "_图片_横屏_GIF"
	if len(horizontalGifImageList) > 0 {
		util.CreateDir(horizontalGifImagePath)
		doMoveFileToDir(horizontalGifImageList, horizontalGifImagePath)
	}
	for widthStr, imagePaths := range horizontalImageList {
		horizontalImagePath := rootDir + pathSeparator + uid + "_图片_横屏_" + widthStr
		util.CreateDir(horizontalImagePath)
		doMoveFileToDir(imagePaths, horizontalImagePath)
	}
}

// 处理垂直图片
func handleVerticalImage(currentImagePath string, height int, suffix string) {
	if strings.EqualFold(suffix, ".gif") {
		verticalGifImageList = append(verticalGifImageList, currentImagePath)
		return
	}
	heightStr := fmt.Sprintf("%d", height)
	if verticalImageList[heightStr] == nil {
		verticalImageList[heightStr] = make([]string, 0)
	}
	verticalImageList[heightStr] = append(verticalImageList[heightStr], currentImagePath)
}

// 处理横向图片
func handleHorizontalImage(currentImagePath string, width int, height int, suffix string) {
	if strings.EqualFold(suffix, ".gif") {
		horizontalGifImageList = append(horizontalGifImageList, currentImagePath)
		return
	}
	widthStr := fmt.Sprintf("%d", width)
	if horizontalImageList[widthStr] == nil {
		horizontalImageList[widthStr] = make([]string, 0)
	}
	horizontalImageList[widthStr] = append(horizontalImageList[widthStr], currentImagePath)
}

// 判断是否属于支持的图片文件
func isSupportImage(imageType string) bool {
	for _, supportImageType := range supportImageTypes {
		if strings.EqualFold(supportImageType, imageType) {
			return true
		}
	}
	return false
}

// 读取一般图片文件信息
func readImageInfo(filePath string) (err error, width int, height int) {
	file, err := os.Open(filePath) // 图片文件路径
	if err != nil {
		return err, 0, 0
	}
	defer file.Close()
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return err, 0, 0
	}
	return nil, img.Width, img.Height
}
