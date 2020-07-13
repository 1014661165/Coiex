package lan

const LANGUAGE_C = "c"
const LANGUAGE_JAVA = "java"
const INCLUDE = "include"
const ENUM = "enum"
const STRUCT = "struct"
const DEFINE = "define"
const C_KEYWORDS_WITH_PAREN = "for,if,return,switch,while"


type File struct {
	Path string `json:"path"`
	Line int `json:"line"`
}

type IFile interface {
	Init()
	Detect(path string)
}