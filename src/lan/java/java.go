package java

import (
	"../../lan"
	"../../util"
	"strings"
)

//Java文件
type JavaFile struct {
	lan.File
	Package string `json:"package"`
	Imports []string `json:"imports"`
	Classes []JavaClass `json:"classes"`
	Methods []JavaMethod `json:"methods"`
}

//Java类
type JavaClass struct {
	Name string `json:"name"`
	Access string `json:"access"`
	SuperClass string `json:"super_class"`
	Interfaces []string `json:"interfaces"`
	StartLine int `json:"start_line"`
	EndLine int `json:"end_line"`
}

//Java方法
type JavaMethod struct {
	MethodName string `json:"method_name"`
	Params string `json:"params"`
	Access string `json:"access"`
	StartLine int `json:"start_line"`
	EndLine int `json:"end_line"`
	Apis []string `json:"apis"`
}

//初始化
func (file *JavaFile) Init(){
	file.Path = ""
	file.Line = 0
	file.Package = ""
	file.Imports = make([]string, 0)
	file.Classes = make([]JavaClass, 0)
	file.Methods = make([]JavaMethod, 0)
}

//检测
func (file *JavaFile) Detect(path string){
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
		if char == "p"{
			if isTargetWord(content, idx, lan.JAVA_PACKAGE){
				if file.Package == ""{
					if len(file.Imports) == 0 && len(file.Classes) == 0{
						file.Package = processPackage(content, &idx, size, &line)
					}
				}
			}
		}
		if char == "i"{
			if isTargetWord(content, idx, lan.JAVA_IMPORT){
				if len(file.Classes) == 0{
					importPackage := processImports(content, &idx, size, &line)
					if !strings.Contains(importPackage, "\""){
						file.Imports = append(file.Imports, importPackage)
					}
				}
			}
		}
		if char == "c"{
			if isTargetWord(content, idx, lan.JAVA_CLASS){
				javaClass := JavaClass{}
				processClass(content, &idx, size, &line, &javaClass)
				file.Classes = append(file.Classes, javaClass)
			}
		}
		if char == "{"{
			if isMethod(content, idx){
				javaMethod := JavaMethod{}
				processMethod(content, &idx, size, &line, &javaMethod)
				file.Methods = append(file.Methods, javaMethod)
			}
		}
		idx++
	}
	file.Line = line
}

//判断是否需要增加行数
func checkLine(char string, line *int){
	if char == "\n"{
		*line++
	}
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

//判断是否为目标词
func isTargetWord(content []byte, idx int, target string) bool{
	char := string(content[idx])
	tmpIndex := idx
	var s string
	for {
		s += char
		tmpIndex++
		if tmpIndex >= len(content){
			break
		}
		char = string(content[tmpIndex])
		if !util.IsIdentifier(char){
			break
		}
	}
	preCharIsSpace := true
	nextCharisSpace := true

	if !util.IsSpace(char){
		nextCharisSpace = false
	}
	if idx - 1 >= 0{
		char = string(content[idx - 1])
		if !util.IsSpace(char){
			preCharIsSpace = false
		}
	}
	return  s == target && preCharIsSpace && nextCharisSpace
}

//处理包名
func processPackage(content []byte, idx *int, size int, line *int) string{
	packageName := ""
	*idx += len(lan.JAVA_PACKAGE)
	char := string(content[*idx])
	statement := ""
	for {
		statement += char
		*idx++
		if *idx >= size{
			break
		}
		char = string(content[*idx])
		checkLine(char, line)
		if char == ";"{
			break
		}
	}
	packageName = strings.Trim(statement, " ")
	return packageName
}

//处理导入的包
func processImports(content []byte, idx *int, size int, line *int) string{
	importPackage := ""
	*idx += len(lan.JAVA_IMPORT)
	char := string(content[*idx])
	statement := ""
	for {
		statement += char
		*idx++
		if *idx >= size{
			break
		}
		char = string(content[*idx])
		checkLine(char, line)
		if char == ";"{
			break
		}
	}
	importPackage = strings.Trim(statement, " ")
	return importPackage
}

//处理Java类
func processClass(content []byte, idx *int, size int, line *int, javaClass *JavaClass){
	javaClass.StartLine = *line + 1

	//获取类访问权限
	words := getFrontWords(content, *idx, 3)
	for _,word := range words{
		if strings.Contains(lan.JAVA_ACCESSES, word){
			javaClass.Access = word
			break
		}
	}
	if javaClass.Access == ""{
		javaClass.Access = lan.JAVA_ACCESS_DEFAULT
	}

	//获取类名
	*idx += len(lan.JAVA_CLASS)
	var char string
	for {
		if *idx >= size{
			break
		}
		char = string(content[*idx])
		checkLine(char, line)
		if !util.IsSpace(char){
			break
		}
		*idx++
	}
	for {
		javaClass.Name += char
		*idx++
		if *idx >= size{
			break
		}
		char = string(content[*idx])
		checkLine(char, line)
		if !util.IsIdentifier(char){
			break
		}
	}

	//获取超类和接口
	javaClass.SuperClass = "java.lang.Object"
	javaClass.Interfaces = make([]string, 0)
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
		checkLine(char, line)
	}
	sentence = strings.ReplaceAll(sentence, ",", " ")
	words = strings.Split(sentence, " ")
	extendsAppear := false
	implementsAppear := false
	for _,word := range words {
		word = strings.Trim(word, " ")
		if util.IsIdentifier(word){
			if word == lan.JAVA_EXTENDS {
				extendsAppear = true
			}else if word == lan.JAVA_IMPLEMENTS {
				implementsAppear = true
			}else{
				if extendsAppear && !implementsAppear{
					javaClass.SuperClass = word
				}else if implementsAppear {
					javaClass.Interfaces = append(javaClass.Interfaces, word)
				}
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
		checkLine(char, &tmpLine)
		if char == "{"{
			leftBracketCnt++
		}else if char == "}"{
			leftBracketCnt--
		}
	}
	javaClass.EndLine = tmpLine + 1
}

//获取前面n个标识符
func getFrontWords(content []byte, idx int, n int) []string{
	tmpIndex := idx
	var char string
	//向前跳过n个标识符
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
				tmpS = append([]string{char}, tmpS...)
			}
			tmpIndex--
			if tmpIndex < 0{
				break
			}
			char = string(content[tmpIndex])
		}
		s[n-1-i] = strings.Join(tmpS, "")
	}
	return s
}

