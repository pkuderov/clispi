package interpreter

import (
    "fmt"
    "os"
    "io"
    "errors"
    "strconv"
    "math/big"
    "../../clispi/decoder"
    "../../clispi/lexer"
    "../../clispi/parser"
    . "../../clispi/globals"
)

//-----------------------------------------------------------------------------------------//
func raise(msg string) {
    panic(errors.New(msg))
}

func debug(x interface{}) {
    if DebugFlag {
        print(to_string(x), "\n")
    }
}

func require(p bool, err_msg string) {
    if !p {
        raise(err_msg)
    }
}
func requireList(x interface{}, fname string) {
    require(is_cons(x), 
        fmt.Sprintf("error in %v: %v is not type of List", fname, to_string(x)))
}
func requireLen(x interface{}, len int, fname string) {
    _len := 0
    for nil != x {
        requireList(x, fname)
        _len ++
        x = cdr(x)
    }
    
    require(_len == len, fmt.Sprintf("%v requires %v arguments - received %v", fname, len, _len))
}

//-----------------------------------------------------------------------------------------//

type Symbol struct {
    value string
}
type VListBlock struct {
    p *VList
    
    size int
    last_used int
    data []interface{}
}
type VList struct {
    base *VListBlock
    offset int
}
type LList struct {
    car, cdr interface{}
}
type Cons interface {
    Car() interface{}
    Cdr() interface{}
    
    SetCar(interface{})
    
    Nth(int) interface{}
    Len() int
}
type Closure struct {
    vars interface{}
    body interface{}
    env *Context
}
type Context struct {
    vars map [string] interface{}    
    outer Cons
}

const (
    t_nil = iota
    t_bool
    t_int
    t_bigint
    t_float
    t_string
    t_symbol
    t_cons
    t_func
    t_closure
    t_context
    t_map
)

var (
    consts = map[string] interface{} {}
    global_env *Context = &Context{nil, nil}
)

func typeof(x interface{}) int {
    switch i := x.(type) {
        case nil : return t_nil
        case bool: return t_bool
        case int64: return t_int
        case float64: return t_float
        case *big.Int: return t_bigint
        case string: return t_string
        case *Symbol: return t_symbol
        case Cons: return t_cons
        case func(interface{}, *Context) interface{}: return t_func
        case *Closure : return t_closure
        case *Context: return t_context
        case map[string] interface{}: return t_map
        default:
            raise(fmt.Sprintf("%v has unknown type", i))
    }
    return -1
}

//---------------------------LList methods-------------------------------------------------//
func (x *LList) Car() interface{} {
    return x.car
}
func (x *LList) Cdr() interface{} {
    return x.cdr
}
func (x *LList) SetCar(y interface{}) {
    x.car = y
}
func (x *LList) Nth(n int) interface{} {
    var y interface{} = x
    
    for n > 1 {
        y = cdr(y)
        n--
    }
    return y
}
func (x *LList) Len() (res int) {
    res = 0
    var y interface{} = x
    for y != nil{
        y = cdr(y)
        res++
    }
    return res
}
func llist_cons(car interface{}, cdr interface{}) interface{} {
    x := LList{car, cdr}
    return &x
}
//---------------------------VList methods-------------------------------------------------//
func (x *VList) Car() interface{} {
    if x == nil || x.base == nil {
        return nil
    }
    return x.base.data[x.offset]
}
func (x *VList) Clone() *VList {
    if x == nil || x.base == nil { 
        return nil
    }    
    
    y := VList{x.base, x.offset}
    return &y
}
func (x *VList) Cdr() interface{} {
    if x == nil || x.base == nil { 
        return nil
    }    
    
    y := x.Clone();
    return dam_cdr(y, 1)
}

