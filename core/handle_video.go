package core

import (
	"OdMediaPicker/util"
	"OdMediaPicker/vars"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/redmask-hb/GoSimplePrint/goPrint"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

//go:embed files/ffprobe.exe
var ffprobeWin64 []byte

var videoTag = "[PickVideo]"                              // 标记文件已经被整理过
var ignoreVideoPathList []string                          // 忽略的文件路径
var readErrorVideoPathList []string                       // 读取信息异常的路径
var videoPath2WidthHeightMap = make(map[string]string)    // 视频路径和宽高比
var videoPath2WidthHeightTagMap = make(map[string]string) // 视频路径和宽高比[640x480]
var videoPath2DurationMap = make(map[string]string)       // 视频路径和时长
// 支持的视频格式
var supportVideoTypes = []string{
	".ts",
	".flv",
	".rm",
	".avi",
	".mp4",
	".mov",
	".mpg",
	".mkv",
	".m4v",
	".rmvb",
	".3gp",
	".3g2",
	".webm",
	".wmv",
}

// 水平视频
var horizontalNormalVideoList []string
var horizontalUnHandleVideoList []string
var horizontalGifVideoList []string

var horizontal1KVideoList []string
var horizontal2KVideoList []string
var horizontal3KVideoList []string
var horizontal4KVideoList []string
var horizontal5KVideoList []string
var horizontal6KVideoList []string
var horizontal7KVideoList []string
var horizontal8KVideoList []string
var horizontal9KVideoList []string
var horizontalHKVideoList []string

// 标准横向视频
var horizontalStandard720PVideoList []string
var horizontalStandard1080PVideoList []string
var horizontalStandard2KVideoList []string
var horizontalStandard4KVideoList []string
var horizontalStandard5KVideoList []string
var horizontalStandard8KVideoList []string

// 垂直视频
var verticalNormalVideoList []string
var verticalGifVideoList []string
var vertical1KVideoList []string
var vertical2KVideoList []string
var vertical3KVideoList []string
var vertical4KVideoList []string
var vertical5KVideoList []string
var vertical6KVideoList []string
var vertical7KVideoList []string
var vertical8KVideoList []string
var vertical9KVideoList []string
var verticalHKVideoList []string
var verticalUnHandleVideoList []string

// 等比视频
var squareNormalVideoList []string
var squareGifVideoList []string
var square1KVideoList []string
var square2KVideoList []string
var square3KVideoList []string
var square4KVideoList []string
var square5KVideoList []string
var square6KVideoList []string
var square7KVideoList []string
var square8KVideoList []string
var square9KVideoList []string
var squareHKVideoList []string
var squareUnHandleVideoList []string

func DoHandleVideo(rootDir string) {
	// 释放ffprobe
	readerFileName := "./ffprobe.exe"
	if util.CheckFileIsExist(readerFileName) {
		_ = os.Remove(readerFileName)
	}
	err := util.WriteByteArraysToFile(ffprobeWin64, readerFileName)
	if err != nil {
		fmt.Println("=== 释放解码器失败, 5秒后将自动退出", err)
		time.Sleep(time.Second * 5)
		return
	}
	total := len(vars.GlobalVideoPathList) // 总数
	successCount := 0                      // 成功数
	errorCount := 0                        // 失败数
	ignoreCount := 0                       // 忽略数
	for _, videoFilePath := range vars.GlobalVideoPathList {
		suffix := vars.GlobalFilePath2FileExtMap[videoFilePath]
		if isSupportVideo(suffix) {
			width, height, err := readVideoWidthHeight(videoFilePath)
			if err == nil {
				successCount = successCount + 1
				videoPath2WidthHeightMap[videoFilePath] = fmt.Sprintf("%d-%d", width, height)
				videoPath2WidthHeightTagMap[videoFilePath] = fmt.Sprintf("[%dx%d]", width, height)
				fmt.Printf("=== Video总数: %d, 已读取Info: %d, 成功数: %d, 失败数: %d \n", total, successCount+errorCount+ignoreCount, successCount, errorCount)
				duration := readVideoDuration(videoFilePath)
				if duration == 0 {
					videoPath2DurationMap[videoFilePath] = "00-00-00"
				} else {
					videoPath2DurationMap[videoFilePath] = util.SecondsToHms(duration)
				}
			} else {
				errorCount = errorCount + 1
				readErrorVideoPathList = append(readErrorVideoPathList, videoFilePath)
				fmt.Printf("=== 异常视频: %s \n", videoFilePath)
			}
			continue
		}
		// 其他的直接先忽略吧, 爱改改, 不改拉倒
		ignoreCount = ignoreCount + 1
		ignoreVideoPathList = append(ignoreVideoPathList, videoFilePath)
	}
	//uuid := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	if len(readErrorVideoPathList) > 0 {
		readInfoErrorPath := rootDir + string(os.PathSeparator) + "读取异常"
		if util.CreateDir(readInfoErrorPath) {
			doMoveFileToDir(readErrorVideoPathList, readInfoErrorPath)
		}
	}
	if len(ignoreVideoPathList) > 0 {
		ignorePath := rootDir + string(os.PathSeparator) + "已忽略"
		if util.CreateDir(ignorePath) {
			doMoveFileToDir(ignoreVideoPathList, ignorePath)
		}
	}
	doPickVideoFile(rootDir, videoPath2WidthHeightMap)
	if util.CheckFileIsExist(readerFileName) {
		_ = os.Remove(readerFileName)
	}
	fmt.Printf("=== 视频处理完毕 \n\n")
}

