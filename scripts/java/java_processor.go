package main

import (
	"Coiex/lan/java"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

//输出文件
const (
	MethodIndexOutputFile = "MethodIndex.csv"
	CallListOutputFile = "CallList.csv"
)

//方法索引数据，方便查询
var (
	methodIndex map[string]map[string][]Method
)

//方法结构体
type Method struct {
	MethodId int
	MethodName string
	Params string
	StartLine int
	EndLine int
	Path string
	Package string
}

func init(){
	methodIndex = make(map[string]map[string][]Method)
}

//初始化方法索引
func initMethodIndex(inputFolder string){
	fis, err := ioutil.ReadDir(inputFolder)
	if err != nil {
		log.Fatal(err)
	}
	cnt := 0
	size := len(fis)
	for _, file := range fis {
		cnt++
		log.Printf("%.2f%%\n", float64(cnt*100)/float64(size))

		path := inputFolder + "/" + file.Name()
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println(err)
			continue
		}

		//将json文件解析到结构体中
		var javaFiles []java.JavaFile
		err = json.Unmarshal(content, &javaFiles)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, javaFile := range javaFiles {
			path := javaFile.Path
			sep := "/"
			if strings.Contains(path, "\\"){
				sep = "\\"
			}

			//获取目录
			pack := path[0: strings.LastIndex(path, sep)]
			javaMethods := javaFile.Methods
			if len(javaMethods) == 0 {
				continue
			}
			for _, javaMethod := range javaMethods {
				//获取methodIndex索引中javaMethod对应的map
				methodSet, ok := methodIndex[javaMethod.MethodName]
				if !ok {
					methodSet = make(map[string][]Method)
					methodIndex[javaMethod.MethodName] = methodSet
				}

				//获取map中在pack包内的方法
				methodSubSet, ok := methodSet[pack]
				if !ok {
					methodSubSet = make([]Method, 0)
					methodIndex[javaMethod.MethodName][pack] = methodSubSet
				}
				method := Method{
					MethodId:   javaMethod.MethodId,
					MethodName: javaMethod.MethodName,
					Params:     javaMethod.Params,
					StartLine:  javaMethod.StartLine,
					EndLine:    javaMethod.EndLine,
					Path:       javaFile.Path,
					Package:    pack,
				}
				methodIndex[javaMethod.MethodName][pack] = append(methodIndex[javaMethod.MethodName][pack], method)
			}
		}
	}
}

func process(inputFolder string){
	f1, err := os.OpenFile(MethodIndexOutputFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}
	f2, err := os.OpenFile(CallListOutputFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer  f1.Close()
	defer  f2.Close()
	writer1 := bufio.NewWriter(f1)
	writer2 := bufio.NewWriter(f2)

	fis, err := ioutil.ReadDir(inputFolder)
	if err != nil {
		log.Fatal(err)
	}
	cnt := 0
	size := len(fis)
	for _, file := range fis {
		cnt++
		log.Printf("%.2f%%\n", float64(cnt*100)/float64(size))

		path := inputFolder + "/" + file.Name()
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println(err)
			continue
		}
		var javaFiles []java.JavaFile
		err = json.Unmarshal(content, &javaFiles)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, javaFile := range javaFiles {
			path := javaFile.Path
			sep := "/"
			if strings.Contains(path, "\\"){
				sep = "\\"
			}
			pack := path[0: strings.LastIndex(path, sep)]
			javaMethods := javaFile.Methods
			if len(javaMethods) == 0 {
				continue
			}

			for _, javaMethod := range javaMethods {
				writer1.WriteString(fmt.Sprintf("%d,%s,%d,%d\n", javaMethod.MethodId, javaFile.Path, javaMethod.StartLine, javaMethod.EndLine))
				if len(javaMethod.Apis) == 0 {
					continue
				}
				for _, api := range javaMethod.Apis {
					set1, ok := methodIndex[api]
					if !ok {
						continue
					}
					set2, ok := set1[pack]
					if ok {
						for _, m := range set2 {
							if javaMethod.MethodId != m.MethodId{
								writer2.WriteString(fmt.Sprintf("%d,%d\n", javaMethod.MethodId, m.MethodId))
							}
						}
					}else{
						if len(javaFile.Imports) == 0 {
							continue
						}
						p := javaFile.Package
						if strings.Contains(javaFile.Package, "."){
							p = strings.ReplaceAll(javaFile.Package, ".", sep)
						}

						if !strings.Contains(path, p){
							continue
						}
						prefix := path[0: strings.LastIndex(path, p)]

						for _, importPackage := range javaFile.Imports {
							n := strings.ReplaceAll(importPackage, ".", sep)
							n = n[0: strings.LastIndex(n, sep)]

							newPack := prefix + n
							set3, ok := set1[newPack]
							if ok {
								for _, m := range set3 {
									if javaMethod.MethodId != m.MethodId{
										writer2.WriteString(fmt.Sprintf("%d,%d\n", javaMethod.MethodId, m.MethodId))
									}
								}
							}
						}
					}
				}
			}
		}
	}
	writer1.Flush()
	writer2.Flush()
}


func main() {
	if len(os.Args) < 2{
		fmt.Println("go run java_processor.go inputFolder")
		os.Exit(0)
	}
	timeStart := time.Now()
	log.Print("create method index...")
	initMethodIndex(os.Args[1])

	log.Print("processing...")
	process(os.Args[1])

	log.Printf("task finish! time cost: %v\n", time.Since(timeStart))
}