func dam_cdr(x *VList, n int) interface{} {
    debug("_----dam cdr")
    if x == nil || x.base == nil {
        return nil
    }
    
    if x.offset - n < 0 {
        n -= x.offset + 1
        x = x.base.p
        return dam_cdr(x, n)
    } 

    x.offset -= n
    return x
}
func (x *VList) SetCar(y interface{}) {
    debug("_----caaar")    
    if x == nil || x.base == nil { 
        return
    }    
    x.base.data[x.offset] = y
}
func (x *VList) Nth(n int) interface{} {
    if x == nil || x.base == nil { 
        return nil
    }    
    y := x.Clone()  
    return dam_cdr(y, n)
}
func (x *VList) Len() (res int) {
    debug("_----car")
    if x == nil || x.base == nil { 
        return
    }
    
    res = 0
    y := x.Clone()
    for y != nil && y.base != nil {
        res += y.offset + 1
        y = y.base.p
    }
    return res
}
func vlist_cons(car interface{}, cdr interface{}) interface{} {
    var x *VList = nil
    var size = 1

    debug("_----cons")
    debug(car)
    debug(cdr)
    if cdr != nil {
        if _, ok := cdr.(*VList); !ok {
            cdr = vlist_cons(cdr, nil)
            x = cdr.(*VList)
        } else {
            x = cdr.(*VList)
        }
                
        y := x.base
        if x.offset == y.last_used && y.last_used < y.size - 1 {
            y.last_used ++
            y.data[y.last_used] = car
            res := VList{x.base, y.last_used}
            return &res
        }
        
        if x.offset == y.last_used {
            size = y.size * 2
        }
    }    
        debug("_----===")
    
    block := VListBlock{x, size, 0, make([] interface{}, size)}
    block.data[0] = car
    res := VList{&block, 0}
    return &res
}
//-----------------------Service functions: USE WITH WARN!!!-------------------------------//
func valueof(x interface{}) interface{} {
    if a, ok := x.(*Symbol); ok {
        return a.value
    }
    return x
}
func is_cons(x interface{}) bool {
    switch typeof(x) {
        case t_nil, t_cons : return true
    }
    return false
}
func is_atom(x interface{}) bool {
    switch typeof(x) {
        case t_nil, t_bool, t_int, t_float, t_bigint, t_string, t_func, t_map : return true
        case t_symbol : return true
    }
    return false
}
func is_equal(x interface{}, y interface{}) bool {
    return valueof(x) == valueof(y)
}
func to_symbol(x interface{}) *Symbol {
    if a, ok := x.(*Symbol); ok {
        return a
    }
    raise(fmt.Sprintf("%v is not of type Symbol", to_string(x)))
    return nil
}
func to_cons(x interface{}) Cons {
    if c, ok := x.(Cons); ok {
        return c
    }
    raise(fmt.Sprintf("%v is not of type Cons", to_string(x)))
    return nil
}
func to_func(x interface{}) func(interface{}, *Context) interface{} {    
    if f, ok := x.(func(interface{}, *Context) interface{}); ok {
        return f
    }
    raise(fmt.Sprintf("%v is not of type { const function }", to_string(x)))
    return nil
}
func to_closure(x interface{}) *Closure {
    if f, ok := x.(*Closure); ok {
        return f
    }
    raise(fmt.Sprintf("%v is not of type { closure }", to_string(x)))
    return nil    
}
func to_context(x interface{}) *Context {
    if f, ok := x.(*Context); ok {
        return f
    }
    raise(fmt.Sprintf("%v is not of type { context }", to_string(x)))
    return nil    
}
func to_map(x interface{}) map[string] interface{} {    
    if c, ok := x.(map[string] interface{}); ok {
        return c
    }
    raise(fmt.Sprintf("%v is not of type Cons", to_string(x)))
    return nil
}
//----------------------------------Base functions-----------------------------------------//
var cons = vlist_cons

func car(x interface{}) interface{} {
    switch i := x.(type) {
        case nil : return nil
        case Cons: return i.Car()
        default:
            requireList(x, "car")
    }
    return nil
}
func cdr(x interface{}) interface{} {
    switch i := x.(type) {
        case nil : return nil
        case Cons: return i.Cdr()
        default:
            requireList(x, "cdr")
    }
    return nil
}

