package utils

import (
	"GoScheduler/internal/modules/global"
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"path/filepath"
)

func FileExist(file string) bool {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	if os.IsPermission(err) {
		return false
	}

	return true
}

// GBK编码转换为UTF8
func GBK2UTF8(s string) (string, bool) {
	// 将GBK字符串转换为字节数组
	gbkData := []byte(s)

	// 创建GBK到UTF-8的解码器
	reader := transform.NewReader(bytes.NewReader(gbkData), simplifiedchinese.GB18030.NewDecoder())

	// 读取解码后的UTF-8数据
	utf8Data, err := io.ReadAll(reader)
	if err != nil {
		return s, false
	}

	return string(utf8Data), true
}

// 检测目录是否存在
func CreateDirIfNotExists(path ...string) {
	for _, value := range path {
		if FileExist(value) {
			continue
		}
		err := os.Mkdir(value, 0755)
		if err != nil {
			zap.S().Fatal(fmt.Sprintf("创建目录失败:%s", err.Error()))
		}
	}
}

// 获取当前运行目录
func GetWorkDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		zap.S().Warnf("无法获取当前程序路径")
		return "", err
	}

	return filepath.Dir(exePath), nil
}

// IsInstalled 判断应用是否已安装
func IsInstalled() bool {
	_, err := os.Stat(filepath.Join(global.ConfDir, "/install.lock"))
	if os.IsNotExist(err) {
		return false
	}

	return true
}

// CreateInstallLock 创建安装锁文件
func CreateInstallLock() error {
	_, err := os.Create(filepath.Join(global.ConfDir, "/install.lock"))
	if err != nil {
		zap.S().Error("创建安装锁文件conf/install.lock失败")
	}

	return err
}