// getVideoDuration 使用ffprobe获取视频时长
func getVideoDuration(ffmpegExecPath string, videoPath string) (float64, error) {
	// ffprobe命令，-v error 用于减少输出信息，-show_entries format=duration -of compact=p=0,nk=1 用于只输出时长
	cmd := exec.Command(ffmpegExecPath, "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", videoPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe failed with error: %v, stderr: %q", err, stderr.String())
	}

	// 解析输出的时长字符串为浮点数
	durationStr := strings.TrimSpace(string(output))
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %q, error: %v", durationStr, err)
	}

	return duration, nil
}

func getVideoResolution(ffmpegExecPath string, filePath string) (width int, height int, err error) {
	// 构建ffprobe命令
	cmd := exec.Command(ffmpegExecPath, "-v", "error", "-show_entries", "stream=width,height", "-of", "csv=p=0:s=x", filePath)
	// 执行命令并捕获输出
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to run ffprobe: %w", err)
	}
	// 解析输出字符串，格式应为 "宽度,高度"
	resolutionStr := strings.TrimSpace(string(output))
	parts := strings.Split(resolutionStr, ",")
	if len(parts) == 2 {
		width = util.String2int(parts[0])
		height = util.String2int(parts[1])
		return width, height, nil
	}
	parts = strings.Split(resolutionStr, "x")
	if len(parts) == 2 {
		width = util.String2int(parts[0])
		height = util.String2int(parts[1])
		return width, height, nil
	}
	parts = strings.Split(resolutionStr, "\r\n\r\n\r\n\r\n")
	if len(parts) == 2 {
		tempHw := parts[0]
		parts = strings.Split(tempHw, "x")
		if len(parts) == 2 {
			width = util.String2int(parts[0])
			height = util.String2int(parts[1])
			return width, height, nil
		}
	}
	parts = strings.Split(resolutionStr, "x")
	if len(parts) == 3 {
		width = util.String2int(parts[0])
		height = util.String2int(parts[1])
		return width, height, nil
	}
	return 0, 0, fmt.Errorf("invalid resolution format: %s", resolutionStr)
}

// 获取视频的时长，单位秒
func readVideoDuration(videoFilePath string) int {
	duration, err := getVideoDuration("./ffprobe.exe", videoFilePath)
	if err != nil {
		fmt.Println("=== Error getting video duration:", err)
		return 0
	}
	//fmt.Printf("=== Video duration: %.2f seconds\n", duration)
	return int(math.Floor(duration)) // 向下取整
}

// 获取视频的分辨率
func readVideoWidthHeight(videoFilePath string) (int, int, error) {
	width, height, err := getVideoResolution("./ffprobe.exe", videoFilePath)
	if err != nil {
		fmt.Printf("=== Error getting resolution: %v\n", err)
		return 0, 0, err
	}
	//fmt.Printf("=== Video resolution: %dx%d\n", width, height)
	return width, height, nil
}

