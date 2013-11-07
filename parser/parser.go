package parser

import (
    "fmt"
    "../../clispi/lexer"
    . "../../clispi/globals"
)

type ParserError struct {
    pos Position
    msg string
}

func (e ParserError) Error() (string) {
    return fmt.Sprintf("Syntax error at [%v, %v]: %v", e.pos.Line, e.pos.Col, e.msg)
}

type Parser struct {
    lex *lexer.Lexer
}

func New(lex *lexer.Lexer) (*Parser) {
    return &Parser{lex}
}

func print(x... interface{}) {
    fmt.Println(fmt.Sprintf("\n%+v", x))
}

func (p *Parser) raiseError(msg interface{}) {
    if DebugFlag {
        panic(ParserError{
            pos: p.lex.Pos(),
            msg: fmt.Sprintf("%v", msg)})
    } else {
        panic(ParserError{
            pos: p.lex.Pos(),
            msg: fmt.Sprintf("%v", msg)})
    }
}

func (p *Parser) Parse() (expr interface{}){
    defer func() {
		if x := recover(); x != nil {
		    if _, ok := x.(lexer.LexerError); ok {
		        panic(x)
	        } else if _, ok := x.(ParserError); ok {
		        panic(x)
	        } else if _, ok := x.(WaitError); ok {
	            panic(x)
	        } else {
		        p.raiseError(x)
	        }
		}
	}()
	
    p.lex.Init()    
    expr = p.parseLexem(p.lex.NextLexem())
    p.lex.Flush()

    return expr
}

func (p *Parser) parseLexem(lexem lexer.Lexem) (expr interface{}) {        
    switch lexem.Value() {
        case "" :
            return nil
        case "(" :
            expr := make([]interface{}, 0)
            for {
                lexem = p.lex.NextLexem()
                if ")" == lexem.Value() {
                    return expr
                } else if inner_expr := p.parseLexem(lexem); nil != inner_expr {
                    expr = append(expr, inner_expr)
                } else {
                    RaiseWaitError()
                }             
            }
        case ")" :
            p.raiseError("unexpected )")
        default:
            if q, ok := MQuotes[lexem.Value()]; ok {                
                lexem := p.lex.NextLexem()
                if inner_expr := p.parseLexem(lexem); nil != inner_expr {
                    return []interface{} {lexer.NewLexem(q, lexem.Pos()), inner_expr}
                } else {
                    RaiseWaitError()
                }                                
            } else {
                 return lexem
            }
    }

    p.raiseError("unreachable code!")
    return nil
}
