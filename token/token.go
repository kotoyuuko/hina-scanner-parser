package token

import (
	"strconv"
)

// Token 表示 Token 类型
type Token int

// Token 列表
const (
	// 特殊 Token
	ILLEGAL Token = iota
	STATEMENT
	EOF
	COMMENT

	literalBegin
	// 标识符和基本数据类型
	IDENT  // main
	INT    // 12345
	FLOAT  // 123.45
	CHAR   // 'a'
	STRING // "abc"
	literalEnd

	operatorBegin
	// 运算符
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND // &
	OR  // |
	XOR // ^
	SHL // <<
	SHR // >>

	LAND // &&
	LOR  // ||
	INC  // ++
	DEC  // --

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ // !=
	LEQ // <=
	GEQ // >=
	operatorEnd

	delimiterBegin
	// 界符
	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN // )
	RBRACK // ]
	RBRACE // }
	COLON  // :
	delimiterEnd

	keywordBegin
	// 关键字
	BREAK
	CONST
	CONTINUE
	ELSE
	FOR
	IF
	RETURN
	LET
	IMPORT
	FUNC
	keywordEnd
)

var tokens = [...]string{
	ILLEGAL:   "ILLEGAL",
	STATEMENT: "STATEMENT",
	EOF:       "EOF",
	COMMENT:   "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	CHAR:   "CHAR",
	STRING: "STRING",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	AND: "&",
	OR:  "|",
	XOR: "^",
	SHL: "<<",
	SHR: ">>",

	LAND: "&&",
	LOR:  "||",
	INC:  "++",
	DEC:  "--",

	EQL:    "==",
	LSS:    "<",
	GTR:    ">",
	ASSIGN: "=",
	NOT:    "!",

	NEQ: "!=",
	LEQ: "<=",
	GEQ: ">=",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",
	PERIOD: ".",

	RPAREN: ")",
	RBRACK: "]",
	RBRACE: "}",
	COLON:  ":",

	BREAK:    "break",
	CONST:    "const",
	CONTINUE: "continue",
	ELSE:     "else",
	FOR:      "for",
	IF:       "if",
	RETURN:   "return",
	LET:      "let",
	IMPORT:   "import",
	FUNC:     "func",
}

// String 返回 Token 的字符串描述
func (t Token) String() string {
	s := ""
	if 0 <= t && t < Token(len(tokens)) {
		s = tokens[t]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}

// 定义优先级
const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 6
	HighestPrec = 7
)

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keywordBegin + 1; i < keywordEnd; i++ {
		keywords[tokens[i]] = i
	}
}

// Lookup 返回关键字或标识符的 Token
func Lookup(ident string) Token {
	if token, isKeyword := keywords[ident]; isKeyword {
		return token
	}
	return IDENT
}

// IsLiteral 判断 Token 是否为标识符或基本数据类型
func (t Token) IsLiteral() bool {
	return literalBegin < t && t < literalEnd
}

// IsOperator 判断 Token 是否为运算符
func (t Token) IsOperator() bool {
	return operatorBegin < t && t < operatorEnd
}

// IsDelimiter 判断 Token 是否为界符
func (t Token) IsDelimiter() bool {
	return delimiterBegin < t && t < delimiterEnd
}

// IsKeyword 判断 Token 是否为关键字
func (t Token) IsKeyword() bool {
	return keywordBegin < t && t < keywordEnd
}

// IsKeyword 判断字符串是否为关键字
func IsKeyword(name string) bool {
	_, ok := keywords[name]
	return ok
}
