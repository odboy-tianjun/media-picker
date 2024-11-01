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
	".psb",
}

// gif图片
var horizontalGifImageList []string
var verticalGifImageList []string
var squareGifImageList []string

// 标准横向图片
var horizontalStandard720PImageList []string
var horizontalStandard1080PImageList []string
var horizontalStandard4KImageList []string
var horizontalStandard8KImageList []string

// psd图片
var psdImageList []string

// 规格对应图片所在的路径
var hSpecs2ImageList = make(map[string][]string)
var vSpecs2ImageList = make(map[string][]string)
var mSpecs2ImageList = make(map[string][]string)

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
		psdImagePath := rootDir + string(os.PathSeparator) + uid + "-图片_PSD"
		if util.CreateDir(psdImagePath) {
			doMoveFileToDir(psdImageList, psdImagePath)
		}
	}
	if len(readErrorImagePathList) > 0 {
		readInfoErrorPath := rootDir + string(os.PathSeparator) + uid + "-图片_读取异常"
		if util.CreateDir(readInfoErrorPath) {
			doMoveFileToDir(readErrorImagePathList, readInfoErrorPath)
		}
	}
	if len(ignoreImagePathList) > 0 {
		ignorePath := rootDir + string(os.PathSeparator) + uid + "-图片_已忽略"
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
			handleVerticalImage(currentImagePath, width, height, suffix)
			continue
		}
		handleSquareImage(currentImagePath, width, height, suffix)
	}
	moveStandardHorizontalImage(rootDir, uid)
	moveStandardVerticalImage(rootDir, uid)
	moveStandardSquareImage(rootDir, uid)
	// 移动特殊规格的
	moveSpecsHorizontalImage(rootDir, uid)
	moveSpecsVerticalImage(rootDir, uid)
	moveSpecsSquareImage(rootDir, uid)
}

func moveStandardSquareImage(rootDir string, uid string) {
	pathSeparator := string(os.PathSeparator)
	squareGifImagePath := rootDir + pathSeparator + uid + "-图片_等比_GIF"
	if len(squareGifImageList) > 0 {
		util.CreateDir(squareGifImagePath)
		doMoveFileToDir(squareGifImageList, squareGifImagePath)
	}
}
func moveSpecsSquareImage(rootDir string, uid string) {
	pathSeparator := string(os.PathSeparator)
	for dirName := range mSpecs2ImageList {
		imagePath := rootDir + pathSeparator + uid + "-" + dirName
		if len(mSpecs2ImageList[dirName]) > 0 {
			util.CreateDir(imagePath)
			doMoveFileToDir(mSpecs2ImageList[dirName], imagePath)
		}
	}
}
func handleSquareImage(currentImagePath string, width int, height int, suffix string) {
	if strings.EqualFold(suffix, ".gif") {
		squareGifImageList = append(squareGifImageList, currentImagePath)
		return
	}
	// 非标准规格图片处理
	prefix := fmt.Sprintf("M_%d_%d", width, height)
	if mSpecs2ImageList[prefix] == nil {
		mSpecs2ImageList[prefix] = make([]string, 1)
	}
	mSpecs2ImageList[prefix] = append(mSpecs2ImageList[prefix], currentImagePath)
}

// 移动垂直图片
func moveStandardVerticalImage(rootDir string, uid string) {
	pathSeparator := string(os.PathSeparator)
	verticalGifImagePath := rootDir + pathSeparator + uid + "-图片_竖屏_GIF"
	if len(verticalGifImageList) > 0 {
		util.CreateDir(verticalGifImagePath)
		doMoveFileToDir(verticalGifImageList, verticalGifImagePath)
	}
}
func moveSpecsVerticalImage(rootDir string, uid string) {
	pathSeparator := string(os.PathSeparator)
	for dirName := range vSpecs2ImageList {
		imagePath := rootDir + pathSeparator + uid + "-" + dirName
		if len(vSpecs2ImageList[dirName]) > 0 {
			util.CreateDir(imagePath)
			doMoveFileToDir(vSpecs2ImageList[dirName], imagePath)
		}
	}
}

