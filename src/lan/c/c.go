package c

import (
	"../../lan"
	"../../util"
	"strings"
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
		checkLine(char, &line)
		if char == "/"{
			if idx + 1>= size{
				break
			}
			nextChar := string(content[idx + 1])
			if nextChar == "/"{
				processComment1(content, &idx, size, &line)
			}else if nextChar == "*"{
				processComment2(content, &idx, size, &line)
			}
		}
		if char == "\""{
			processString(content, &idx, size)
		}
		if char  == "#"{
			if isInclude(content, &idx, size, &line){
				var header string
				processHeader(content, &idx, size, &line, &header)
				file.Headers = append(file.Headers, header)
			}
		}
		if char == "{"{
			words := getFrontWords(content, idx, 2)
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
					file.Methods = append(file.Methods, cm)
				}
			}
		}
		idx++
	}
	file.Line = line
}

//处理注释
func processComment1(content []byte, idx *int, size int, line *int){
	var char string
	for {
		*idx++
		if *idx >= size{
			break
		}
		char = string(content[*idx])
		checkLine(char, line)
		if char == "\n"{
			break
		}
	}
}
func processComment2(content []byte, idx *int, size int, line *int){
	var char string
	for {
		*idx++
		if *idx >= size{
			break
		}
		char = string(content[*idx])
		checkLine(char, line)
		if char == "*"{
			*idx++
			if *idx >= size{
				break
			}
			char = string(content[*idx])
			checkLine(char, line)
			if char == "/"{
				break
			}
		}
	}
}

//处理字符串变量
func processString(content []byte, idx *int, size int){
	char := string(content[*idx])
	for {
		*idx++
		if *idx >= size{
			break
		}
		if content[*idx] == '\\'{
			*idx++
			continue
		}
		char = string(content[*idx])
		if char == "\""{
			break
		}
	}
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

//获取前面n个标识符
func getFrontWords(content []byte, idx int, n int) []string{
	tmpIndex := idx
	var char string
	//跳过n个标识符
	s := make([]string, n)
	for i:=0;i<n; i++ {
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
			if char != "*"{
				tmpS = append(tmpS, char)
			}
			tmpIndex--
			if tmpIndex < 0{
				break
			}
			char = string(content[tmpIndex])
		}
		s[n-1-i] = util.ReverseString(strings.Join(tmpS, ""))
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
	words := getFrontWords(content, tmpIndex, 2)
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
		checkLine(char, line)
		methodBody = append(methodBody, content[*idx])
		if char == "{"{
			leftBracketCnt++
		}else if char == "}"{
			leftBracketCnt--
		}
	}
	cm.EndLine = *line + 1
	cm.Apis = findApis(methodBody)
}

//从方法体中查找api调用
func findApis(chars []byte) []string{
	var char string
	index := 0
	apis := make([]string, 0)
	for {
		if index >= len(chars){
			break
		}
		char = string(chars[index])
		if char == "\""{
			processString(chars, &index, len(chars))
		}

		if char == "("{
			tmpIndex := index
			var tmpChar string
			apiName := make([]string, 0)
			for {
				tmpIndex--
				if tmpIndex < 0{
					break
				}
				tmpChar = string(chars[tmpIndex])
				if !util.IsSpace(tmpChar){
					break
				}
			}
			for {
				if !util.IsIdentifier(tmpChar){
					break
				}
				apiName = append([]string{tmpChar}, apiName...)
				tmpIndex--
				if tmpIndex < 0{
					break
				}
				tmpChar = string(chars[tmpIndex])
			}
			if len(apiName) != 0{
				api := strings.Join(apiName, "")
				if api != ""{
					if !strings.Contains(lan.C_KEYWORDS_WITH_PAREN, api){
						apis = append(apis, api)
					}
				}
			}
		}
		index++
	}
	return apis
}