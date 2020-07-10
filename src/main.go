package main

import (
	"./config"
	"./lan"
	"./lan/c"
	"./util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

func main() {
	timeStart := time.Now()
	conf:= config.Config{}
	conf.Load(config.CONFIG_FILE)
	files := util.ListFiles(conf.RepoPath, strings.Split(conf.Extensions, ","))
	util.Mkdir(conf.OutputFolder)
	if conf.Language == lan.LANGUAGE_C{
		result := c.Process(files)
		if conf.OutputInterval == 0{
			filename := fmt.Sprintf("%s/1.json", conf.OutputFolder)
			content,_ := json.Marshal(result)
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
				content,_ := json.Marshal(result[index: end])
				ioutil.WriteFile(filename, content, 0744)
				index = end
				fileIndex++
			}
		}

	}
	timeEnd := time.Now()
	fmt.Printf("%d file had been processed!\n", len(files))
	fmt.Printf("task finish! time cost:%.1f s\n", timeEnd.Sub(timeStart).Seconds())
}