// 条件视频并分组存放
func doPickVideoFile(rootDir string, videoPath2WidthHeightMap map[string]string) {
	if len(videoPath2WidthHeightMap) == 0 {
		fmt.Printf("=== 当前目录下没有扫描到视频文件, %s \n", rootDir)
		readerFileName := "./ffprobe.exe"
		if util.CheckFileIsExist(readerFileName) {
			_ = os.Remove(readerFileName)
		}
		return
	}
	for currentVideoPath, infoStr := range videoPath2WidthHeightMap {
		width2Height := strings.Split(infoStr, "-")
		width := util.String2int(width2Height[0])
		height := util.String2int(width2Height[1])
		suffix := vars.GlobalFilePath2FileExtMap[currentVideoPath]
		if width > height {
			handleHorizontalVideo(currentVideoPath, width, height, suffix)
			continue
		}
		if width < height {
			handleVerticalVideo(currentVideoPath, height, suffix)
			continue
		}
		handleSquareVideo(currentVideoPath, width, height, suffix)
	}
	moveNormalVideo(rootDir)
	moveHorizontalVideo(rootDir)
	moveVerticalVideo(rootDir)
	moveSquareVideo(rootDir)
}

// 移动垂直视频
func moveVerticalVideo(rootDir string) {
	if len(vertical1KVideoList) > 0 {
		renameFileV2("[V][1k]", vertical1KVideoList)
	}
	if len(vertical2KVideoList) > 0 {
		renameFileV2("[V][2k]", vertical2KVideoList)
	}
	if len(vertical3KVideoList) > 0 {
		renameFileV2("[V][3k]", vertical3KVideoList)
	}
	if len(vertical4KVideoList) > 0 {
		renameFileV2("[V][4k]", vertical4KVideoList)
	}
	if len(vertical5KVideoList) > 0 {
		renameFileV2("[V][5k]", vertical5KVideoList)
	}
	if len(vertical6KVideoList) > 0 {
		renameFileV2("[V][6k]", vertical6KVideoList)
	}
	if len(vertical7KVideoList) > 0 {
		renameFileV2("[V][7k]", vertical7KVideoList)
	}
	if len(vertical8KVideoList) > 0 {
		renameFileV2("[V][8k]", vertical8KVideoList)
	}
	if len(vertical9KVideoList) > 0 {
		renameFileV2("[V][9k]", vertical9KVideoList)
	}
	if len(verticalHKVideoList) > 0 {
		renameFileV2("[V][原]", verticalHKVideoList)
	}
	if len(verticalUnHandleVideoList) > 0 {
		renameFileV2("[V][UnHandle]", verticalUnHandleVideoList)
	}
}

// 移动文件到根目录
func renameFile(rootDir string, modelType string, videoList []string, pathSeparator string) {
	total := len(videoList)
	var count = 0
	bar := goPrint.NewBar(100)
	bar.SetNotice("=== 重命名文件：")
	bar.SetGraph(">")
	for _, videoFilePath := range videoList {
		wh := videoPath2WidthHeightTagMap[videoFilePath]
		fileName := vars.GlobalFilePath2FileNameMap[videoFilePath]
		if strings.Contains(fileName, videoTag) { // 处理过了
			fileNames := strings.Split(fileName, videoTag)
			if len(fileNames) == 2 {
				fileName = fileNames[1]
				targetFilePath := rootDir + pathSeparator + "[" + videoPath2DurationMap[videoFilePath] + "]" + modelType + wh + videoTag + fileName
				err := os.Rename(videoFilePath, targetFilePath)
				if err != nil {
					fmt.Printf("=== 重命名异常: %s \n", videoFilePath)
				}
			}
		} else {
			targetFilePath := rootDir + pathSeparator + "[" + videoPath2DurationMap[videoFilePath] + "]" + modelType + wh + videoTag + " - " + fileName
			err := os.Rename(videoFilePath, targetFilePath)
			if err != nil {
				fmt.Printf("=== 重命名异常: %s \n", videoFilePath)
			}
		}
		count = count + 1
		bar.PrintBar(util.CalcPercentage(count, total))
	}
	bar.PrintEnd("=== Finish")
}