func nth(x interface{}, n int64) interface{} {
    return to_cons(x).Nth(int(n))
}
func llen(x interface{}) (res int64) {
    return int64(to_cons(x).Len())
}
func first(x interface{}) interface{} {
    return cons(car(x), nil)
}
func caar(x interface{}) interface{} { return car(car(x)) }
func cadr(x interface{}) interface{} { return car(cdr(x)) }
func cdar(x interface{}) interface{} { return cdr(car(x)) }
func cddr(x interface{}) interface{} { return cdr(cdr(x)) }

//------------------------------------BASE LISP--------------------------------------------//
func fn_eq(args interface{}, env *Context) interface{} {
    if x := eval(first(args), env); is_atom(x) {
        if y := eval(cdr(args), env); is_atom(y) && is_equal(x, y) { 
            return true 
        }
    }
    return nil
}
func fn_atom(args interface{}, env *Context) interface{} {
    if x := eval(args, env); is_atom(x) { 
        return true 
    }
    return nil
}
func fn_cons(args interface{}, env *Context) interface{} {    
    x := eval(first(args), env)
    y := eval(cdr(args), env)
    return cons(x, y)
}
func fn_car(args interface{}, env *Context) interface{} {    
    return car(eval(args, env))
}
func fn_cdr(args interface{}, env *Context) interface{} {    
    return cdr(eval(args, env))
}
func fn_nth(args interface{}, env *Context) interface{} {
    x := eval(first(args), env).(int64)
    l := eval(cdr(args), env)
    return nth(l, x)
}
func fn_len(args interface{}, env *Context) interface{} {
    return llen(eval(args, env))    
}
func fn_cond(args interface{}, env *Context) interface{} {
    x := args
    for nil != x {
        pair := car(x)        
        if nil != eval(first(pair), env) { 
            return eval(cdr(pair), env)
        }
        x = cdr(x)
    }    
    return nil
}
func fn_defun(args interface{}, env *Context) interface{} {   
    _lambda := cons(parse_atom("lambda"), cdr(args))
//    _name := cons(parse_atom("quote"), first(args))
    _name := car(args)
    _res := cons(_name, cons(_lambda, nil))
    return fn_setq(_res, env)
}
func fn_quote(args interface{}, env *Context) interface{} {
    debug("quote===")
    return car(args)
}
func fn_lambda(args interface{}, env *Context) interface{} {
    vars := car(args)
    body := cdr(args)
    
    closure := Closure{vars, body, env}
    return &closure
}
func fn_eval(args interface{}, env *Context) interface{} {
    res := eval(args, env)
    return eval(cons(res, nil), env)
}
func eval(args interface{}, env *Context) interface{} {
    x := car(args)
    debug("--------------------------------")
    debug(x)
    switch typeof(x) {        
        case t_nil, t_bool, t_int, t_float, t_bigint, t_string, t_func, t_closure : 
            return x
        case t_symbol :
            t := to_symbol(x)
            name := t.value

            //check consts
            if pair := find_consts(name); pair != nil {
                return car(pair)
            }            
            if res, ok := eval_arg(t, env); ok {
                debug("eval found")
                debug(res)
                return res
            } else {
                debug("eval not founf")
                raise(fmt.Sprintf("%v is unbound", t.value))
            }
        case t_cons :
            //const function or closure
            switch fn := eval(first(x), env).(type) {
                case *Closure:
                    debug("--closure--")
                    return eval_closure(fn, cdr(x), env)
                case func(interface{}, *Context) interface{} :
                    //const function
                    t := fn(cdr(x), env)
                    return t
                default:
                    debug("ERROR")
                    debug(fn)
            }
    }    
    debug(x)
    if _, ok := x.(*Cons); ok {
        debug("cons")
        
    }
    //debug(env)
    debug("--------------------------------")
    raise(fmt.Sprintf("can't evaluate %v", args))
    return nil
}
func eval_closure(fn *Closure, f_args interface{}, env *Context) interface{} {
    vars := fn.vars
    body := fn.body
    local_env := &Context{nil, nil}
    
    //bind vars
    for vars != nil && f_args != nil {
        var_name := to_symbol(car(vars)).value
        var var_value interface{} = nil
        
        if '*' == var_name[0] {
            //variable num of params
            var_name = var_name[1: ]
            if '$' == var_name[0] {
                var_name = var_name[1: ]
                var_value = eval_listarg(f_args, env, true)
            } else {
                var_value = eval_listarg(f_args, env, false)
            }
            f_args = nil
        } else if '$' == var_name[0] {
            var_name = var_name[1: ]
            var_value = fn_lambda(cons(nil, first(f_args)), env)
        } else if '&' == var_name[0] {
            var_name = var_name[1: ]
            var_value = car(f_args)
        } else {
//            debug(env)
            var_value = eval(first(f_args), env)
        }
        
        debug("*)**)")
        debug(var_name)
        debug(env)
        debug(var_value)
        debug("*)**)")

        bind(var_name, var_value, local_env)
        vars = cdr(vars)
        f_args = cdr(f_args)
    }
    ok := vars == nil && f_args == nil
    require(ok, "number of variables and received arguments are different")
    
    local_env.outer = to_cons(cons(fn.env, cons(env, nil)))
        
//    debug(local_env)
    debug("closure body eval")
    res := eval(body, local_env)
    debug("closure result")
    debug(body)
    debug(res)
    return res
}
func eval_listarg(f_args interface{}, env *Context, block bool) interface{} {
    if f_args != nil {
        if block {
            x := fn_lambda(cons(nil, first(f_args)), env)
            y := eval_listarg(cdr(f_args), env, block)
            return cons(x, y)
        }
        return cons(eval(first(f_args), env), eval_listarg(cdr(f_args), env, block))
    }
    return nil
}
func eval_arg(x *Symbol, env *Context) (interface{}, bool) {
    debug("eval_arg")
    debug(x)
    debug(env.vars)
    
    if val, ok := env.vars[x.value]; ok {
        debug("okok")
        return val, true
    } else {
        outers := env.outer
        for outers != nil {
            t_env := to_context(car(outers))

            if res, ok := eval_arg(x, t_env); ok {
                return res, true
            }        
            if cdr(outers) != nil {    
                outers = to_cons(cdr(outers))
            } else {
                break
            }
        }
    }
    debug("nono")
    debug(env.vars[x.value])
    return nil, false
}
func fn_list(args interface{}, env *Context) interface{} {
    if nil == args {
        return nil
    }    
    return cons(eval(first(args), env), fn_list(cdr(args), env))
}
func fn_set(args interface{}, env *Context) interface{} {
//    go requireLen(args, 2, "set")
    
    x := to_symbol(eval(first(args), env))
    name := x.value
    require(nil == find_consts(name), "can't set value to constant")
    
    value := eval(cdr(args), env)
    bind(name, value, env)
    
    return value
}
func fn_setq(args interface{}, env *Context) interface{} {
    x := to_symbol(car(args))
    name := x.value
    require(nil == find_consts(name), "can't set value to constant")
    
    value := eval(cdr(args), env)
    
    if _, ok := env.vars[name]; ok {
        env.vars[name] = value
    } else {
        bind(name, value, env)
    }            
    return value
}
//-------------------Functions for work with enterpreter environment-----------------------//
func find_consts(x string) interface{} {
    if y, ok := consts[x]; ok {
        return cons(y, consts)
    }
    return nil
}
func bind(x string, value interface{}, env *Context) {
//    debug(value)   
    for nil != env.outer {
        env = to_context(car(env.outer))
    }         
    if nil == env.vars {
        env.vars = map[string] interface{} {}
    }    
    env.vars[x] = value
}
//-------------------------Functions for convert atoms to and from strings----------------//
func parse_atom(lexem string) interface{} {      
    if "nil" == lexem { 
        return nil 
    }
    if "t" == lexem { 
        return true 
    }
    if x, err := strconv.ParseInt(lexem, 0, 32); nil == err { 
        return x
    } 
    bigInt := new(big.Int)
    if _, err := fmt.Sscan(lexem, bigInt); nil == err { 
        return bigInt
    }
    if x, err := strconv.ParseFloat(lexem, 64); nil == err { 
        return x
    }
    if '"' == lexem[0] { 
        return lexem[1: len(lexem) - 1]
    }
    
    x := Symbol{lexem}
    return &x
}
func to_string(x interface{}) (res string) {
    switch y := x.(type) {
        case nil :
            return "nil"
        case bool : 
            if true == y { 
                return "t" 
            } else { 
                raise("bool value can't be false")
                return "nil" 
            }
        case string : 
            return y
        case func(interface{}, *Context) interface{} : 
            return "{ const function }"
        case *Symbol : 
            return y.value
        case *Cons, *LList, *VList :
            head := car(y)
            tail := cdr(y)
            
            if nil == tail {
                return fmt.Sprintf("(%s)", to_string(head))
            }
            s_car := to_string(head)
            s_cdr := to_string(tail)
        
            if is_atom(tail) {
                return fmt.Sprintf("(%s . %s)", to_string(s_car), to_string(s_cdr))
            }
            return fmt.Sprintf("(%s %s", to_string(s_car), to_string(s_cdr[1: ]))
        case *Closure:
            s := "closure = "            
            if y.vars != nil { 
                s += fmt.Sprintf(" {vars: %s} ", to_string(y.vars))
            }
            if y.body != nil {
                s += fmt.Sprintf(" {body: %s} ", to_string(y.body))
            } 
            if y.env != nil {
                s += fmt.Sprintf(" {context: %s} ", to_string(y.env))
            } 
            return s
        case *Context:
            s := "context = "
            if y.vars != nil { 
                s += fmt.Sprintf(" {local: %s} ", to_string(y.vars))
            }
            if y.outer != nil {
                s += fmt.Sprintf(" {outer: %s} ", to_string(y.outer))
            }
            return s
        case int64, float64, *big.Int : 
            return fmt.Sprint(y)
        case map[string] interface{} :
            return fmt.Sprint(y)
        case *VListBlock:
            return to_string(y.p)
    }
    
    debug(x)
    return "#UNKNOWN#"
}
//-------------------------------Math functions--------------------------------------------//
func fn_clone(args interface{}, env *Context) interface{} {
    x := fn_eval(first(args), env)
    
    switch y := x.(type) {
        case int64:
            return big.NewInt(y)
        case *big.Int:
            res := big.NewInt(0)
            return res.Add(res, y)
    }
    raise("require numeric in clone")
    return nil
}
const (
    op_add = iota
    op_sub
    op_mul
    op_div
    op_mod
    op_pow
)
func fn_add(args interface{}, env *Context) interface{} {
    x := fn_eval(first(args), env)
    y := fn_eval(cdr(args), env)
    return eval_operator(x, y, op_add)
}
func fn_sub(args interface{}, env *Context) interface{} {
    x := fn_eval(first(args), env)
    y := fn_eval(cdr(args), env)    
    return eval_operator(x, y, op_sub)
}
func fn_mul(args interface{}, env *Context) interface{} {
    x := fn_eval(first(args), env)
    y := fn_eval(cdr(args), env)
    return eval_operator(x, y, op_mul)    
}
func fn_div(args interface{}, env *Context) interface{} {
    x := fn_eval(first(args), env)
    y := fn_eval(cdr(args), env)
    return eval_operator(x, y, op_div) 
}
func fn_mod(args interface{}, env *Context) interface{} {
    x := fn_eval(first(args), env)
    y := fn_eval(cdr(args), env)
    return eval_operator(x, y, op_mod) 
}
func fn_pow(args interface{}, env *Context) interface{} {
    x := fn_eval(first(args), env)
    y := fn_eval(cdr(args), env)
    return eval_operator(x, y, op_pow) 
}
func eval_operator(_x interface{}, _y interface{}, op int) interface{} {
    switch x := _x.(type) {
        case int64:
            switch y := _y.(type) {
                case int64:                    
                    return eval_ints(x, y, op)
                case *big.Int:
                    tx := big.NewInt(x)
                    return eval_bigints(tx, y, op)
            }
        case *big.Int:
            switch y := _y.(type) {
                case int64:
                    ty := big.NewInt(y)
                    return eval_bigints(x, ty, op)
                case *big.Int:
                    return eval_bigints(x, y, op)
            }
    }
    
    raise("can't evaluate add operation")
    return nil    
}
func eval_ints(x int64, y int64, op int) interface{} {
    switch op {
        case op_add:
            x = x + y
        case op_sub:
            x = x - y
        case op_mul:
            x = x * y
        case op_div:
            x = x / y
        case op_mod:
            x = x % y
        case op_pow:
            return eval_operator(big.NewInt(x), big.NewInt(y), op_pow)
    }
    if x >= (1 << 32) || x <= -(1 << 32){
        return big.NewInt(x)
    }
    return x
}

