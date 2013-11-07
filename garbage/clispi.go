package main

import (
    "fmt"
    "os"
//    "strings"
    "errors"
    "strconv"
    "math/big"
    "../clispi/decoder"
    "../clispi/lexer"
    "../clispi/parser"
    . "../clispi/globals"
)

//-------------------------------------------------------------------------------------//
func require(p bool, err_msg string) {
    if !p {
        raise(err_msg)
    }
}

func requireLen(x Object, len int, fname string) {
    _len := 0
    for !x.IsNil() {
        len ++
        x = cdr(x)
    }
    
    require(_len == len, fmt.Sprintf("%v requires %v arguments - got %v", fname, len, _len))
}

//--------------------------------Types Declaration------------------------------------//
type Object interface {
    Car() Object
    Cdr() Object
    
    IsEq(Object) bool
    IsAtom() bool
    IsNil() bool
    IsCons() bool
    IsTrueConst() bool
    IsTrue() bool
    IsConst() bool
    
    String() string
}

type AtomObject struct {
    name string
    is_const bool
    value interface{}
}

type ConsObject struct {
    car, cdr Object
}

type Function func(Object, Object) Object

//---------------------------------Atom constants--------------------------------------//
var (
    atom_T = &AtomObject{"t", true, true}
    atom_Nil = &AtomObject{"nil", true, nil}
)

//------------------Object Interface Implementation: AtomObject------------------------//
func (x *AtomObject) Cdr() Object {
    require(x.IsNil(), "cdr")
    return atom_Nil
}
func (x *AtomObject) Car() Object { 
    require(x.IsNil(), "car")
    return x 
}

func (x *AtomObject) IsEq(y Object) bool {
    if x == y { return true }
    if !y.IsAtom() { return false }   
    return x.Value() == y.(*AtomObject).Value()
}
func (x *AtomObject) IsAtom() bool { return true }
func (x *AtomObject) IsNil() bool { return x == atom_Nil }
func (x *AtomObject) IsCons() bool { return false }
func (x *AtomObject) IsConst() bool { return x.is_const }

func (a *AtomObject) String() string {
    switch x := a.value.(type) {
        case nil : return "nil"
        case bool : 
            if true == x { 
                return "t" 
            } else { 
                raise("bool value false")
                return "nil" 
            }
        case int64, float64 : return fmt.Sprint(x)
        case string : 
            if a.IsConst() { 
                return strconv.Quote(x) 
            } else {
                return x
            }
        case *func(Object, Object) Object : 
            return fmt.Sprintf("{ const function %v }", a.name)
        case Object: 
            return fmt.Sprintf("{ closure: %v }", x.String())
        default:
            if y, ok := x.(*big.Int); ok { return fmt.Sprint(y) }
            raise(fmt.Sprintf("unable to convert %v to string", t))
    }
    
    return "#UNKNOWN#"
}

//------------------AtomObject non Object Interface methods----------------------------//
func (x *AtomObject) SetValue(y Object) {
    require(!x.IsConst(), "set to constant symbol")
    x.value = y
}
func (x *AtomObject) Value() Object {
    if x.IsConst() { return x }
    return atom(x.value)
}

