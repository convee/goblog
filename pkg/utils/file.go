package utils

import "os"

// 判断所给路径文件/文件夹是否存在

func FileExists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
