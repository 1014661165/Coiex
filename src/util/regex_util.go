package util
import "regexp"

//判断是否为空格
func IsSpace(s string) bool{
	p := "\\s"
	result,_ := regexp.MatchString(p, s)
	return result
}
