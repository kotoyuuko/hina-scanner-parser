package scanner

import (
	"bufio"
	"io"
	"os"

	"github.com/kotoyuuko/hina-scanner-parser/token"
)

// Scanner 读取源文件的 Reader
type Scanner struct {
	File     *os.File
	Start    int
	Position int
	bytes    []byte
	reader   *bufio.Reader
}

// NewScanner 新建 Scanner 对象
func NewScanner(file *os.File) *Scanner {
	scanner := &Scanner{
		File:   file,
		reader: bufio.NewReader(file),
	}
	return scanner
}

// Next 读取下一个字节
func (s *Scanner) Next() byte {
	b, err := s.reader.ReadByte()
	if err == io.EOF {
		return 0
	}
	if err != nil {
		panic(err)
	}

	s.bytes = append(s.bytes, b)
	s.Position++

	return b
}

// Prev 返回上一个字节
func (s *Scanner) Prev() {
	err := s.reader.UnreadByte()
	if err != nil {
		panic(err)
	}

	s.Position--
	s.bytes = s.bytes[0:s.Position]
}

// Ignore 从当前位置重新开始读取
func (s *Scanner) Ignore() {
	s.Start = s.Position
}

// TokenItem 保存从文件中读取的 Token 二元组
type TokenItem struct {
	Token token.Token
	Code  string
}

func isLetter(b byte) bool {
	if b >= 'a' && b <= 'z' {
		return true
	}
	if b >= 'A' && b <= 'Z' {
		return true
	}
	if b == '_' {
		return true
	}
	return false
}

func isDigit(b byte) bool {
	if b >= '0' && b <= '9' {
		return true
	}
	return false
}

func isOperator(b byte) bool {
	if b == '+' || b == '-' || b == '*' || b == '/' || b == '%' {
		return true
	}
	if b == '>' || b == '<' || b == '=' || b == '!' {
		return true
	}
	if b == '&' || b == '|' || b == '^' {
		return true
	}
	return false
}

func isDelimiter(b byte) bool {
	if b == '(' || b == ')' || b == '{' || b == '}' || b == '[' || b == ']' {
		return true
	}
	if b == ',' || b == '.' || b == ':' {
		return true
	}
	return false
}

// Scan 读取源文件并返回对应的 Token 二元组的数组
func (s *Scanner) Scan() []TokenItem {
	tokens := []TokenItem{}
	for {
		b := s.Next()

		if b == 0 {
			break
		}

		item := TokenItem{}
		if b == ' ' || b == '\r' || b == '\t' {
			s.Ignore()
			continue
		} else if b == '\n' {
			item = TokenItem{
				Token: token.STATEMENT,
				Code:  "",
			}
			s.Ignore()
		} else if b == '#' {
			item = s.scanComment()
		} else if b == '"' {
			item = s.scanString()
		} else if isLetter(b) {
			item = s.scanLetter()
		} else if isDigit(b) {
			item = s.scanDigit()
		} else if isOperator(b) {
			item = s.scanOperator(b)
		} else if isDelimiter(b) {
			item = s.scanDelimiter(b)
		} else {
			item = s.scanIllegal()
		}
		tokens = append(tokens, item)
	}
	tokens = append(tokens, TokenItem{
		Token: token.EOF,
		Code:  "",
	})
	return tokens
}

func (s *Scanner) scanComment() TokenItem {
	item := TokenItem{
		Token: token.COMMENT,
	}
	for {
		if b := s.Next(); b == '\n' {
			break
		}
	}
	item.Code = string(s.bytes[s.Start : s.Position-1])
	s.Ignore()
	return item
}

func (s *Scanner) scanIllegal() TokenItem {
	item := TokenItem{
		Token: token.ILLEGAL,
		Code:  string(s.bytes[s.Start:s.Position]),
	}
	s.Ignore()
	return item
}

func (s *Scanner) scanString() TokenItem {
	item := TokenItem{
		Token: token.STRING,
	}
	for {
		if b := s.Next(); b == '"' {
			break
		}
	}
	item.Code = string(s.bytes[s.Start:s.Position])
	s.Ignore()
	return item
}

