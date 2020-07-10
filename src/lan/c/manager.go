package c

import "fmt"

//处理程序
func Process(files []string) []CFile{
	result := make([]CFile, len(files))
	for i, file := range files {
		fmt.Printf("%.2f%%\n", float32((i+1)*100)/float32(len(files)))
		cfile := CFile{}
		cfile.Init()
		cfile.Detect(file)
		result[i] = cfile
	}
	return result
}