func atom(x interface{}) Object {
    switch y := x.(type) {
        case *AtomObject : return y
        case *func(Object, Object) Object : return &AtomObject{"", true, y }
        case string : return stringToAtom(y) 
        default :
            raise(fmt.Sprintf("unable to convert %v to atom", x))
            
    return atom_Nil
}

//------------------Object Interface Implementation: ConsObject------------------------//
func (x *ConsObject) Cdr() Object { return x.cdr }
func (x *ConsObject) Car() Object { return x.car }

func (x *ConsObject) IsEq(y Object) bool {
    if x == y { return true }
    if y.!IsCons() { return false }
    
    return car(x).IsEq(car(y)) && cdr(x).IsEq(cdr(y))
}
func (x *ConsObject) IsAtom() bool { return false }
func (x *ConsObject) IsNil() bool { return false }
func (x *ConsObject) IsCons() bool { return true }
func (x *ConsObject) IsConst() bool { return false}

func (x *ConsObject) String() string {
    if cdr(x).IsNil() {
        return fmt.Sprintf("(%v)", car(x).String())
    }

    car := car(x).String()
    cdr := cdr(x).String()
    
    if cdr(x).IsAtom && !cdr(x).IsNil() {
        return fmt.Sprintf("(%v . %v)", car, cdr)
    }
    return fmt.Sprintf("(%v %v", car, cdr[1: ])
}

//------------------ConsObject non Object Interface methods----------------------------//
func (x *ConsObject) SetCar(y Object) { x.car = y }
func (x *ConsObject) SetCdr(y Object) { x.cdr = y }

//-----------------------------Base unctions------------------------------------------//
func cons(car Object, cdr Object) Object { return &ConsObject{car, cdr} }
func car(x Object) Object { return x.Car() }
func cdr(x Object) Object { return x.Cdr() }
func caar(x Object) Object { return car(car(x)) }
func cadr(x Object) Object { return car(cdr(x)) }
func cdar(x Object) Object { return cdr(car(x)) }
func cddr(x Object) Object { return cdr(cdr(x)) }

func p_eq(args Object) bool {
    requireLen(args, 2, "eq")
    x := car(args)
    return x.IsAtom() && x.IsEq(cadr(args))
}
func p_atom(args Object) bool {
    requireLen(args, 1, "atom")
    return car(args).IsAtom()
}

func fn_cons(args Object, env Object) Object {
    requireLen(args, 2, "cons")
    return cons(car(args), cdr(args))
}

func fn_car(args Object, env Object) Object {
    requireLen(args, 1, "car")
    return car(args)
}

func fn_cdr(args Object, env Object) Object {
    requireLen(args, 1, "cdr")
    return cdr(args)
}

f

func init_env() {

    consts = map[string] *AtomObject {
        "t": atom_T, "nil": atom_Nil,
        "cons" : &AtomObject{"cons", true, fn_cons},
        "car" : &AtomObject{"car", true, fn_car},
        "cdr" : &AtomObject{"cdr", true, fn_cdr},
        "atom" : &AtomObject{"atom", true, fn_atom},
        "eq" : &AtomObject{"eq", true, fn_eq}}
}

func FromString(lexem string) (Object) {    
    if x, err := strconv.ParseInt(lexem, 0, 32); nil == err { 
        return &AtomObject{"", true, x} 
    } 

    bigInt := new(big.Int)
    if _, err := fmt.Sscan(lexem, bigInt); nil == err { 
        return &AtomObject{"", true, bigInt} }
    if x, err := strconv.ParseFloat(lexem, 64); nil == err { 
        return &AtomObject{"", true, x} }
    if '"' == lexem[0] { 
        if x, err := strconv.Unquote(lexem); nil == err { 
        return &AtomObject{"", true, x} }
    }
            
    return &AtomObject{"", false, lexem}
}
//-----------------------------------REPL----------------------------------------------//

var (
    input = os.Stdin
    output = os.Stdout
)

const (
    prompt = "clispi>"
    wait_prompt = ">"
)

func fromStdin() bool { return input == os.Stdin }

func raise(msg string) {
    panic(errors.New(msg))
}

func print(x... interface{}) {
    output.WriteString(fmt.Sprintf("%v\n", x))
    output.Flush()
}

func read(input *os.File, rune_stream chan rune) {
	
	raw_stream := make(chan byte)
	defer close(raw_stream)
	
	useEncoding := decoder.E_UNKNOWN
	if fromStdin() { useEncoding = decoder.E_UTF8 }	
	
	go decoder.Decode(raw_stream, rune_stream, useEncoding)
		
	buffer := make([]byte, 64, 64)		
	for {
		c, err := input.Read(buffer)
						
		if nil != err {
			//error
			return
		}
		
		for i, val := range buffer{
			if i < c {
				raw_stream <- val
			}
		}		
	}
}


func repl(in *os.File) {

    rune_stream := make(chan rune)
        
    go read(in, rune_stream)
    
    lex := lexer.New()
    tryEval(lex)
    
    for r := range rune_stream {
        lex.Append(r)
                
        if '\n' == r && fromStdin() {
            tryEval(lex)
        }         
    }

    tryEval(lex)
}
func tryEval(lex *lexer.Lexer) {
    stdout := os.Stdout  
    
    defer func() {
		if x := recover(); x != nil {
		    if _, ok := x.(WaitError); ok {
		        //WAIT ERROR
	            stdout.WriteString(wait_prompt)
	        } else {
			    stdout.WriteString(fmt.Sprintf("\n%v\n", x))
			    stdout.WriteString(prompt)
			
			    lex = lexer.New()
		    }
		}
	}()
	
	p := parser.New(lex)
	
	for {
	    tree := p.Parse()
	    if nil == tree { break }
        
        stdout.WriteString("\nexpr: ")
        expr := toExpr(tree)
        print(expr)
//        res := fn_eval(expr, global_env)
//        
//        print("------")
//        print(res)
//        print("------")
	}
	
    lex = lexer.New()
    stdout.WriteString(prompt)
}

func toExpr(x interface{}) (obj Object){
//x MUST NOT BE NIL!!! 
    defer func() {
//        if is_atom(obj) { print(obj.String()) }
    }()
      
    if t, ok := x.(lexer.Lexem); ok {        
        return atom(t.Value())
    } 
    arr, ok := x.([]interface{})
    if !ok {
        raise("wtf in toExpr?")
    }
    if len(arr) == 0 { 
        return atom_Nil
    }
     
    car := toExpr(arr[0])
    cdr := toExpr(arr[1: ])
    
    return cons(car, cdr)
}

func main() {
    defer input.Close()
    
 	if len(os.Args) == 2 {
     	f_input, e := file.Open(os.Args[1])     
     	if e != nil {
     		print(e)
     		os.Exit(1)
     	}
     	input = f_input
    }
 	    
    global_env = atom_Nil
    init_env()
    
    repl(input, false)
//    repl(os.Stdin, true)
    
    fmt.Println("\nexit from clispi")
}

