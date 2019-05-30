package parser

import (
	"errors"

	"github.com/kotoyuuko/hina-scanner-parser/scanner"
	"github.com/kotoyuuko/hina-scanner-parser/token"
)

// Parser 语法分析器
type Parser struct {
	tokens  []scanner.TokenItem
	current int
	length  int
}

// New 创建语法分析器
func (p *Parser) New(tokens []scanner.TokenItem) {
	p.tokens = tokens
	p.current = 0
	p.length = len(tokens)
}

// Key 返回当前 Key
func (p *Parser) Key() int {
	return p.current
}

// Token 返回当前 Token
func (p *Parser) Token() scanner.TokenItem {
	return p.tokens[p.current]
}

//Len 返回 Tokens 数组长度
func (p *Parser) Len() int {
	return p.length
}

// Next 分析下一个 Token
func (p *Parser) Next() scanner.TokenItem {
	p.current++

	if p.current >= p.Len() {
		p.current = p.Len() - 1
	}

	return p.Token()
}

// Prev 返回上一个 Token
func (p *Parser) Prev() {
	p.current--
}

var parser Parser

// Init 初始化语法分析器
func Init(tokens []scanner.TokenItem) {
	parser.New(tokens)
}

// Parse 执行语法分析
func Parse() error {
	for parser.Token().Token != token.EOF {
		err := statement()
		if err != nil {
			return err
		}
	}

	return nil
}

func statement() error {
	if parser.Token().Token == token.LET {
		// 初始化语句
		if parser.Next().Token != token.IDENT {
			return errors.New("illegal variable definition statement")
		}
		if parser.Next().Token != token.ASSIGN {
			return errors.New("illegal variable definition statement")
		}

		var err error
		if parser.Next().Token == token.FUNC {
			err = function()
		} else {
			err = expression()
		}
		if err != nil {
			return err
		}
		parser.Next()
	}

	if parser.Token().Token == token.IDENT {
		// 赋值语句 or 函数调用
		if tok := parser.Next(); tok.Token == token.ASSIGN {
			// 赋值语句
			var err error
			if parser.Next().Token == token.FUNC {
				err = function()
			} else {
				err = expression()
			}
			if err != nil {
				return err
			}
			parser.Next()
		} else if tok.Token == token.LPAREN {
			// 函数调用
			err := functionCall()
			if err != nil {
				return err
			}
		}
	}

	if parser.Token().Token.IsKeyword() {
		parser.Next()
	}

	if parser.Token().Token != token.STATEMENT {
		return errors.New("illegal statement")
	}
	parser.Next()

	return nil
}

func functionCall() error {
	if parser.Next().Token != token.LPAREN {
		return errors.New("illegal function call statement")
	}
	if parser.Next().Token != token.RPAREN {
		return errors.New("illegal function call statement")
	}

	parser.Next()
	return nil
}

func function() error {
	if parser.Token().Token != token.FUNC {
		return errors.New("illegal function definition statement")
	}
	if parser.Next().Token != token.LPAREN {
		return errors.New("illegal function definition statement")
	}
	if parser.Next().Token != token.RPAREN {
		return errors.New("illegal function definition statement")
	}

	parser.Next()
	err := block()
	if err != nil {
		return err
	}

	if parser.Token().Token == token.STATEMENT {
		parser.Next()
		return nil
	}

	return errors.New("illegal function definition statement")
}

func expression() error {
	err := term()
	if err != nil {
		return err
	}

	for parser.Token().Token == token.ADD || parser.Token().Token == token.SUB {
		parser.Next()
		err = term()
		if err != nil {
			return err
		}
	}

	return nil
}

func term() error {
	err := factor()
	if err != nil {
		return err
	}

	for parser.Token().Token == token.MUL || parser.Token().Token == token.QUO {
		parser.Next()
		err = factor()
		if err != nil {
			return err
		}
	}

	return nil
}

func factor() error {
	if parser.Token().Token.IsLiteral() {
		parser.Next()
	} else if parser.Token().Token == token.LPAREN {
		parser.Next()

		err := expression()
		if err != nil {
			return err
		}

		if parser.Token().Token != token.RPAREN {
			return errors.New("illegal usage of brackets")
		}
	} else {
		return errors.New("illegal expression")
	}

	return nil
}

func block() error {
	if parser.Token().Token != token.LBRACE {
		return errors.New("illegal code block")
	}

	parser.Next()
	parser.Next()

	for parser.Token().Token != token.RBRACE {
		err := statement()
		if err != nil {
			return err
		}
	}

	parser.Next()

	return nil
}
