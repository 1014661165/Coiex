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
	if conf.Language == lan.LANGUAGE_C{
		result := c.Process(files)
		content,_ := json.Marshal(result)
		ioutil.WriteFile(conf.OutputFile, content, 0744)
	}
	timeEnd := time.Now()
	fmt.Printf("%d file had been processed!\n", len(files))
	fmt.Printf("task finish! time cost:%.1f s\n", timeEnd.Sub(timeStart).Seconds())
}