func eval_bigints(x *big.Int, y *big.Int, op int) interface{} {    
    switch op {
        case op_add:
            x.Add(x, y)
        case op_sub:
            x.Sub(x, y)
        case op_mul:
            x.Mul(x, y)
        case op_div:
            x.Div(x, y)
        case op_mod:
            x.Mod(x, y)
        case op_pow:
            x.Exp(x, y, big.NewInt(0))
    }
    return x
}
//-------------------------------Public functions------------------------------------------//
func InitInterpreter() {    
    consts["t"] = true
    consts["nil"] = nil
 	consts["eq"] = fn_eq
 	consts["atom"] = fn_atom
 	consts["cond"] = fn_cond
 	consts["cons"] = fn_cons
 	consts["car"] = fn_car
 	consts["cdr"] = fn_cdr
 	consts["set"] = fn_set
 	consts["quote"] = fn_quote
 	consts["lambda"] = fn_lambda
 	consts["eval"] = fn_eval
 	consts["defun"] = fn_defun
 	consts["setq"] = fn_setq
 	consts["list"] = fn_list
 	
 	consts["len"] = fn_len
 	consts["nth"] = fn_nth
 	
 	consts["add"] = fn_add
 	consts["sub"] = fn_sub
 	consts["mul"] = fn_mul
 	consts["div"] = fn_div
 	consts["mod"] = fn_mod
 	consts["pow"] = fn_pow
 	consts["long"] = fn_clone
}
func ToExpression(x interface{}) interface{} {      
    if t, ok := x.(lexer.Lexem); ok {   
        return parse_atom(t.Value())
    } 
    arr, _ := x.([]interface{})
    if len(arr) == 0 { 
        return nil
    }
     
    car := ToExpression(arr[0])
    cdr := ToExpression(arr[1: ])
    
    return cons(car, cdr)
}
func Repl(in *os.File, out *os.File) {
    input = in
    output = out
    
    rune_stream := make(chan rune)
        
    go read(rune_stream)
    
    lex := lexer.New()    
    print(prompt)
    
    for r := range rune_stream {
        lex.Append(r)
                
        if '\n' == r && fromStdin() {
            tryeval(lex)
        }         
    }
    tryeval(lex)
    
    print("\n")
}

