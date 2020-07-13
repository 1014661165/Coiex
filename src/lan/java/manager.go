package java

import (
	"../../config"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//处理程序
func Process(files []string) []JavaFile{
	result := make([]JavaFile, len(files))
	for i, file := range files {
		fmt.Printf("%.2f%%\n", float32((i+1)*100)/float32(len(files)))
		jfile := JavaFile{}
		jfile.Init()
		jfile.Detect(file)
		result[i] = jfile
	}
	return result
}

//输出
func Output(result []JavaFile, conf *config.Config){
	if conf.OutputInterval == 0{
		filename := fmt.Sprintf("%s/1.json", conf.OutputFolder)
		content,_ := json.MarshalIndent(result, "", "\t")
		ioutil.WriteFile(filename, content, 0744)
	}else if conf.OutputInterval > 0 {
		index := 0
		fileIndex := 1
		for {
			if index >= len(result){
				break
			}
			filename := fmt.Sprintf("%s/%d.json", conf.OutputFolder, fileIndex)
			var end int
			if index + conf.OutputInterval > len(result){
				end = len(result)
			}else{
				end = index + conf.OutputInterval
			}
			content,_ := json.MarshalIndent(result[index: end], "", "\t")
			ioutil.WriteFile(filename, content, 0744)
			index = end
			fileIndex++
		}
	}
}
