package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//读取文件的字节流
func ReadBytes(path string) []byte{
	data,_ := ioutil.ReadFile(path)
	return data
}

//递归遍历目录
func ListFiles(folder string, extensions []string) []string{
	files := make([]string, 0)
	fs,_ := ioutil.ReadDir(folder)
	for _,f := range fs{
		fullPath := fmt.Sprintf("%s/%s", folder, f.Name())
		if f.IsDir() {
			tmp := ListFiles(fullPath, extensions)
			if len(tmp) != 0{
				files = append(files, tmp...)
			}
		}else {
			if strings.Contains(fullPath, "."){
				index := strings.LastIndex(fullPath, ".")
				extension := fullPath[index+1:]
				exist := false
				for _,ext := range extensions{
					if ext == extension{
						exist = true
						break
					}
				}
				if exist{
					files = append(files, fullPath)
				}
			}
		}
	}
	return files
}

//判断文件是否存在
func Exists(path string) bool{
	_,err := os.Stat(path)
	if err != nil{
		return os.IsExist(err)
	}
	return true
}

//创建文件夹
func Mkdir(folder string){
	_ = os.Mkdir(folder, 0744)
}