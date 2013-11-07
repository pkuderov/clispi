package main

import (
    "fmt"
    "os"
    clispi "../clispi/interpreter"
)

var (    
    input = os.Stdin
    output = os.Stdout
    init_file = "/home/petr/go/src/pkuderov/clispi/init.lisp"
)

func init_repl() {
 	fmt.Println("===========init interpreter=============")
 	
 	clispi.InitInterpreter()
 	
 	f, e := os.Open("init.lisp")
 	defer f.Close()
 	
 	if e != nil {
 		fmt.Print(e)
 		return
 	} 	
 	clispi.Repl(f, output)
 	
 	fmt.Println("============you're welcome==============")
}

func main() {       
    defer func() {
        input.Close()
        output.Close()
    }()    
     
    init_repl()
     
 	if len(os.Args) > 1 {
     	f_input, e := os.Open(os.Args[1])
     	if e != nil {
     		fmt.Print(e)
     		os.Exit(1)
     	}
     	input = f_input
    }
 	if len(os.Args) > 2 {
     	f_output, e := os.Open(os.Args[2])
     	if e != nil {
     		fmt.Print(e)
     		os.Exit(1)
     	}
     	output = f_output
    }   
 	     
//    clispi.Repl(input, output)
    
    fmt.Println("\nexit from clispi\n")
}

////-----------------------------------REPL----------------------------------------------//
//var (
//    input = os.Stdin
//    output = os.Stdout
//)

//const (
//    prompt = "clispi>"
//    wait_prompt = ">"
//)

//func fromStdin() bool { 
//    return input == os.Stdin 
//}

//func raise(msg string) {
//    panic(errors.New(msg))
//}

//func read(input *os.File, rune_stream chan rune) {	
//	raw_stream := make(chan byte)
//	defer close(raw_stream)
//	
//	useEncoding := decoder.E_UNKNOWN
//	if fromStdin() { 
//	    useEncoding = decoder.E_UTF8 
//    }	
//	go decoder.Decode(raw_stream, rune_stream, useEncoding)
//		
//	buffer := make([]byte, 64, 64)		
//	for {
//		c, err := input.Read(buffer)
//						
//		if nil != err {
//		    panic(err)
//		}		
//		for i, val := range buffer{
//			if i < c {
//				raw_stream <- val
//			}
//		}		
//	}
//}

//func repl(in *os.File) {
//    rune_stream := make(chan rune)
//        
//    go read(in, rune_stream)
//    
//    lex := lexer.New()
//    tryEval(lex)
//    
//    for r := range rune_stream {
//        lex.Append(r)
//                
//        if '\n' == r && fromStdin() {
//            tryEval(lex)
//        }         
//    }

//    tryEval(lex)
//}
//func tryEval(lex *lexer.Lexer) {    
//    defer func() {
//		if x := recover(); x != nil {
//		    if _, ok := x.(WaitError); ok {
//		        //WAIT ERROR
//	            output.WriteString(wait_prompt)
//	        } else {
//			    output.WriteString(fmt.Sprintf("\n%v\n", x))
//			    output.WriteString(prompt)
//			
//			    lex = lexer.New()
//		    }
//		}
//	}()
//	
//	p := parser.New(lex)
//	
//	for {
//	    tree := p.Parse()
//	    if nil == tree { break }
//        
//        stdout.WriteString("expr: ")
//        expr := toExpr(tree)
//        print(expr)
//        
//        res := fn_eval(expr, global_env)
//        
//        print("------")
//        print(res)
//        print("------")
//        stdout.WriteString("\n")
//	}
//	
//    lex = lexer.New()
//    stdout.WriteString(prompt)
//}

////-------------------------------------------------------------------------------------//
//func require(p bool, err_msg string) {
//    if !p {
//        raise(err_msg)
//    }
//}

//func requireLen(x interface{}, len int, fname string) {
//    _len := 0
//    for is_cons(x) {
//        len ++
//        x = cdr(x)
//    }
//    
//    if !is_nil(x) {
//        len++
//    }
//    
//    require(_len == len, fmt.Sprintf("%v requires %v arguments - got %v", fname, len, _len))
//}

//type AtomObject struct {
//    value interface{}
//}

//type ConsObject struct {
//    car, cdr interface{}
//}

//func new_atom(name string) *AtomObject {
//    x := AtomObject{name}
//    return &x
//}

//const (
//    t_nil = iota
//    t_bool
//    t_int
//    t_bigint
//    t_float
//    t_string
//    t_atom
//    t_cons
//    t_func
//)

//var (
//    consts = map[string] interface{} {
//        "t": true, "nil": nil}
//    global_env = cons(nil, nil)
//)

