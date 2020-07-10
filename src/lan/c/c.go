package c

import (
	"../../util"
	"strings"
)
import "../../lan"

//C文件结构体
type CFile struct {
	Path string `json:"path"`
	Line int `json:"line"`
	Headers []string `json:"headers"`
	Methods []CMethod `json:"methods"`
	Structs []CStruct `json:"structs"`
	Enums []CEnum `json:"enums"`
}

//C文件方法
type CMethod struct {
	MethodName string `json:"method_name"`
	Params string `json:"params"`
	StartLine int `json:"start_line"`
	EndLine int `json:"end_line"`
}

//C文件结构体
type CStruct struct {
	StructName string `json:"struct_name"`
	StartLine int `json:"start_line"`
	EndLine int `json:"end_line"`
}

//C文件枚举
type CEnum struct {
	EnumName string `json:"enum_name"`
	StartLine int `json:"start_line"`
	EndLine int `json:"end_line"`
}

//C文件接口
type ICFile interface {
	Init()
	Detect(path string)
}

//初始化
func (file *CFile) Init(){
	file.Path = ""
	file.Line = 0
	file.Headers = make([]string, 0)
	file.Methods = make([]CMethod, 0)
	file.Structs = make([]CStruct, 0)
	file.Enums = make([]CEnum, 0)
}

//检测
func (file *CFile) Detect(path string){
	file.Path = path
	content := util.ReadBytes(path)
	size := len(content)
	idx := 0
	line := 0
	for {
		if idx >= size{
			break
		}
		char := string(content[idx])
		checkLine(char, &line)
		if char  == "#"{
			if isInclude(content, &idx, size, &line){
				var header string
				processHeader(content, &idx, size, &line, &header)
				file.Headers = append(file.Headers, header)
			}
		}
		if char == "{"{
			words := getFrontTwoWords(content, idx)
			if words[0] == lan.ENUM {
				ce := CEnum{EnumName:words[1]}
				processEnum(content, &idx, size, &line ,&ce)
				file.Enums = append(file.Enums, ce)
			}else if words[1] == lan.ENUM {
				ce := CEnum{EnumName:""}
				processEnum(content, &idx, size, &line ,&ce)
				file.Enums = append(file.Enums, ce)
			}else if words[0] == lan.STRUCT{
				cs := CStruct{StructName:words[1]}
				processStruct(content, &idx, size, &line, &cs)
				file.Structs = append(file.Structs, cs)
			}else if words[1] == lan.STRUCT {
				cs := CStruct{StructName:""}
				processStruct(content, &idx, size, &line, &cs)
				file.Structs = append(file.Structs, cs)
			}else if strings.Contains(words[0], ")") || strings.Contains(words[1], ")"){
				cm := CMethod{}
				processMethod(content, &idx, size, &line, &cm)
				file.Methods = append(file.Methods, cm)
			}
		}
		idx++
	}
	file.Line = line
}

//判断#后面是否为include
func isInclude(content []byte, idx *int, size int, line *int) bool{
	s := ""
	var char string
	for {
		*idx++
		if *idx >= size{
			break
		}
		char = string(content[*idx])
		checkLine(char, line)
		if char == " "{
			break
		}
		s += char
	}
	return s == lan.INCLUDE
}

//记录头文件
func processHeader(content []byte, idx *int, size int, line *int, header *string){
	var char string
	for {
		*idx++
		if *idx >= size{
			break
		}
		char = string(content[*idx])
		checkLine(char, line)
		if char == "<" || char == "\""{
			break
		}
	}
	s := ""
	for {
		*idx++
		if *idx >= size{
			break
		}
		char = string(content[*idx])
		checkLine(char, line)
		if char == ">" || char == "\""{
			break
		}
		s += char
	}
	*header = s
}

//获取左大括号前面两个标识符，用于判断枚举变量或结构体
func getFrontTwoWords(content []byte, idx int) []string{
	tmpIndex := idx
	var char string
	//跳过两个标识符
	s := make([]string, 2)
	for i:=0;i<2; i++ {
		for {
			tmpIndex--
			if tmpIndex < 0{
				break
			}
			char = string(content[tmpIndex])
			if !util.IsSpace(char){
				break
			}
		}
		tmpS := make([]string, 0)
		for {
			if util.IsSpace(char){
				break
			}
			tmpS = append([]string{char}, tmpS...)
			tmpIndex--
			if tmpIndex < 0{
				break
			}
			char = string(content[tmpIndex])
		}
		s[1-i] = strings.Join(tmpS, "")
	}
	return s
}

//判断是否需要增加行数
func checkLine(char string, line *int){
	if char == "\n"{
		*line++
	}
}

//处理枚举变量
func processEnum(content []byte, idx *int, size int, line *int, ce *CEnum){
	leftBracketCnt := 1
	ce.StartLine = *line + 1
	var char string
	for {
		if leftBracketCnt == 0{
			break
		}
		*idx++
		if *idx >= size{
			break
		}
		char = string(content[*idx])
		checkLine(char, line)
		if char == "{"{
			leftBracketCnt++
		}else if char == "}"{
			leftBracketCnt--
		}
	}
	ce.EndLine = *line + 1
}

//处理结构体
func processStruct(content []byte, idx *int, size int, line *int, cs *CStruct){
	leftBracketCnt := 1
	cs.StartLine = *line + 1
	var char string
	for {
		if leftBracketCnt == 0{
			break
		}
		*idx++
		if *idx >= size{
			break
		}
		char = string(content[*idx])
		checkLine(char, line)
		if char == "{"{
			leftBracketCnt++
		}else if char == "}"{
			leftBracketCnt--
		}
	}
	cs.EndLine = *line + 1
}

//处理函数
func processMethod(content []byte, idx *int, size int, line *int, cm *CMethod){
	cm.StartLine = *line + 1
	var char string
	tmpIndex := *idx
	for {
		tmpIndex--
		if tmpIndex < 0 {
			break
		}
		char = string(content[tmpIndex])
		if char == ")"{
			break
		}
	}
	//查找参数列表
	params := make([]string, 0)
	for {
		params = append([]string{char}, params...)
		if char == "("{
			break
		}
		tmpIndex--
		if tmpIndex < 0 {
			break
		}
		char = string(content[tmpIndex])
	}
	cm.Params = strings.Join(params, "")

	//查找方法名
	for {
		tmpIndex--
		if tmpIndex < 0{
			break
		}
		char = string(content[tmpIndex])
		if !util.IsSpace(char){
			break
		}
	}
	methodName := make([]string, 0)
	for {
		if util.IsSpace(char){
			break
		}
		methodName = append([]string{char}, methodName...)
		tmpIndex--
		if tmpIndex < 0{
			break
		}
		char = string(content[tmpIndex])
	}
	cm.MethodName = strings.Join(methodName, "")

	//查找结束行
	leftBracketCnt := 1
	for {
		if leftBracketCnt == 0{
			break
		}
		*idx++
		if *idx >= size{
			break
		}
		char = string(content[*idx])
		checkLine(char, line)
		if char == "{"{
			leftBracketCnt++
		}else if char == "}"{
			leftBracketCnt--
		}
	}
	cm.EndLine = *line + 1
}