func tryeval(lex *lexer.Lexer) {    
    defer func() {
		if x := recover(); x != nil {
		    if _, ok := x.(WaitError); ok {
	            print(wait_prompt)
	        } else {
			    print(fmt.Sprintf("%v\n", x), prompt)		
			    lex = lexer.New()
		    }
		}
	}()
	
	p := parser.New(lex)
	
	for tree := p.Parse(); nil != tree; tree = p.Parse() {
        if expr := ToExpression(tree); nil != expr {  
            res := eval(cons(expr, nil), global_env)
            print(to_string(res), "\n")
        }
	}
		
    lex = lexer.New()
    print(prompt)
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

func print(ss... string) {
    if len(ss) == 0 {
        return
    }
    output.WriteString(ss[0])
    print(ss[1:]...)
}

func fromStdin() bool { 
    return input == os.Stdin 
}

func read(rune_stream chan rune) {
	raw_stream := make(chan byte)
	defer close(raw_stream)
	
	useEncoding := decoder.E_UNKNOWN
	if fromStdin() { 
	    useEncoding = decoder.E_UTF8 
    }	
	go decoder.Decode(raw_stream, rune_stream, useEncoding)
		
	buffer := make([]byte, 64, 64)		
	for {
		c, err := input.Read(buffer)
		
		if 0 == c && io.EOF == err {
		    break
	    } else if nil != err {
		    panic(err)
		}		
		
		for i, val := range buffer{
			if i < c {
				raw_stream <- val
			}
		}		
	}
}
