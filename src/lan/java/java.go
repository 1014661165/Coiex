package java

import(
	"../../lan"
)

//Java文件
type JavaFile struct {
	lan.File
	Package string `json:"package"`
	Imports []string `json:"imports"`
	Classes []JavaClass `json:"classes"`
}

//Java类
type JavaClass struct {
	Name string `json:"name"`
	Access string `json:"access"`
	SuperClass string `json:"super_class"`
	Interfaces []string `json:"interfaces"`
	StartLine int `json:"start_line"`
	EndLine int `json:"end_line"`
	ChildClass []JavaClass `json:"child_class"`
	Methods []JavaMethod `json:"methods"`
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
func (jfile *JavaFile) Init(){
	jfile.Path = ""
	jfile.Line = 0
	jfile.Package = ""
	jfile.Imports = make([]string, 0)
	jfile.Classes = make([]JavaClass, 0)
}

//检测
func (jfile *JavaFile) Detect(path string){

}


