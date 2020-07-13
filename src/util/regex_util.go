package util
import "regexp"

//判断是否为空格
func IsSpace(s string) bool{
	p := "\\s"
	result,_ := regexp.MatchString(p, s)
	return result
}

//判断是否为标识符
func IsIdentifier(s string) bool{
	p := "[A-Za-z0-9_]"
	result,_ := regexp.MatchString(p, s)
	return result
}
