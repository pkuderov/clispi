package globals

import (
    "fmt"
)

var DebugFlag bool = false

func IsSpecialSymbol(r rune) (bool) {
    switch r {
        case '\'', '`', ',', '(', ')' : return true
    }    
    return false
}

const (    
    TQuote = "quote"
    THashtag = "hashtag"
    TQuasiauote = "quasiquote"
    TUnquote = "unquote"
    TUnquoteSplicing = "unquote-splicing"
)

var MQuotes = map[string] string {"'": TQuote}
    

//------------------------Position------------------------------------------------//

type Position struct {
    Line	int
    Col		int
}

func (pos *Position) Next() {
    pos.Col ++
}

func (pos *Position) NextLine() {
    pos.Col = 1
    pos.Line ++
}

//------------------------WaitError------------------------------------------------//
type WaitError struct {
    msg string
}

func RaiseWaitError() {
    panic(WaitError{""})
}

func (e WaitError) Error() (string) {
    return fmt.Sprintf("Unexpected end of stream", e.msg)
}