//func typeof(x interface{}) int {
//    switch i := x.(type) {
//        case nil : return t_nil
//        case bool: return t_bool
//        case int64: return t_int
//        case float64: return t_float
//        case *big.Int: return t_bigint
//        case string: return t_string
//        case *AtomObject: return t_atom
//        case *ConsObject: return t_cons
//        case func(interface{}, interface{}) interface{} : return t_func
//        default:
//            raise(fmt.Sprintf("%v has unknown type", i))
//    }
//    return -1
//}

//func valueof(x interface{}) interface{} {
//    switch i := x.(type) { 
//        case *AtomObject : return i.value
//        case *ConsObject : return i.car
//    }
//    return x
//}

//func is_nil(x interface{}) bool { return nil == x }
//func is_cons(x interface{}) bool {
//    switch typeof(x) {
//        case t_nil, t_cons : return true
//    }
//    return false
//}
//func is_atom(x interface{}) bool {
//    switch typeof(x) {
//        case t_nil, t_bool, t_int, t_float, t_bigint, t_string, t_func : return true
//        case t_atom : return true
//    }
//    return false
//}
//func is_equal(x interface{}, y interface{}) bool {
//    return valueof(x) == valueof(y)
//}

//func car(x interface{}) interface{} {
//    switch i := x.(type) {
//        case nil : return nil
//        case *ConsObject: return i.car
//        default:
//            raise(fmt.Sprintf("car from noncons object %v", i))
//    }
//    return nil
//}
//func cdr(x interface{}) interface{} {
//    switch i := x.(type) {
//        case nil : return nil
//        case *ConsObject: return i.cdr
//        default:
//            raise(fmt.Sprintf("cdr from noncons object %v", i))
//    }
//    return nil
//}
//func cons(car interface{}, cdr interface{}) interface{} {
//    x := ConsObject{car, cdr}
//    return &x
//}

//func caar(x interface{}) interface{} { return car(car(x)) }
//func cadr(x interface{}) interface{} { return car(cdr(x)) }
//func cdar(x interface{}) interface{} { return cdr(car(x)) }
//func cddr(x interface{}) interface{} { return cdr(cdr(x)) }
////----------------------------------BASE LISP------------------------------------------//
//func p_eq(args interface{}, env interface{}) interface{} {
////    defer func() {
////        if err := recover(); err != nil {
////            requireLen(args, 2, "eq")
////        }
////    }()
//    x := fn_eval(car(args), env)
//    if is_atom(x) {
//        y := fn_eval(cadr(args), env)
//        if is_atom(y) && is_equal(x, y) { return true }
//    }
//    return nil
//}
//func p_atom(args interface{}, env interface{}) interface{} {
////    defer func() {
////        if err := recover(); err != nil {
////            requireLen(args, 1, "atom")
////            panic(err)
////        }
////    }()
//    if is_atom(fn_eval(car(args), env)) { return true }
//    return nil
//}
//func fn_cons(args interface{}, env interface{}) interface{} {
////    defer func() {
////        if err := recover(); err != nil {
////            requireLen(args, 2, "cons")
////        }
////    }()
//    return cons(fn_eval(car(args), env), fn_eval(cadr(args), env))
//}

//func fn_car(args interface{}, env interface{}) interface{} {
////    defer func() {
////        if err := recover(); err != nil {
////            requireLen(args, 1, "car")
////        }
////    }()
//    return fn_eval(caar(args), env)
//}

//func fn_cdr(args interface{}, env interface{}) interface{} {
////    defer func() {
////        if err := recover(); err != nil {
////            requireLen(args, 1, "cdr")
////        }
////    }()
//    return fn_eval(cdar(args), env)
//}

//func fn_cond(args interface{}, env interface{}) interface{} {
//    for nil != args {
//        p := fn_eval(caar(args), env)
//        if nil != p { 
//            return fn_eval(cadr(car(args)), env)
//        }
//        args = cdr(args)
//    }
//    
//    return nil
//}

//func fn_setf(args interface{}, env interface{}) interface{} {
////    defer func() {
////        if err := recover(); err != nil {
////            requireLen(args, 2, "setf")
////            panic(err)
////        }
////    }()
//    
//    x := car(args)
//    if atom, ok := x.(*AtomObject); ok {
//        name := atom.value.(string)
//        if _, ok := consts[name]; ok {
//            raise("can't set value to constant")
//        }        
//        value := fn_eval(cadr(args), env)        
//        
//        if pair := find(name, env); pair != nil {
//            m := cdr(pair).(map[string] interface{})
//            m[name] = value
//        } else {
//            bind(name, value, env)
//        }        
//        return value
//    } else {
//        raise(fmt.Sprintf("%v is not of type symbol", x))
//    }
//        
//    return nil
//}

//func fn_quote(args interface{}, env interface{}) interface{} {
//    requireLen(args, 1, "quote")
//    print("__")
//    print(args)
//    return car(args)
//}

//func fn_lambda(args interface{}, env interface{}) interface{} {
//    vars := car(args)
//    body := cadr(args)
//    
//    return cons(cons(vars, body), cons(nil, env))
//}

