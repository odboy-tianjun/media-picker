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
var normalImageList []string
var horizontalUnHandleImageList []string
var horizontalGifImageList []string
var horizontal2KImageList []string
var horizontal1KImageList []string
var horizontal3KImageList []string
var horizontal4KImageList []string
var horizontal5KImageList []string
var horizontal6KImageList []string
var horizontal7KImageList []string
var horizontal8KImageList []string
var horizontal9KImageList []string
var horizontalHKImageList []string

// 标准横向图片
var horizontalStandard720PImageList []string
var horizontalStandard1080PImageList []string
var horizontalStandard2KImageList []string
var horizontalStandard4KImageList []string
var horizontalStandard5KImageList []string
var horizontalStandard8KImageList []string

// 垂直图片
var verticalGifImageList []string
var verticalUnHandleImageList []string
var vertical1KImageList []string
var vertical2KImageList []string
var vertical3KImageList []string
var vertical4KImageList []string
var vertical5KImageList []string
var vertical6KImageList []string
var vertical7KImageList []string
var vertical8KImageList []string
var vertical9KImageList []string
var verticalHKImageList []string

// 等比图片
var squareUnHandleImageList []string
var squareGifImageList []string
var square1KImageList []string
var square2KImageList []string
var square3KImageList []string
var square4KImageList []string
var square5KImageList []string
var square6KImageList []string
var square7KImageList []string
var square8KImageList []string
var square9KImageList []string
var squareHKImageList []string
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
		handleSquareImage(currentImagePath, width, suffix)
	}
	moveNormalImage(rootDir, uid)
	moveHorizontalImage(rootDir, uid)
	moveVerticalImage(rootDir, uid)
	moveSquareImage(rootDir, uid)
}

func moveSquareImage(rootDir string, uid string) {
	pathSeparator := string(os.PathSeparator)
	squareUnHandleImagePath := rootDir + pathSeparator + uid + "_图片_等比_未处理"
	squareGifImagePath := rootDir + pathSeparator + uid + "_图片_等比_GIF"
	square1KImagePath := rootDir + pathSeparator + uid + "_图片_等比_1k"
	square2KImagePath := rootDir + pathSeparator + uid + "_图片_等比_2k"
	square3KImagePath := rootDir + pathSeparator + uid + "_图片_等比_3k"
	square4KImagePath := rootDir + pathSeparator + uid + "_图片_等比_4k"
	square5KImagePath := rootDir + pathSeparator + uid + "_图片_等比_5k"
	square6KImagePath := rootDir + pathSeparator + uid + "_图片_等比_6k"
	square7KImagePath := rootDir + pathSeparator + uid + "_图片_等比_7k"
	square8KImagePath := rootDir + pathSeparator + uid + "_图片_等比_8k"
	square9KImagePath := rootDir + pathSeparator + uid + "_图片_等比_9k"
	squareHKImagePath := rootDir + pathSeparator + uid + "_图片_等比_原图"
	if len(squareUnHandleImageList) > 0 {
		util.CreateDir(squareUnHandleImagePath)
		doMoveFileToDir(squareUnHandleImageList, squareUnHandleImagePath)
	}
	if len(squareGifImageList) > 0 {
		util.CreateDir(squareGifImagePath)
		doMoveFileToDir(squareGifImageList, squareGifImagePath)
	}
	if len(square1KImageList) > 0 {
		util.CreateDir(square1KImagePath)
		doMoveFileToDir(square1KImageList, square1KImagePath)
	}
	if len(square2KImageList) > 0 {
		util.CreateDir(square2KImagePath)
		doMoveFileToDir(square2KImageList, square2KImagePath)
	}
	if len(square3KImageList) > 0 {
		util.CreateDir(square3KImagePath)
		doMoveFileToDir(square3KImageList, square3KImagePath)
	}
	if len(square4KImageList) > 0 {
		util.CreateDir(square4KImagePath)
		doMoveFileToDir(square4KImageList, square4KImagePath)
	}
	if len(square5KImageList) > 0 {
		util.CreateDir(square5KImagePath)
		doMoveFileToDir(square5KImageList, square5KImagePath)
	}
	if len(square6KImageList) > 0 {
		util.CreateDir(square6KImagePath)
		doMoveFileToDir(square6KImageList, square6KImagePath)
	}
	if len(square7KImageList) > 0 {
		util.CreateDir(square7KImagePath)
		doMoveFileToDir(square7KImageList, square7KImagePath)
	}
	if len(square8KImageList) > 0 {
		util.CreateDir(square8KImagePath)
		doMoveFileToDir(square8KImageList, square8KImagePath)
	}
	if len(square9KImageList) > 0 {
		util.CreateDir(square9KImagePath)
		doMoveFileToDir(square9KImageList, square9KImagePath)
	}
	if len(squareHKImageList) > 0 {
		util.CreateDir(squareHKImagePath)
		doMoveFileToDir(squareHKImageList, squareHKImagePath)
	}
}