func (s *Scanner) scanLetter() TokenItem {
	item := TokenItem{}
	for {
		if b := s.Next(); !isLetter(b) && !isDigit(b) {
			break
		}
	}
	s.Prev()
	item.Code = string(s.bytes[s.Start:s.Position])
	item.Token = token.Lookup(item.Code)
	s.Ignore()
	return item
}

func (s *Scanner) scanDigit() TokenItem {
	item := TokenItem{}
	digitType := 1
	for {
		b := s.Next()
		if !isDigit(b) && b != '.' {
			break
		}
		if b == '.' {
			digitType = 2
		}
	}
	s.Prev()
	item.Code = string(s.bytes[s.Start:s.Position])
	item.Token = token.INT
	if digitType == 2 {
		item.Token = token.FLOAT
	}
	s.Ignore()
	return item
}

func (s *Scanner) scanOperator(b byte) TokenItem {
	item := TokenItem{}
	if b == '+' {
		if b = s.Next(); b != '+' {
			s.Prev()
		}
		item.Code = string(s.bytes[s.Start:s.Position])
		if item.Code == "++" {
			item.Token = token.INC
		} else {
			item.Token = token.ADD
		}
	} else if b == '-' {
		if b = s.Next(); b != '-' {
			s.Prev()
		}
		item.Code = string(s.bytes[s.Start:s.Position])
		if item.Code == "--" {
			item.Token = token.DEC
		} else {
			item.Token = token.SUB
		}
	} else if b == '*' {
		item.Code = string(s.bytes[s.Start:s.Position])
		item.Token = token.MUL
	} else if b == '/' {
		item.Code = string(s.bytes[s.Start:s.Position])
		item.Token = token.QUO
	} else if b == '%' {
		item.Code = string(s.bytes[s.Start:s.Position])
		item.Token = token.REM
	} else if b == '>' || b == '<' || b == '=' || b == '!' {
		if bNext := s.Next(); bNext != '=' {
			s.Prev()
		}
		if bNext := s.Next(); !(bNext == '<' && b == '<') {
			s.Prev()
		}
		if bNext := s.Next(); !(bNext == '>' && b == '>') {
			s.Prev()
		}
		item.Code = string(s.bytes[s.Start:s.Position])
		if item.Code == ">" {
			item.Token = token.GTR
		} else if item.Code == "<" {
			item.Token = token.LSS
		} else if item.Code == "=" {
			item.Token = token.ASSIGN
		} else if item.Code == "!" {
			item.Token = token.NOT
		} else if item.Code == ">=" {
			item.Token = token.GEQ
		} else if item.Code == "<=" {
			item.Token = token.LEQ
		} else if item.Code == "==" {
			item.Token = token.EQL
		} else if item.Code == "!=" {
			item.Token = token.NEQ
		} else if item.Code == "<<" {
			item.Token = token.SHL
		} else if item.Code == ">>" {
			item.Token = token.SHR
		}
	} else if b == '&' || b == '|' || b == '^' {
		if bNext := s.Next(); b != bNext {
			s.Prev()
		}
		item.Code = string(s.bytes[s.Start:s.Position])
		if item.Code == "&" {
			item.Token = token.AND
		} else if item.Code == "|" {
			item.Token = token.OR
		} else if item.Code == "^" {
			item.Token = token.XOR
		} else if item.Code == "&&" {
			item.Token = token.LAND
		} else if item.Code == "||" {
			item.Token = token.LOR
		}
	}
	s.Ignore()
	return item
}

func (s *Scanner) scanDelimiter(b byte) TokenItem {
	item := TokenItem{
		Code: string(s.bytes[s.Start:s.Position]),
	}
	if b == '(' {
		item.Token = token.LPAREN
	} else if b == '[' {
		item.Token = token.LBRACK
	} else if b == '{' {
		item.Token = token.LBRACE
	} else if b == ')' {
		item.Token = token.RPAREN
	} else if b == ']' {
		item.Token = token.RBRACK
	} else if b == '}' {
		item.Token = token.RBRACE
	} else if b == ',' {
		item.Token = token.COMMA
	} else if b == '.' {
		item.Token = token.PERIOD
	} else if b == ':' {
		item.Token = token.COLON
	}
	s.Ignore()
	return item
}
