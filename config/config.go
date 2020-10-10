package config

import (
	"Coiex/util"
	"encoding/xml"
	"io/ioutil"
)

const CONFIG_FILE  = "./CoiexConfig.xml"

type Config struct {
	RepoPath string `xml:"repo_path"`
	Extensions string `xml:"extensions"`
	Language string `xml:"language"`
	OutputFolder string `xml:"output_folder"`
	OutputInterval int `xml:"output_interval"`
}

type IConfig interface {
	Save(file string)
	Load(file string)
}

func (c *Config) Save(file string){
	c.RepoPath = ""
	c.Language = ""
	c.Extensions = ""
	c.OutputFolder = "result"
	c.OutputInterval = 0
	content,_ := xml.MarshalIndent(c, "", "  ")
	ioutil.WriteFile(file, content, 0744)
}

func (c *Config) Load(file string){
	exist := util.Exists(file)
	if !exist{
		c.Save(CONFIG_FILE)
		panic("please update config file")
	}
	content,err := ioutil.ReadFile(file)
	if err != nil{
		panic("fail to load config file")
	}
	xml.Unmarshal(content, c)
}