func handleSquareImage(currentImagePath string, width int, suffix string) {
	if strings.EqualFold(suffix, ".gif") {
		squareGifImageList = append(squareGifImageList, currentImagePath)
		return
	}
	if width < caleHorizontalPix(1) {
		normalImageList = append(normalImageList, currentImagePath)
	} else if width >= caleHorizontalPix(1) && width < caleHorizontalPix(2) {
		square1KImageList = append(square1KImageList, currentImagePath)
	} else if width >= caleHorizontalPix(2) && width < caleHorizontalPix(3) {
		square2KImageList = append(square2KImageList, currentImagePath)
	} else if width >= caleHorizontalPix(3) && width < caleHorizontalPix(4) {
		square3KImageList = append(square3KImageList, currentImagePath)
	} else if width >= caleHorizontalPix(4) && width < caleHorizontalPix(5) {
		square4KImageList = append(square4KImageList, currentImagePath)
	} else if width >= caleHorizontalPix(5) && width < caleHorizontalPix(6) {
		square5KImageList = append(square5KImageList, currentImagePath)
	} else if width >= caleHorizontalPix(6) && width < caleHorizontalPix(7) {
		square6KImageList = append(square6KImageList, currentImagePath)
	} else if width >= caleHorizontalPix(7) && width < caleHorizontalPix(8) {
		square7KImageList = append(square7KImageList, currentImagePath)
	} else if width >= caleHorizontalPix(8) && width < caleHorizontalPix(9) {
		square8KImageList = append(square8KImageList, currentImagePath)
	} else if width >= caleHorizontalPix(9) && width < caleHorizontalPix(10) {
		square9KImageList = append(square9KImageList, currentImagePath)
	} else if width >= caleHorizontalPix(10) {
		squareHKImageList = append(squareHKImageList, currentImagePath)
	} else {
		// 未处理的等比图片
		squareUnHandleImageList = append(squareUnHandleImageList, currentImagePath)
	}
}

