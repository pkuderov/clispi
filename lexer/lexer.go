package lexer

import (
    "fmt"
    "unicode"
    "strings"
    . "../../clispi/globals"
)

//-------------------------LexerError-----------------------------------------//
type LexerError struct {
    pos Position
    msg string
}

func (e LexerError) Error() (string) {
    return fmt.Sprintf("Lexical error at [%v, %v] : %v", e.pos.Line, e.pos.Col, e.msg)
}

//---------------------------Lexem--------------------------------------------//
type Lexem struct {
    value string
    
    pos Position
}

func NewLexem(value string, pos Position) Lexem {
    return Lexem{value, pos}
}

func (l Lexem) Value() string {
    return l.value
}

func (l Lexem) Pos() Position {
    return l.pos
}

//---------------------------Lexer--------------------------------------------//
type Lexer struct {
    buffer []rune       //reading buffer
    r rune              //last read rune
    ind int             //next rune to read
    eob bool            //can read next rune
    
    currentPos Position //current position of the rune r
    nextPos Position    //next position
}

func New() (*Lexer) {
    return &Lexer{make([]rune, 0), 0, 0, true, Position{1, 0}, Position{1, 1}}
}

func (lex *Lexer) raiseError(msg interface{}) {
    if DebugFlag {
        panic(LexerError{
            pos: lex.Pos(),
            msg: fmt.Sprintf("%v\n%+v", msg, lex)})
    } else {
        panic(LexerError{
            pos: lex.Pos(),
            msg: fmt.Sprintf("%v", msg)})
    }
}

func print(x... interface{}) {
    fmt.Println(fmt.Sprintf("%+v", x))
}

func (lex *Lexer) next(){
    if lex.eob {   
        lex.raiseError("unexpected end of input stream")
    } else {        
        lex.r = lex.buffer[lex.ind]        
        lex.ind ++
        lex.eob = lex.ind >= len(lex.buffer)
        
        lex.currentPos = lex.nextPos
        if '\n' == lex.r {
            (&lex.nextPos).NextLine()
        } else {
            (&lex.nextPos).Next()
        }
    }
}

func (lex *Lexer) getNext() (next_r rune, ok bool) {
    if !lex.eob {
        return lex.buffer[lex.ind], true
    }
    return 0, false
}

func (lex *Lexer) Init() {
    lex.ind = 0
    lex.eob = lex.ind >= len(lex.buffer)
}

func (lex *Lexer) Flush() {
    if !lex.eob {
        lex.buffer = lex.buffer[lex.ind: ]
    } else {
        lex.buffer = make([]rune, 0)
    }
}

func (lex *Lexer) Append(r rune) {
    lex.buffer = append(lex.buffer, r)
}

func (lex *Lexer) Pos() (Position) { return lex.currentPos }
func (lex *Lexer) Index() (int) { return lex.ind }

func (lex *Lexer) NextLexem() (lexem Lexem) {
    defer func() {
		if x := recover(); x != nil {
		    if _, ok := x.(LexerError); ok {
		        panic(x)
	        } else if _, ok := x.(WaitError); ok {
	            panic(x)
	        } else {
		        lex.raiseError(x)
	        }
		}
	}()
		
    for !lex.eob {
        lex.next()
                                
        if unicode.IsSpace(lex.r) {
            //nothing to do here
        } else if IsSpecialSymbol(lex.r) {
            //special symbol as lexem
            lexem = Lexem{fmt.Sprintf("%c", lex.r), lex.Pos()}
//            if r, ok := lex.getNext(); ok {
//                //check if hashtag
//                temp_pair := fmt.Sprintf("%c%c", lex.r, r)
//                _lexem := Lexem{temp_pair, lex.Pos()}
//                if q, ok := MQuotes[temp_pair]; ok && THashtag == q {
//                    lex.next()
//                    lexem = _lexem
//                }
//            }
            return lexem
        } else if '"' == lex.r { 
            //string
            start_str := lex.ind - 1
            pos := lex.Pos()                    
            for !lex.eob {           
                lex.next()                                
                if '\\' == lex.r {
                    lex.next()
                } else if '"' == lex.r {
                    return Lexem{string(lex.buffer[start_str:lex.ind]), pos}
                }
            }     
            RaiseWaitError()
        } else if ';' == lex.r {
            //comments
            for !lex.eob && '\n' != lex.r {
                lex.next()
            }            
        } else {//symbol
            str_start := lex.ind - 1
            pos := lex.Pos()
            for !lex.eob {                
                if r, _ := lex.getNext(); IsSpecialSymbol(r) || unicode.IsSpace(r) {
                    break
                }                                
                lex.next()                
            }                        
            return Lexem{strings.ToLower(string(lex.buffer[str_start:lex.ind])), pos}
        }        
    }    
    return Lexem{"", lex.Pos()}
}
