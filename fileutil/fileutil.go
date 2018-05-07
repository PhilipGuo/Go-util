// fileutil project fileutil.go
package fileutil

import (
	"os"
)

//
// 判断路径是否存在
//
func IsDirExists(dir string) (bool, error) {
	_, err := os.Stat(dir)
	if nil == err {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, nil
}

//
// 创建路径
//
func CreateDir(dir string) bool {
	if ok, _ := IsDirExists(dir); ok {
		return true
	}
	err := os.Mkdir(dir, os.ModeDir)
	if nil == err {
		return true
	}
	return false
}

//
// 判断文件是否存在
//

//
// 创建文件
//

//
// 读取文件（共享）
//

//
// 读取文件（只读）
//

//
// 写文件
//

//
// 删除文件
//

//
// 复制文件
//

//
// 移动文件
//

//
// 取文件后缀
//

//
// 修改后缀
//