// 移动垂直图片
func moveVerticalImage(rootDir string, uid string) {
	pathSeparator := string(os.PathSeparator)
	verticalGifImagePath := rootDir + pathSeparator + uid + "_图片_竖屏_GIF"
	verticalUnHandleImagePath := rootDir + pathSeparator + uid + "_图片_竖屏_未处理"
	vertical1KImagePath := rootDir + pathSeparator + uid + "_图片_竖屏_1k"
	vertical2KImagePath := rootDir + pathSeparator + uid + "_图片_竖屏_2k"
	vertical3KImagePath := rootDir + pathSeparator + uid + "_图片_竖屏_3k"
	vertical4KImagePath := rootDir + pathSeparator + uid + "_图片_竖屏_4k"
	vertical5KImagePath := rootDir + pathSeparator + uid + "_图片_竖屏_5k"
	vertical6KImagePath := rootDir + pathSeparator + uid + "_图片_竖屏_6k"
	vertical7KImagePath := rootDir + pathSeparator + uid + "_图片_竖屏_7k"
	vertical8KImagePath := rootDir + pathSeparator + uid + "_图片_竖屏_8k"
	vertical9KImagePath := rootDir + pathSeparator + uid + "_图片_竖屏_9k"
	verticalHKImagePath := rootDir + pathSeparator + uid + "_图片_竖屏_原图"
	if len(verticalGifImageList) > 0 {
		util.CreateDir(verticalGifImagePath)
		doMoveFileToDir(verticalGifImageList, verticalGifImagePath)
	}
	if len(verticalUnHandleImageList) > 0 {
		util.CreateDir(verticalUnHandleImagePath)
		doMoveFileToDir(verticalUnHandleImageList, verticalUnHandleImagePath)
	}
	if len(vertical1KImageList) > 0 {
		util.CreateDir(vertical1KImagePath)
		doMoveFileToDir(vertical1KImageList, vertical1KImagePath)
	}
	if len(vertical2KImageList) > 0 {
		util.CreateDir(vertical2KImagePath)
		doMoveFileToDir(vertical2KImageList, vertical2KImagePath)
	}
	if len(vertical3KImageList) > 0 {
		util.CreateDir(vertical3KImagePath)
		doMoveFileToDir(vertical3KImageList, vertical3KImagePath)
	}
	if len(vertical4KImageList) > 0 {
		util.CreateDir(vertical4KImagePath)
		doMoveFileToDir(vertical4KImageList, vertical4KImagePath)
	}
	if len(vertical5KImageList) > 0 {
		util.CreateDir(vertical5KImagePath)
		doMoveFileToDir(vertical5KImageList, vertical5KImagePath)
	}
	if len(vertical6KImageList) > 0 {
		util.CreateDir(vertical6KImagePath)
		doMoveFileToDir(vertical6KImageList, vertical6KImagePath)
	}
	if len(vertical7KImageList) > 0 {
		util.CreateDir(vertical7KImagePath)
		doMoveFileToDir(vertical7KImageList, vertical7KImagePath)
	}
	if len(vertical8KImageList) > 0 {
		util.CreateDir(vertical8KImagePath)
		doMoveFileToDir(vertical8KImageList, vertical8KImagePath)
	}
	if len(vertical9KImageList) > 0 {
		util.CreateDir(vertical9KImagePath)
		doMoveFileToDir(vertical9KImageList, vertical9KImagePath)
	}
	if len(verticalHKImageList) > 0 {
		util.CreateDir(verticalHKImagePath)
		doMoveFileToDir(verticalHKImageList, verticalHKImagePath)
	}
}