// 移动水平图片
func moveStandardHorizontalImage(rootDir string, uid string) {
	pathSeparator := string(os.PathSeparator)
	horizontalGifImagePath := rootDir + pathSeparator + uid + "-图片_横屏_GIF"
	horizontalStandard720PImagePath := rootDir + pathSeparator + uid + "-图片_横屏_720P"
	horizontalStandard1080PImagePath := rootDir + pathSeparator + uid + "-图片_横屏_1080P"
	horizontalStandard4KImagePath := rootDir + pathSeparator + uid + "-图片_横屏_4KP"
	horizontalStandard8KImagePath := rootDir + pathSeparator + uid + "-图片_横屏_8KP"
	if len(horizontalGifImageList) > 0 {
		util.CreateDir(horizontalGifImagePath)
		doMoveFileToDir(horizontalGifImageList, horizontalGifImagePath)
	}
	if len(horizontalStandard720PImageList) > 0 {
		util.CreateDir(horizontalStandard720PImagePath)
		doMoveFileToDir(horizontalStandard720PImageList, horizontalStandard720PImagePath)
	}
	if len(horizontalStandard1080PImageList) > 0 {
		util.CreateDir(horizontalStandard1080PImagePath)
		doMoveFileToDir(horizontalStandard1080PImageList, horizontalStandard1080PImagePath)
	}
	if len(horizontalStandard4KImageList) > 0 {
		util.CreateDir(horizontalStandard4KImagePath)
		doMoveFileToDir(horizontalStandard4KImageList, horizontalStandard4KImagePath)
	}
	if len(horizontalStandard8KImageList) > 0 {
		util.CreateDir(horizontalStandard8KImagePath)
		doMoveFileToDir(horizontalStandard8KImageList, horizontalStandard8KImagePath)
	}
}
func moveSpecsHorizontalImage(rootDir string, uid string) {
	pathSeparator := string(os.PathSeparator)
	for dirName := range hSpecs2ImageList {
		imagePath := rootDir + pathSeparator + uid + "-" + dirName
		if len(hSpecs2ImageList[dirName]) > 0 {
			util.CreateDir(imagePath)
			doMoveFileToDir(hSpecs2ImageList[dirName], imagePath)
		}
	}
}

// 处理垂直图片
func handleVerticalImage(currentImagePath string, width int, height int, suffix string) {
	if strings.EqualFold(suffix, ".gif") {
		verticalGifImageList = append(verticalGifImageList, currentImagePath)
		return
	}
	// 非标准规格图片处理
	prefix := fmt.Sprintf("V_%d_%d", width, height)
	if vSpecs2ImageList[prefix] == nil {
		vSpecs2ImageList[prefix] = make([]string, 1)
	}
	vSpecs2ImageList[prefix] = append(vSpecs2ImageList[prefix], currentImagePath)
}

// 处理横向图片
func handleHorizontalImage(currentImagePath string, width int, height int, suffix string) {
	if strings.EqualFold(suffix, ".gif") {
		horizontalGifImageList = append(horizontalGifImageList, currentImagePath)
		return
	}
	if width >= 1000 && width < 2000 {
		// 1280 x 720 -> 720p
		if width == 1280 && height == 720 {
			horizontalStandard720PImageList = append(horizontalStandard720PImageList, currentImagePath)
			return
		}
		// 1920 x 1080 -> 1080p
		if width == 1920 && height == 1080 {
			horizontalStandard1080PImageList = append(horizontalStandard1080PImageList, currentImagePath)
			return
		}
	} else if width >= 3000 && width < 4000 {
		// 3840 x 2160 -> 4k
		if width == 3840 && height == 2160 {
			horizontalStandard4KImageList = append(horizontalStandard4KImageList, currentImagePath)
			return
		}
	} else if width >= 7000 && width < 8000 {
		// 7680 x 4320 -> 8k
		if width == 7680 && height == 4320 {
			horizontalStandard8KImageList = append(horizontalStandard8KImageList, currentImagePath)
			return
		}
	}
	// 非标准规格图片处理
	prefix := fmt.Sprintf("H_%d_%d", width, height)
	if hSpecs2ImageList[prefix] == nil {
		hSpecs2ImageList[prefix] = make([]string, 1)
	}
	hSpecs2ImageList[prefix] = append(hSpecs2ImageList[prefix], currentImagePath)
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