// 移动文件到原目录
func renameFileV2(modelType string, videoList []string) {
	total := len(videoList)
	var count = 0
	bar := goPrint.NewBar(100)
	bar.SetNotice("=== 重命名文件：")
	bar.SetGraph(">")
	for _, videoFilePath := range videoList {
		wh := videoPath2WidthHeightTagMap[videoFilePath]
		fileName := vars.GlobalFilePath2FileNameMap[videoFilePath]
		filePath := util.GetFileDirectory(videoFilePath)
		if strings.Contains(fileName, videoTag) { // 处理过了
			fileNames := strings.Split(fileName, videoTag)
			if len(fileNames) == 2 {
				fileName = fileNames[1]
				targetFilePath := filePath + "[" + videoPath2DurationMap[videoFilePath] + "]" + modelType + wh + videoTag + fileName
				err := os.Rename(videoFilePath, targetFilePath)
				if err != nil {
					fmt.Printf("=== 重命名异常: %s \n", videoFilePath)
				}
			}
		} else {
			targetFilePath := filePath + "[" + videoPath2DurationMap[videoFilePath] + "]" + modelType + wh + videoTag + " - " + fileName
			err := os.Rename(videoFilePath, targetFilePath)
			if err != nil {
				fmt.Printf("=== 重命名异常: %s \n", videoFilePath)
			}
		}
		count = count + 1
		bar.PrintBar(util.CalcPercentage(count, total))
	}
	bar.PrintEnd("=== Finish")
}

// 移动水平视频
func moveHorizontalVideo(rootDir string) {
	if len(horizontal1KVideoList) > 0 {
		renameFileV2("[H][1k]", horizontal1KVideoList)
	}
	if len(horizontal2KVideoList) > 0 {
		renameFileV2("[H][2k]", horizontal2KVideoList)
	}
	if len(horizontal3KVideoList) > 0 {
		renameFileV2("[H][3k]", horizontal3KVideoList)
	}
	if len(horizontal4KVideoList) > 0 {
		renameFileV2("[H][4k]", horizontal4KVideoList)
	}
	if len(horizontal5KVideoList) > 0 {
		renameFileV2("[H][5k]", horizontal5KVideoList)
	}
	if len(horizontal6KVideoList) > 0 {
		renameFileV2("[H][6k]", horizontal6KVideoList)
	}
	if len(horizontal7KVideoList) > 0 {
		renameFileV2("[H][7k]", horizontal7KVideoList)
	}
	if len(horizontal8KVideoList) > 0 {
		renameFileV2("[H][8k]", horizontal8KVideoList)
	}
	if len(horizontal9KVideoList) > 0 {
		renameFileV2("[H][9k]", horizontal9KVideoList)
	}
	if len(horizontalHKVideoList) > 0 {
		renameFileV2("[H][原]", horizontalHKVideoList)
	}
	if len(horizontalStandard720PVideoList) > 0 {
		renameFileV2("[H][720P]", horizontalStandard720PVideoList)
	}
	if len(horizontalStandard1080PVideoList) > 0 {
		renameFileV2("[H][1080P]", horizontalStandard1080PVideoList)
	}
	if len(horizontalStandard2KVideoList) > 0 {
		renameFileV2("[H][2KP]", horizontalStandard2KVideoList)
	}
	if len(horizontalStandard4KVideoList) > 0 {
		renameFileV2("[H][4KP]", horizontalStandard4KVideoList)
	}
	if len(horizontalStandard5KVideoList) > 0 {
		renameFileV2("[H][5KP]", horizontalStandard5KVideoList)
	}
	if len(horizontalStandard8KVideoList) > 0 {
		renameFileV2("[H][8KP]", horizontalStandard8KVideoList)
	}
	if len(horizontalUnHandleVideoList) > 0 {
		renameFileV2("[H][UnHandle]", horizontalUnHandleVideoList)
	}
}

// 移动等比视频
func moveSquareVideo(rootDir string) {
	if len(square1KVideoList) > 0 {
		renameFileV2("[M][1k]", square1KVideoList)
	}
	if len(square2KVideoList) > 0 {
		renameFileV2("[M][2k]", square2KVideoList)
	}
	if len(square3KVideoList) > 0 {
		renameFileV2("[M][3k]", square3KVideoList)
	}
	if len(square4KVideoList) > 0 {
		renameFileV2("[M][4k]", square4KVideoList)
	}
	if len(square5KVideoList) > 0 {
		renameFileV2("[M][5k]", square5KVideoList)
	}
	if len(square6KVideoList) > 0 {
		renameFileV2("[M][6k]", square6KVideoList)
	}
	if len(square7KVideoList) > 0 {
		renameFileV2("[M][7k]", square7KVideoList)
	}
	if len(square8KVideoList) > 0 {
		renameFileV2("[M][8k]", square8KVideoList)
	}
	if len(square9KVideoList) > 0 {
		renameFileV2("[M][9k]", square9KVideoList)
	}
	if len(squareHKVideoList) > 0 {
		renameFileV2("[M][原]", squareHKVideoList)
	}
	if len(squareUnHandleVideoList) > 0 {
		renameFileV2("[M][UnHandle]", squareUnHandleVideoList)
	}
}