// 移动水平图片
func moveHorizontalImage(rootDir string, uid string) {
	pathSeparator := string(os.PathSeparator)
	horizontalUnHandleImagePath := rootDir + pathSeparator + uid + "_图片_横屏_未处理"
	horizontalGifImagePath := rootDir + pathSeparator + uid + "_图片_横屏_GIF"
	horizontal1KImagePath := rootDir + pathSeparator + uid + "_图片_横屏_1k"
	horizontal2KImagePath := rootDir + pathSeparator + uid + "_图片_横屏_2k"
	horizontal3KImagePath := rootDir + pathSeparator + uid + "_图片_横屏_3k"
	horizontal4KImagePath := rootDir + pathSeparator + uid + "_图片_横屏_4k"
	horizontal5KImagePath := rootDir + pathSeparator + uid + "_图片_横屏_5k"
	horizontal6KImagePath := rootDir + pathSeparator + uid + "_图片_横屏_6k"
	horizontal7KImagePath := rootDir + pathSeparator + uid + "_图片_横屏_7k"
	horizontal8KImagePath := rootDir + pathSeparator + uid + "_图片_横屏_8k"
	horizontal9KImagePath := rootDir + pathSeparator + uid + "_图片_横屏_9k"
	horizontalHKImagePath := rootDir + pathSeparator + uid + "_图片_横屏_原图"
	horizontalStandard720PImagePath := rootDir + pathSeparator + uid + "_图片_横屏_720P"
	horizontalStandard1080PImagePath := rootDir + pathSeparator + uid + "_图片_横屏_1080P"
	horizontalStandard2KImagePath := rootDir + pathSeparator + uid + "_图片_横屏_2KP"
	horizontalStandard4KImagePath := rootDir + pathSeparator + uid + "_图片_横屏_4KP"
	horizontalStandard5KImagePath := rootDir + pathSeparator + uid + "_图片_横屏_5KP"
	horizontalStandard8KImagePath := rootDir + pathSeparator + uid + "_图片_横屏_8KP"
	if len(horizontalUnHandleImageList) > 0 {
		util.CreateDir(horizontalUnHandleImagePath)
		doMoveFileToDir(horizontalUnHandleImageList, horizontalUnHandleImagePath)
	}
	if len(horizontalGifImageList) > 0 {
		util.CreateDir(horizontalGifImagePath)
		doMoveFileToDir(horizontalGifImageList, horizontalGifImagePath)
	}
	if len(horizontal1KImageList) > 0 {
		util.CreateDir(horizontal1KImagePath)
		doMoveFileToDir(horizontal1KImageList, horizontal1KImagePath)
	}
	if len(horizontal2KImageList) > 0 {
		util.CreateDir(horizontal2KImagePath)
		doMoveFileToDir(horizontal2KImageList, horizontal2KImagePath)
	}
	if len(horizontal3KImageList) > 0 {
		util.CreateDir(horizontal3KImagePath)
		doMoveFileToDir(horizontal3KImageList, horizontal3KImagePath)
	}
	if len(horizontal4KImageList) > 0 {
		util.CreateDir(horizontal4KImagePath)
		doMoveFileToDir(horizontal4KImageList, horizontal4KImagePath)
	}
	if len(horizontal5KImageList) > 0 {
		util.CreateDir(horizontal5KImagePath)
		doMoveFileToDir(horizontal5KImageList, horizontal5KImagePath)
	}
	if len(horizontal6KImageList) > 0 {
		util.CreateDir(horizontal6KImagePath)
		doMoveFileToDir(horizontal6KImageList, horizontal6KImagePath)
	}
	if len(horizontal7KImageList) > 0 {
		util.CreateDir(horizontal7KImagePath)
		doMoveFileToDir(horizontal7KImageList, horizontal7KImagePath)
	}
	if len(horizontal8KImageList) > 0 {
		util.CreateDir(horizontal8KImagePath)
		doMoveFileToDir(horizontal8KImageList, horizontal8KImagePath)
	}
	if len(horizontal9KImageList) > 0 {
		util.CreateDir(horizontal9KImagePath)
		doMoveFileToDir(horizontal9KImageList, horizontal9KImagePath)
	}
	if len(horizontalHKImageList) > 0 {
		util.CreateDir(horizontalHKImagePath)
		doMoveFileToDir(horizontalHKImageList, horizontalHKImagePath)
	}
	if len(horizontalStandard720PImageList) > 0 {
		util.CreateDir(horizontalStandard720PImagePath)
		doMoveFileToDir(horizontalStandard720PImageList, horizontalStandard720PImagePath)
	}
	if len(horizontalStandard1080PImageList) > 0 {
		util.CreateDir(horizontalStandard1080PImagePath)
		doMoveFileToDir(horizontalStandard1080PImageList, horizontalStandard1080PImagePath)
	}
	if len(horizontalStandard2KImageList) > 0 {
		util.CreateDir(horizontalStandard2KImagePath)
		doMoveFileToDir(horizontalStandard2KImageList, horizontalStandard2KImagePath)
	}
	if len(horizontalStandard4KImageList) > 0 {
		util.CreateDir(horizontalStandard4KImagePath)
		doMoveFileToDir(horizontalStandard4KImageList, horizontalStandard4KImagePath)
	}
	if len(horizontalStandard5KImageList) > 0 {
		util.CreateDir(horizontalStandard5KImagePath)
		doMoveFileToDir(horizontalStandard5KImageList, horizontalStandard5KImagePath)
	}
	if len(horizontalStandard8KImageList) > 0 {
		util.CreateDir(horizontalStandard8KImagePath)
		doMoveFileToDir(horizontalStandard8KImageList, horizontalStandard8KImagePath)
	}
}

// 移动图片
func moveNormalImage(rootDir string, uid string) {
	pathSeparator := string(os.PathSeparator)
	allNormalImagePath := rootDir + pathSeparator + uid + "_图片_普通"
	if len(normalImageList) > 0 {
		util.CreateDir(allNormalImagePath)
		doMoveFileToDir(normalImageList, allNormalImagePath)
	}
}

