package util

import (
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// CheckFileIsExist 判断文件是否存在，存在返回 true，不存在返回false
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// WriteByteArraysToFile 写字节数组到文件
func WriteByteArraysToFile(content []byte, filename string) error {
	return ioutil.WriteFile(filename, content, 0777)
}

// String2int 字符串转int
func String2int(str string) int {
	intValue, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return intValue
}

// CreateDir 创建目录
func CreateDir(dirPath string) bool {
	err := os.Mkdir(dirPath, 0755)
	return err == nil
}

// ReadFileMimeInfo 获取文件mime信息
func ReadFileMimeInfo(filepath string) *mimetype.MIME {
	mt, err := mimetype.DetectFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	return mt
}

// SecondsToHms 将秒数转换为小时、分钟、秒的格式
func SecondsToHms(seconds int) string {
	t := time.Duration(seconds) * time.Second
	h := t / time.Hour
	t -= h * time.Hour
	m := t / time.Minute
	t -= m * time.Minute
	s := t / time.Second
	return fmt.Sprintf("%02d-%02d-%02d", h, m, s)
}

// CalcPercentage 计算percentage相对于total的百分比
func CalcPercentage(percentage int, total int) int {
	return int(float64(percentage) / float64(total) * 100)
}

// 获取文件所在文件夹
func GetFileDirectory(filePath string) string {
	directoryPath, _ := filepath.Split(filePath)
	return directoryPath
}
