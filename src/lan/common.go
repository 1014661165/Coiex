package lan

import "../util"
import "strings"

const LANGUAGE_C = "c"
const C_INCLUDE = "include"
const C_ENUM = "enum"
const C_STRUCT = "struct"
const C_DEFINE = "define"
const C_KEYWORDS_WITH_PAREN = "for,if,return,switch,while"

const LANGUAGE_JAVA = "java"
const JAVA_PACKAGE = "package"
const JAVA_IMPORT = "import"
const JAVA_CLASS = "class"
const JAVA_ACCESSES = "public,private,protected"
const JAVA_ACCESS_DEFAULT = "default"
const JAVA_EXTENDS = "extends"
const JAVA_IMPLEMENTS = "implements"
const JAVA_THROWS = "throws"
const JAVA_KEYWORDS_WITH_PAREN = "for,switch,catch,if,while"

const LANGUAGE_CPP = "cpp"
const CPP_CLASS = "class"
const CPP_NAMESPACE = "namespace"
const CPP_ACCESS = "public,private,protected"


type File struct {
	Path string `json:"path"`
	Line int `json:"line"`
}

type IFile interface {
	Init()
	Detect(path string)
}

//处理注释
func ProcessComment1(content []byte, idx *int, size int, line *int){
	var char string
	for {
		*idx++
		if *idx >= size{
			break
		}
		char = string(content[*idx])
		CheckLine(char, line)
		if char == "\n"{
			break
		}
	}
}
func ProcessComment2(content []byte, idx *int, size int, line *int){
	var char string
	for {
		*idx++
		if *idx >= size{
			break
		}
		char = string(content[*idx])
		CheckLine(char, line)
		if char == "*"{
			if *idx+1 >= size{
				break
			}
			char = string(content[*idx+1])
			if char == "/"{
				*idx++
				break
			}
		}
	}
}

//判断是否需要增加行数
func CheckLine(char string, line *int){
	if char == "\n"{
		*line++
	}
}

//处理字符串变量
func ProcessString(content []byte, idx *int, size int){
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


//获取前面n个标识符
func GetFrontWords(content []byte, idx int, n int) []string{
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


//从方法体中查找api调用
func FindApis(chars []byte, filterKeywords string) []string{
	var char string
	index := 0
	apis := make([]string, 0)
	for {
		if index >= len(chars){
			break
		}
		char = string(chars[index])
		if char == "\""{
			ProcessString(chars, &index, len(chars))
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
					if !strings.Contains(filterKeywords, api){
						apis = append(apis, api)
					}
				}
			}
		}
		index++
	}
	return apis
}

//判断是否为目标词
func IsTargetWord(content []byte, idx int, target string) bool{
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