// 处理垂直图片
func handleVerticalImage(currentImagePath string, height int, suffix string) {
	if strings.EqualFold(suffix, ".gif") {
		verticalGifImageList = append(verticalGifImageList, currentImagePath)
		return
	}
	if height < caleVerticalPix(1) {
		normalImageList = append(normalImageList, currentImagePath)
	} else if height >= caleVerticalPix(1) && height < caleVerticalPix(2) {
		vertical1KImageList = append(vertical1KImageList, currentImagePath)
	} else if height >= caleVerticalPix(2) && height < caleVerticalPix(3) {
		vertical2KImageList = append(vertical2KImageList, currentImagePath)
	} else if height >= caleVerticalPix(3) && height < caleVerticalPix(4) {
		vertical3KImageList = append(vertical3KImageList, currentImagePath)
	} else if height >= caleVerticalPix(4) && height < caleVerticalPix(5) {
		vertical4KImageList = append(vertical4KImageList, currentImagePath)
	} else if height >= caleVerticalPix(5) && height < caleVerticalPix(6) {
		vertical5KImageList = append(vertical5KImageList, currentImagePath)
	} else if height >= caleVerticalPix(6) && height < caleVerticalPix(7) {
		vertical6KImageList = append(vertical6KImageList, currentImagePath)
	} else if height >= caleVerticalPix(7) && height < caleVerticalPix(8) {
		vertical7KImageList = append(vertical7KImageList, currentImagePath)
	} else if height >= caleVerticalPix(8) && height < caleVerticalPix(9) {
		vertical8KImageList = append(vertical8KImageList, currentImagePath)
	} else if height >= caleVerticalPix(9) && height < caleVerticalPix(10) {
		vertical9KImageList = append(vertical9KImageList, currentImagePath)
	} else if height >= caleVerticalPix(10) {
		verticalHKImageList = append(verticalHKImageList, currentImagePath)
	} else {
		// 未处理的垂直图片
		verticalUnHandleImageList = append(verticalUnHandleImageList, currentImagePath)
	}
}

// 处理横向图片
func handleHorizontalImage(currentImagePath string, width int, height int, suffix string) {
	if strings.EqualFold(suffix, ".gif") {
		horizontalGifImageList = append(horizontalGifImageList, currentImagePath)
		return
	}
	if width < caleHorizontalPix(1) {
		normalImageList = append(normalImageList, currentImagePath)
	} else if width >= caleHorizontalPix(1) && width < caleHorizontalPix(2) {
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
		horizontal1KImageList = append(horizontal1KImageList, currentImagePath)
	} else if width >= caleHorizontalPix(2) && width < caleHorizontalPix(3) {
		// 2560 x 1440 -> 2k
		if width == 2560 && height == 1440 {
			horizontalStandard2KImageList = append(horizontalStandard2KImageList, currentImagePath)
			return
		}
		horizontal2KImageList = append(horizontal2KImageList, currentImagePath)
	} else if width >= caleHorizontalPix(3) && width < caleHorizontalPix(4) {
		// 3840 x 2160 -> 4k
		if width == 3840 && height == 2160 {
			horizontalStandard4KImageList = append(horizontalStandard4KImageList, currentImagePath)
			return
		}
		horizontal3KImageList = append(horizontal3KImageList, currentImagePath)
	} else if width >= caleHorizontalPix(4) && width < caleHorizontalPix(5) {
		horizontal4KImageList = append(horizontal4KImageList, currentImagePath)
	} else if width >= caleHorizontalPix(5) && width < caleHorizontalPix(6) {
		// 5120 x 2880 -> 5k
		if width == 5120 && height == 2880 {
			horizontalStandard5KImageList = append(horizontalStandard5KImageList, currentImagePath)
			return
		}
		horizontal5KImageList = append(horizontal5KImageList, currentImagePath)
	} else if width >= caleHorizontalPix(6) && width < caleHorizontalPix(7) {
		horizontal6KImageList = append(horizontal6KImageList, currentImagePath)
	} else if width >= caleHorizontalPix(7) && width < caleHorizontalPix(8) {
		// 7680 x 4320 -> 8k
		if width == 7680 && height == 4320 {
			horizontalStandard8KImageList = append(horizontalStandard8KImageList, currentImagePath)
			return
		}
		horizontal7KImageList = append(horizontal7KImageList, currentImagePath)
	} else if width >= caleHorizontalPix(8) && width < caleHorizontalPix(9) {
		horizontal8KImageList = append(horizontal8KImageList, currentImagePath)
	} else if width >= caleHorizontalPix(9) && width < caleHorizontalPix(10) {
		horizontal9KImageList = append(horizontal9KImageList, currentImagePath)
	} else if width >= caleHorizontalPix(10) {
		horizontalHKImageList = append(horizontalHKImageList, currentImagePath)
	} else {
		// 未分类的横向图片
		horizontalUnHandleImageList = append(horizontalUnHandleImageList, currentImagePath)
	}
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
