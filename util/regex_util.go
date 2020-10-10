package util
import "regexp"

//判断是否为空格
func IsSpace(s string) bool{
	p := "[\\s]+"
	result,_ := regexp.MatchString(p, s)
	return result
}

//判断是否为标识符
func IsIdentifier(s string) bool{
	p := "[A-Za-z0-9_]+"
	result,_ := regexp.MatchString(p, s)
	return result
}

//字符串反转
func ReverseString(s string) string{
	res := ""
	for i:=len(s)-1; i>=0; i--{
		res += string(s[i])
	}
	return res
}

