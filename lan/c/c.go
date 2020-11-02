package c

import (
	"Coiex/lan"
	"Coiex/util"
	"strings"
)

var (
	methodId = 0
)

//C文件
type CFile struct {
	lan.File
	Headers []string `json:"headers"`
	Methods []CMethod `json:"methods"`
	Structs []CStruct `json:"structs"`
	Enums []CEnum `json:"enums"`
}

//C文件方法
type CMethod struct {
	MethodId int `json:"method_id"`
	MethodName string `json:"method_name"`
	Params string `json:"params"`
	StartLine int `json:"start_line"`
	EndLine int `json:"end_line"`
	Apis []string `json:"apis"`
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
		lan.CheckLine(char, &line)
		if char == "/"{
			if idx + 1>= size{
				break
			}
			nextChar := string(content[idx + 1])
			if nextChar == "/"{
				lan.ProcessComment1(content, &idx, size, &line)
			}else if nextChar == "*"{
				lan.ProcessComment2(content, &idx, size, &line)
			}
		}
		if char == "\""{
			lan.ProcessString(content, &idx, size)
		}
		if char  == "#"{
			if isInclude(content, &idx, size, &line){
				var header string
				processHeader(content, &idx, size, &line, &header)
				file.Headers = append(file.Headers, header)
			}
		}
		if char == "{"{
			words := lan.GetFrontWords(content, idx, 2)
			if words[0] == lan.C_ENUM {
				ce := CEnum{EnumName:words[1]}
				processEnum(content, &idx, size, &line ,&ce)
				file.Enums = append(file.Enums, ce)
			}else if words[1] == lan.C_ENUM {
				ce := CEnum{EnumName:""}
				processEnum(content, &idx, size, &line ,&ce)
				file.Enums = append(file.Enums, ce)
			}else if words[0] == lan.C_STRUCT{
				cs := CStruct{StructName:words[1]}
				processStruct(content, &idx, size, &line, &cs)
				file.Structs = append(file.Structs, cs)
			}else if words[1] == lan.C_STRUCT {
				cs := CStruct{StructName:""}
				processStruct(content, &idx, size, &line, &cs)
				file.Structs = append(file.Structs, cs)
			}else if strings.Contains(words[0], ")") || strings.Contains(words[1], ")"){
				cm := CMethod{}
				processMethod(content, &idx, size, &line, &cm)
				if cm.MethodName != lan.C_DEFINE{
					methodId++
					cm.MethodId = methodId
					file.Methods = append(file.Methods, cm)
				}
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
		lan.CheckLine(char, line)
		if char == " "{
			break
		}
		s += char
	}
	return s == lan.C_INCLUDE
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
		lan.CheckLine(char, line)
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
		lan.CheckLine(char, line)
		if char == ">" || char == "\""{
			break
		}
		s += char
	}
	*header = s
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
		lan.CheckLine(char, line)
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
		lan.CheckLine(char, line)
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
	rightParen := 1
	for {
		params = append([]string{char}, params...)
		if rightParen == 0{
			break
		}
		tmpIndex--
		if tmpIndex < 0 {
			break
		}
		char = string(content[tmpIndex])
		if char == ")"{
			rightParen++
		}else if char == "("{
			rightParen--
		}
	}
	cm.Params = strings.Join(params, "")

	//查找方法名
	words := lan.GetFrontWords(content, tmpIndex, 2)
	if strings.Contains(words[0], lan.C_DEFINE){
		cm.MethodName = lan.C_DEFINE
		return
	}
	cm.MethodName = words[1]

	//查找结束行
	leftBracketCnt := 1
	methodBody := []byte{content[*idx]}
	for {
		if leftBracketCnt == 0{
			break
		}
		*idx++
		if *idx >= size{
			break
		}
		char = string(content[*idx])
		lan.CheckLine(char, line)
		methodBody = append(methodBody, content[*idx])
		if char == "{"{
			leftBracketCnt++
		}else if char == "}"{
			leftBracketCnt--
		}
	}
	cm.EndLine = *line + 1
	cm.Apis = lan.FindApis(methodBody, lan.C_KEYWORDS_WITH_PAREN)
}