func isMethod(content []byte, idx int) bool{
	tmpIndex := idx
	char := string(content[tmpIndex])
	identifiers := make([]string, 0)
	s := ""
	for {
		if char == ")"{
			break
		}
		tmpIndex--
		if tmpIndex < 0{
			break
		}
		char = string(content[tmpIndex])
		if util.IsIdentifier(char){
			s += char
		}else if util.IsSpace(char) && s != ""{
			identifiers = append(identifiers, util.ReverseString(s))
			s = ""
		}
	}
	throwsAppear := false
	for _,id := range identifiers{
		if id == lan.JAVA_THROWS{
			throwsAppear = true
			break
		}
	}
	if throwsAppear{
		return true
	}
	if !throwsAppear && len(identifiers) > 0{
		return false
	}

	//判断左括号前面的字符是否为java关键字
	rightParen := 1
	for {
		if rightParen == 0{
			break
		}
		tmpIndex--
		if tmpIndex < 0{
			break
		}
		char = string(content[tmpIndex])
		if char == ")"{
			rightParen++
		}else if char == "("{
			rightParen--
		}
	}
	word := getFrontWords(content, tmpIndex, 1)
	if strings.Contains(lan.JAVA_KEYWORDS_WITH_PAREN, word[0]){
		return false
	}
	return true
}

//处理Java方法
func processMethod(content []byte, idx *int, size int, line *int, javaMethod *JavaMethod){
	javaMethod.StartLine = *line + 1

	//提取参数列表
	char := string(content[*idx])
	tmpIndex := *idx
	for {
		if char == ")"{
			break
		}
		tmpIndex--
		if tmpIndex < 0{
			break
		}
		char = string(content[tmpIndex])
	}
	javaMethod.Params = ""
	for {
		javaMethod.Params += char
		if char == "("{
			break
		}
		tmpIndex--
		if tmpIndex < 0{
			break
		}
		char = string(content[tmpIndex])
	}
	javaMethod.Params = util.ReverseString(javaMethod.Params)

	//提取方法名和访问权限
	words := getFrontWords(content, tmpIndex, 4)
	javaMethod.MethodName = words[3]
	for _,word := range words{
		if strings.Contains(lan.JAVA_ACCESSES, word){
			javaMethod.Access = word
			break
		}
	}
	if javaMethod.Access == ""{
		javaMethod.Access = lan.JAVA_ACCESS_DEFAULT
	}

	//计算结束行
	leftBracketCnt := 1
	methodBody := []string{string(content[*idx])}
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
		methodBody = append(methodBody, char)
		if char == "{"{
			leftBracketCnt++
		}else if char == "}"{
			leftBracketCnt--
		}
	}
	javaMethod.EndLine = *line + 1
	javaMethod.Apis = findApis(methodBody)
}

//从方法体中查找api调用
func findApis(chars []string) []string{
	var char string
	index := 0
	apis := make([]string, 0)
	for {
		if index >= len(chars){
			break
		}
		char = chars[index]
		if char == "("{
			tmpIndex := index
			var tmpChar string
			apiName := make([]string, 0)
			for {
				tmpIndex--
				if tmpIndex < 0{
					break
				}
				tmpChar = chars[tmpIndex]
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
				tmpChar = chars[tmpIndex]
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