// 移动普通视频
func moveNormalVideo(rootDir string) {
	//pathSeparator := string(os.PathSeparator)
	if len(horizontalNormalVideoList) > 0 {
		renameFileV2("[L]", horizontalNormalVideoList)
	}
	if len(verticalNormalVideoList) > 0 {
		renameFileV2("[L]", verticalNormalVideoList)
	}
	if len(squareNormalVideoList) > 0 {
		renameFileV2("[L]", squareNormalVideoList)
	}
}

// 处理垂直视频
func handleVerticalVideo(currentVideoPath string, height int, suffix string) {
	if strings.EqualFold(suffix, ".gif") {
		verticalGifVideoList = append(verticalGifVideoList, currentVideoPath)
		return
	}
	if height < 1000 {
		verticalNormalVideoList = append(verticalNormalVideoList, currentVideoPath)
	} else if height >= 1000 && height < 2000 {
		vertical1KVideoList = append(vertical1KVideoList, currentVideoPath)
	} else if height >= 2000 && height < 3000 {
		vertical2KVideoList = append(vertical2KVideoList, currentVideoPath)
	} else if height >= 3000 && height < 4000 {
		vertical3KVideoList = append(vertical3KVideoList, currentVideoPath)
	} else if height >= 4000 && height < 5000 {
		vertical4KVideoList = append(vertical4KVideoList, currentVideoPath)
	} else if height >= 5000 && height < 6000 {
		vertical5KVideoList = append(vertical5KVideoList, currentVideoPath)
	} else if height >= 6000 && height < 7000 {
		vertical6KVideoList = append(vertical6KVideoList, currentVideoPath)
	} else if height >= 7000 && height < 8000 {
		vertical7KVideoList = append(vertical7KVideoList, currentVideoPath)
	} else if height >= 8000 && height < 9000 {
		vertical8KVideoList = append(vertical8KVideoList, currentVideoPath)
	} else if height >= 9000 && height < 10000 {
		vertical9KVideoList = append(vertical9KVideoList, currentVideoPath)
	} else if height >= 10000 {
		verticalHKVideoList = append(verticalHKVideoList, currentVideoPath)
	} else {
		// 未分类的垂直视频
		verticalUnHandleVideoList = append(verticalUnHandleVideoList, currentVideoPath)
	}
}

// 计算水平分辨率
func caleHorizontalPix(level int) int {
	return 1024 * level
}

// 计算垂直分辨率
func caleVerticalPix(level int) int {
	return 768 * level
}

