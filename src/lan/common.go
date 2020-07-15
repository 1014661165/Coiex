package lan

const LANGUAGE_C = "c"
const LANGUAGE_JAVA = "java"
const C_INCLUDE = "include"
const C_ENUM = "enum"
const C_STRUCT = "struct"
const C_DEFINE = "define"
const C_KEYWORDS_WITH_PAREN = "for,if,return,switch,while"
const JAVA_PACKAGE = "package"
const JAVA_IMPORT = "import"
const JAVA_CLASS = "class"
const JAVA_ACCESS_PUBLIC = "public"
const JAVA_ACCESS_PRIVATE = "private"
const JAVA_ACCESS_PROTECTED = "protected"
const JAVA_ACCESSES = "public,private,protected"
const JAVA_ACCESS_DEFAULT = "default"
const JAVA_EXTENDS = "extends"
const JAVA_IMPLEMENTS = "implements"
const JAVA_THROWS = "throws"
const JAVA_KEYWORDS_WITH_PAREN = "for,switch,catch,if,while"

type File struct {
	Path string `json:"path"`
	Line int `json:"line"`
}

type IFile interface {
	Init()
	Detect(path string)
}