//func fn_eval(args interface{}, env interface{}) interface{} {
////    defer func() {
////        if err := recover(); err != nil {
////            requireLen(args, 1, "eval")
////            panic(err)
////        }
////    }()
//    
//    switch typeof(args) {        
//        case t_nil, t_bool, t_int, t_float, t_bigint, t_string, t_func : 
//            return args
//        case t_atom :
//            name := valueof(args).(string)
//            if x := find(name, env); x != nil {
//                return car(x)
//            } else {
//                raise(fmt.Sprintf("%v is unbound", name))
//            }
//        case t_cons :
//            t := fn_eval(car(args), env)
//            switch fn := t.(type) {
//                case *ConsObject:
//                    //closure
//                    print("==welcome to closure==")
//                    
//                    f_args := cdr(args)
//                    vars := caar(fn)
//                    body := cdar(fn)
//                    local_env := cdr(fn)
//                    
//                    for vars != nil && f_args != nil {
//                        var_name := car(vars).(*AtomObject).value.(string)
//                        var_value := fn_eval(car(f_args), env)
//                        
//                        bind(var_name, var_value, local_env)
//                        
//                        vars = cdr(vars)
//                        f_args = cdr(f_args)
//                    }
//                    if vars != nil || f_args != nil {
//                        raise("number of variables and received arguments are different")
//                    }
//                    return fn_eval(body, local_env)
//                case func(interface{}, interface{}) interface{} :
//                    //const function
//                    print("==welcome to const function==")
//                    print(fn)
//                    return fn(cdr(args), env)
//            }
//    }
//    
//    raise(fmt.Sprintf("can't evaluate %v", args))
//    return true
//}


//func find(x string, env interface{}) interface{} {
//    if y, ok := consts[x]; ok {
//        return cons(y, nil)
//    }
//    
//    for is_cons(env) && nil != env {
//        t := car(env)
//        if nil != t {
//            m := t.(map[string] interface{})
//            if y, ok := m[x]; ok {
//                return cons(y, m)
//            }
//        }        
//        env = cdr(env)        
//    }
//    
//    return nil
//}

//func bind(x string, value interface{}, env interface{}) {
//    t := env.(*ConsObject)
//    
//    if nil == t.car {
//        t.car = map[string] interface{} {}
//    }
//    
//    if m, ok := t.car.(map[string] interface{}); ok {
//        m[x] = value
//    } else {
//        raise("environment has damaged")
//    }
//    
//}

//func parseAtom(lexem string) interface{} {      
//    if "nil" == lexem { 
//        return nil 
//    }
//    if "t" == lexem { 
//        return true 
//    }
//    
//    if x, err := strconv.ParseInt(lexem, 0, 32); nil == err { 
//        return x
//    } 
//    bigInt := new(big.Int)
//    if _, err := fmt.Sscan(lexem, bigInt); nil == err { 
//        return bigInt
//    }
//    if x, err := strconv.ParseFloat(lexem, 64); nil == err { 
//        return x
//    }
//    if '"' == lexem[0] { 
//        if s, err := strconv.Unquote(lexem); err != nil {
//            return s
//        }
//        raise("cannot unquote string")
//    }
//    
//    x := new_atom(lexem)    
//    return x
//}
//func to_string(x interface{}) string {    
//    switch y := x.(type) {
//        case nil :
//            return "nil"
//        case bool : 
//            if true == y { 
//                return "t" 
//            } else { 
//                raise("bool value false")
//                return "nil" 
//            }
//        case string : 
//            return strconv.Quote(y)
//        case func(interface{}, interface{}) interface{} : 
//            return "{ const function }"
//        case *AtomObject : 
//            return y.value.(string)
//        case *ConsObject :            
//            if is_nil(cdr(y)) {
//                return fmt.Sprintf("(%v)", to_string(car(y)))
//            }

//            s_car := to_string(car(y))
//            s_cdr := to_string(cdr(y))
//        
//            if is_atom(cdr(x)) && !is_nil(cdr(x)) {
//                return fmt.Sprintf("(%v . %v)", s_car, s_cdr)
//            }
//            return fmt.Sprintf("(%v %v", s_car, s_cdr[1: ])
//        case int64, float64, *big.Int : 
//            return fmt.Sprint(y)
//        case map[string] interface{} :
//            return fmt.Sprint(y)
//    }
//    
//    return "#UNKNOWN#"
//}

//func toExpr(x interface{}) interface{} {      
//    if t, ok := x.(lexer.Lexem); ok {   
//        return parseAtom(t.Value())
//    } 
//    arr, ok := x.([]interface{})
//    if len(arr) == 0 { 
//        return nil
//    }
//     
//    car := toExpr(arr[0])
//    cdr := toExpr(arr[1: ])
//    
//    return cons(car, cdr)
//}
