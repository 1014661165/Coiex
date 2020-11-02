package main

import (
	"Coiex/lan/cpp"
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
		var cppFiles []cpp.CppFile
		err = json.Unmarshal(content, &cppFiles)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, cppFile := range cppFiles {
			path := cppFile.Path
			sep := "/"
			if strings.Contains(path, "\\"){
				sep = "\\"
			}

			//获取目录
			pack := path[0: strings.LastIndex(path, sep)]
			cppMethods := cppFile.Methods
			if len(cppMethods) == 0 {
				continue
			}
			for _, cppMethod := range cppMethods {
				//获取methodIndex索引中cppMethod对应的map
				methodSet, ok := methodIndex[cppMethod.MethodName]
				if !ok {
					methodSet = make(map[string][]Method)
					methodIndex[cppMethod.MethodName] = methodSet
				}

				//获取map中在pack包内的方法
				methodSubSet, ok := methodSet[pack]
				if !ok {
					methodSubSet = make([]Method, 0)
					methodIndex[cppMethod.MethodName][pack] = methodSubSet
				}
				method := Method{
					MethodId:   cppMethod.MethodId,
					MethodName: cppMethod.MethodName,
					Params:     cppMethod.Params,
					StartLine:  cppMethod.StartLine,
					EndLine:    cppMethod.EndLine,
					Path:       cppFile.Path,
					Package:    pack,
				}
				methodIndex[cppMethod.MethodName][pack] = append(methodIndex[cppMethod.MethodName][pack], method)
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

		var cppFiles []cpp.CppFile
		err = json.Unmarshal(content, &cppFiles)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, cppFile := range cppFiles {
			path := cppFile.Path
			sep := "/"
			if strings.Contains(path, "\\"){
				sep = "\\"
			}

			//获取目录
			pack := path[0: strings.LastIndex(path, sep)]
			cppMethods := cppFile.Methods
			if len(cppMethods) == 0 {
				continue
			}
			for _, cppMethod := range cppMethods {
				writer1.WriteString(fmt.Sprintf("%d,%s,%d,%d\n", cppMethod.MethodId, cppFile.Path, cppMethod.StartLine, cppMethod.EndLine))
				if len(cppMethod.Apis) == 0{
					continue
				}

				for _, api := range cppMethod.Apis {
					set1, ok := methodIndex[api]
					if !ok {
						continue
					}
					set2, ok := set1[pack]
					if ok {
						for _, m := range set2 {
							if cppMethod.MethodId != m.MethodId{
								writer2.WriteString(fmt.Sprintf("%d,%d\n", cppMethod.MethodId, m.MethodId))
							}
						}
					}
				}
			}

		}

	}
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
