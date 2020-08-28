package cpp

import (
	"../../lan"
	"../../util"
	"strings"
)

//Cpp文件
type CppFile struct {
	lan.File
	Headers []string `json:"headers"`
	Namespaces []CppNamespace `json:"namespaces"`
	Classes []CppClass `json:"classes"`
	Methods []CppMethod `json:"methods"`
	Structs []CppStruct `json:"structs"`
	Enums []CppEnum `json:"enums"`
}

type CppNamespace struct {
	Name string `json:"name"`
	StartLine int `json:"start_line"`
	EndLine int `json:"end_line"`
}

type CppClass struct {
	Name string `json:"name"`
	SuperClasses []string `json:"super_classes"`
	StartLine int `json:"start_line"`
	EndLine int `json:"end_line"`
}

//Cpp文件方法
type CppMethod struct {
	MethodName string `json:"method_name"`
	Params string `json:"params"`
	StartLine int `json:"start_line"`
	EndLine int `json:"end_line"`
	Apis []string `json:"apis"`
}

//Cpp文件结构体
type CppStruct struct {
	StructName string `json:"struct_name"`
	StartLine int `json:"start_line"`
	EndLine int `json:"end_line"`
}

//Cpp文件枚举
type CppEnum struct {
	EnumName string `json:"enum_name"`
	StartLine int `json:"start_line"`
	EndLine int `json:"end_line"`
}

//初始化
func (file *CppFile) Init(){
	file.Path = ""
	file.Line = 0
	file.Headers = make([]string, 0)
	file.Namespaces = make([]CppNamespace, 0)
	file.Classes = make([]CppClass, 0)
	file.Methods = make([]CppMethod, 0)
	file.Structs = make([]CppStruct, 0)
	file.Enums = make([]CppEnum, 0)
}

//初始化
func (file *CppFile) Detect(path string){
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
		if char == "c"{
			if lan.IsTargetWord(content, idx, lan.CPP_CLASS){
				cppClass := CppClass{}
				processClass(content, &idx, size, &line, &cppClass)
				file.Classes = append(file.Classes, cppClass)
			}
		}
		if char == "{"{
			words := lan.GetFrontWords(content, idx, 2)
			if words[0] == lan.C_ENUM {
				ce := CppEnum{EnumName:words[1]}
				processEnum(content, &idx, size, &line ,&ce)
				file.Enums = append(file.Enums, ce)
			}else if words[1] == lan.C_ENUM {
				ce := CppEnum{EnumName:""}
				processEnum(content, &idx, size, &line ,&ce)
				file.Enums = append(file.Enums, ce)
			}else if words[0] == lan.C_STRUCT{
				cs := CppStruct{StructName:words[1]}
				processStruct(content, &idx, size, &line, &cs)
				file.Structs = append(file.Structs, cs)
			}else if words[1] == lan.C_STRUCT {
				cs := CppStruct{StructName:""}
				processStruct(content, &idx, size, &line, &cs)
				file.Structs = append(file.Structs, cs)
			}else if words[0] == lan.CPP_NAMESPACE {
				cn := CppNamespace{Name:words[1]}
				processNamespace(content, &idx, size, &line, &cn)
				file.Namespaces = append(file.Namespaces, cn)
			} else if strings.Contains(words[0], ")") || strings.Contains(words[1], ")"){
				cm := CppMethod{}
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

//处理Cpp类
func processClass(content []byte, idx *int, size int, line *int, cppClass *CppClass){
	cppClass.StartLine = *line + 1

	//获取类名
	*idx += len(lan.CPP_CLASS)
	var char string
	for {
		if *idx >= size{
			break
		}
		char = string(content[*idx])
		lan.CheckLine(char, line)
		if !util.IsSpace(char){
			break
		}
		*idx++
	}
	for {
		cppClass.Name += char
		*idx++
		if *idx >= size{
			break
		}
		char = string(content[*idx])
		lan.CheckLine(char, line)
		if !util.IsIdentifier(char){
			break
		}
	}

	//获取超类
	cppClass.SuperClasses = make([]string, 0)
	sentence := ""
	for {
		if char == "{"{
			break
		}
		sentence += char
		*idx++
		if *idx >= size{
			break
		}
		char = string(content[*idx])
		lan.CheckLine(char, line)
	}
	sentence = strings.ReplaceAll(sentence, ",", " ")
	words := strings.Split(sentence, " ")
	for _,word := range words {
		word = strings.Trim(word, " :\n")
		if util.IsIdentifier(word){
			if !strings.Contains(lan.CPP_ACCESS, word){
				cppClass.SuperClasses = append(cppClass.SuperClasses, word)
			}
		}
	}

	//获取结束行
	leftBracketCnt := 1
	tmpIndex := *idx
	tmpLine := *line
	for {
		if leftBracketCnt == 0{
			break
		}
		tmpIndex++
		if tmpIndex >= size{
			break
		}
		char = string(content[tmpIndex])
		lan.CheckLine(char, &tmpLine)
		if char == "{"{
			leftBracketCnt++
		}else if char == "}"{
			leftBracketCnt--
		}
	}
	cppClass.EndLine = tmpLine + 1
}

//处理枚举变量
func processEnum(content []byte, idx *int, size int, line *int, ce *CppEnum){
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
func processStruct(content []byte, idx *int, size int, line *int, cs *CppStruct){
	leftBracketCnt := 1
	cs.StartLine = *line + 1
	var char string
	tmpIndex := *idx
	tmpLine := *line
	for {
		if leftBracketCnt == 0{
			break
		}
		tmpIndex++
		if tmpIndex >= size{
			break
		}
		char = string(content[tmpIndex])
		lan.CheckLine(char, &tmpLine)
		if char == "{"{
			leftBracketCnt++
		}else if char == "}"{
			leftBracketCnt--
		}
	}
	cs.EndLine = tmpLine + 1
}

//处理命名空间
func processNamespace(content []byte, idx *int, size int, line *int, cn *CppNamespace){
	cn.StartLine = *line + 1
	char := string(content[*idx])

	//获取结束行
	leftBracketCnt := 1
	tmpIndex := *idx
	tmpLine := *line
	for {
		if leftBracketCnt == 0{
			break
		}
		tmpIndex++
		if tmpIndex >= size{
			break
		}
		char = string(content[tmpIndex])
		lan.CheckLine(char, &tmpLine)
		if char == "{"{
			leftBracketCnt++
		}else if char == "}"{
			leftBracketCnt--
		}
	}
	cn.EndLine = tmpLine + 1
}

//处理函数
func processMethod(content []byte, idx *int, size int, line *int, cm *CppMethod){
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