// 处理横向视频
func handleHorizontalVideo(currentVideoPath string, width int, height int, suffix string) {
	if strings.EqualFold(suffix, ".gif") {
		horizontalGifVideoList = append(horizontalGifVideoList, currentVideoPath)
		return
	}
	if width < caleHorizontalPix(1) {
		horizontalNormalVideoList = append(horizontalNormalVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(1) && width < caleHorizontalPix(2) {
		// 1280 x 720 -> 720p
		if width == 1280 && height == 720 {
			horizontalStandard720PVideoList = append(horizontalStandard720PVideoList, currentVideoPath)
			return
		}
		// 1920 x 1080 -> 1080p
		if width == 1920 && height == 1080 {
			horizontalStandard1080PVideoList = append(horizontalStandard1080PVideoList, currentVideoPath)
			return
		}
		horizontal1KVideoList = append(horizontal1KVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(2) && width < caleHorizontalPix(3) {
		// 2560 x 1440 -> 2k
		if width == 2560 && height == 1440 {
			horizontalStandard2KVideoList = append(horizontalStandard2KVideoList, currentVideoPath)
			return
		}
		horizontal2KVideoList = append(horizontal2KVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(3) && width < caleHorizontalPix(4) {
		// 3840 x 2160 -> 4k
		if width == 3840 && height == 2160 {
			horizontalStandard4KVideoList = append(horizontalStandard4KVideoList, currentVideoPath)
			return
		}
		horizontal3KVideoList = append(horizontal3KVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(4) && width < caleHorizontalPix(5) {
		horizontal4KVideoList = append(horizontal4KVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(5) && width < caleHorizontalPix(6) {
		// 5120 x 2880 -> 5k
		if width == 5120 && height == 2880 {
			horizontalStandard5KVideoList = append(horizontalStandard5KVideoList, currentVideoPath)
			return
		}
		horizontal5KVideoList = append(horizontal5KVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(6) && width < caleHorizontalPix(7) {
		horizontal6KVideoList = append(horizontal6KVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(7) && width < caleHorizontalPix(8) {
		// 7680 x 4320 -> 8k
		if width == 7680 && height == 4320 {
			horizontalStandard8KVideoList = append(horizontalStandard8KVideoList, currentVideoPath)
			return
		}
		horizontal7KVideoList = append(horizontal7KVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(8) && width < caleHorizontalPix(9) {
		horizontal8KVideoList = append(horizontal8KVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(9) && width < caleHorizontalPix(10) {
		horizontal9KVideoList = append(horizontal9KVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(10) {
		horizontalHKVideoList = append(horizontalHKVideoList, currentVideoPath)
	} else {
		// 未分类的视频
		horizontalUnHandleVideoList = append(horizontalUnHandleVideoList, currentVideoPath)
	}
}

// 处理等比视频
func handleSquareVideo(currentVideoPath string, width int, height int, suffix string) {
	if strings.EqualFold(suffix, ".gif") {
		squareGifVideoList = append(squareGifVideoList, currentVideoPath)
		return
	}
	if width < caleHorizontalPix(1) {
		squareNormalVideoList = append(squareNormalVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(1) && width < caleHorizontalPix(2) {
		square1KVideoList = append(square1KVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(2) && width < caleHorizontalPix(3) {
		square2KVideoList = append(square2KVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(3) && width < caleHorizontalPix(4) {
		square3KVideoList = append(square3KVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(4) && width < caleHorizontalPix(5) {
		square4KVideoList = append(square4KVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(5) && width < caleHorizontalPix(6) {
		square5KVideoList = append(square5KVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(6) && width < caleHorizontalPix(7) {
		square6KVideoList = append(square6KVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(7) && width < caleHorizontalPix(8) {
		square7KVideoList = append(square7KVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(8) && width < caleHorizontalPix(9) {
		square8KVideoList = append(square8KVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(9) && width < caleHorizontalPix(10) {
		square9KVideoList = append(square9KVideoList, currentVideoPath)
	} else if width >= caleHorizontalPix(10) {
		squareHKVideoList = append(squareHKVideoList, currentVideoPath)
	} else {
		// 未处理的等比视频
		squareUnHandleVideoList = append(squareUnHandleVideoList, currentVideoPath)
	}
}

// 判断是否属于支持的视频
func isSupportVideo(videoType string) bool {
	for _, supportVideoType := range supportVideoTypes {
		if strings.EqualFold(videoType, supportVideoType) {
			return true
		}
	}
	return false
}

// 批量移动文件到目录
func doMoveFileToDir(filePatnList []string, videoDirPath string) {
	total := len(filePatnList)
	var count = 0
	bar := goPrint.NewBar(100)
	bar.SetNotice("=== 移动文件到目录：")
	bar.SetGraph(">")
	pathSeparator := string(os.PathSeparator)
	for _, videoFilePath := range filePatnList {
		moveFileToDir(videoFilePath, videoDirPath+pathSeparator)
		count = count + 1
		bar.PrintBar(util.CalcPercentage(count, total))
	}
	bar.PrintEnd("=== Finish")
}

// 移动文件到目录
func moveFileToDir(sourceFilePath string, targetDirectory string) bool {
	splits := strings.Split(sourceFilePath, string(os.PathSeparator))
	fileName := splits[len(splits)-1]
	targetFilePath := targetDirectory + fileName
	err := os.Rename(sourceFilePath, targetFilePath)
	//fmt.Printf("=== 移动文件, 源: %s, 目标: %s \n", sourceFilePath, targetFilePath)
	return err == nil
}
