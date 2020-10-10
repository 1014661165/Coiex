package main

import (
	"Coiex/config"
	"Coiex/lan"
	"Coiex/lan/c"
	"Coiex/lan/cpp"
	"Coiex/lan/java"
	"Coiex/util"
	"fmt"
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
		c.Output(result, &conf)
	}else if conf.Language == lan.LANGUAGE_JAVA{
		result := java.Process(files)
		java.Output(result, &conf)
	}else if conf.Language == lan.LANGUAGE_CPP{
		result := cpp.Process(files)
		cpp.Output(result, &conf)
	}
	timeEnd := time.Now()
	fmt.Printf("%d file had been processed!\n", len(files))
	fmt.Printf("task finish! time cost:%.1f s\n", timeEnd.Sub(timeStart).Seconds())
}
