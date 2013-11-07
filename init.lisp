0
'(1 2 2)
(setq define defun)
define
(define null (x) (eq nil x))
(define zero (x) (eq 0 x))
(define and (x y) 
    (cond ((cond (x t)) (cond (y t))))
)
(define or (x y)
    (cond (x t) (y t))
)
(define true (x) (and (atom x) (eq (null x) nil)))
(define twice (x) (cons x x))
(define if (p $truex $falsex) (cond (p (truex)) (t (falsex))))
(define map (f arr) 
    (if 
        (null arr) 
        nil 
        (cons (f (car arr)) (map f (cdr arr)))
    )
)
(define fold (f base n)
    (if
        (zero n) base (fold f (f base n) (sub n 1))
    )
)
(define fact (x)
    (fold mul (long 1) x)
)
(define unpairl (x)
    (if 
        (null x) nil (cons (car (car x)) (unpairl (cdr x)))
    )
)
(define prog1 (x)
    x
)
(define prog2 (x y)
    (if x y y)
)
(define progn_ (x)
    (if (atom x) (prog1 x) 
        (if (null (cdr x)) (prog1 (car x)) (prog2 (car x) (progn_ (cdr x))))
    )
)
(define progn (*x)
    (progn_ x)
)
(define lock ($x)
    x
)
(setq a (lock x))
(setq x 10)
(a)
(len '(a b c d))
(nth 2 '(a b c d))

end   
"=============================================="
(unpairl ((x 10) (1 2)))
"=============================================="

end

(prog1 (cons 1 2))
(prog2 (cons 1 2) (eq 1 2))
(progn (cons 1 2) (cons 2 3) (cons 3 4))
===============test=============
(null 1)
(null t)
(null nil)

(and 110 nil)
(and nil t)
(and 'a 'b)

(true '(asdam))
(true nil)
(true 'x)

(twice '(ads))
(twice and)
(if and 1 2)
(if (true nil) 1 2)
(if nil 3 2)

(fact 100)
(zero 1)
(zero 0)
(fold add 0 4)
(pow 10 200)
===============test=============

(define concat (x y)
    (if (null (cdr x)) (cons (car x) y) (cons (car x) (concat (cdr x) y)))
)
(define reverse_ (l) 
    (if (null (cdr l)) (car l) (cons (reverse_ (cdr l)) (cons (car l) nil)))
)
(define to_list (l)
    (if (atom l) l (concat (to_list (car l)) (cdr l)))
)
1111111111111111111111


(define f (x) (eq x nil))
999999999999998888
(define g (x) x)
(g '(nil))
z
(define f (x) (cond ((eq 'x nil) 19) (t (f (cdr 'x)))))
(f '(nil))
z

100000
(map null '(nil))   
100000
z

(defun g (x y) (cond 
    ((null x) y) 
    ((atom x) x) 
    (t (g (cdr x) 
        (cons (car x) y))
    ) 
))
(defun f (x y) (g (cons '(8 9 10) y) (cons x y)))
(g '(1 2 3) nil)
(f '(1 2 3) '(4 5 6))
1111
(set 'to_list (lambda (x)
    (cond ((true x) (cons x nil)))
))
(set 'p (lambda (x) 
    (cons 'lambda x))
)
(set 'q (lambda (x y)
    (p '(x y)))
)
(p '((x y) (cons 1 2)))
(q 1 2)
1111
(to_list 'x)
(progn '(1 2 3 4))
14241
(cond (1 10) (2 12)
)

((lambda (x) (cons x 1)) 2)
(set 'f (lambda (x) (cons x 1)))
(f 2)


(set 'p
    (lambda (x y) 
        (cons 'lambda (cons x (cons y nil)))
    )
)
(set 'p_
    (lambda (x y) 
        (cons 'cons (cons x (cons y nil)))
    )
)
(p () x)
(set 'q (lambda (x y) (eval (p x y)) ))

((q () (cons 1 2)) 2 3 4)
22222222222222222222
(set 'w (lambda (a args body)
        (set a 
            (eval (p (args) body))
        )
    )
)
(set 'apply (lambda (f x)
    (eval (f x)))
)
22222222222
(apply (w f (x) x) 10)
f
22222222222
(set 'q (lambda (a b)
        (set a (eval (eval b)))
    )
)
999999
(p (x) (cons x 10))
(set 'y (eval (p (x) (cons x 10))))
(y 7)
888888888
(q f (p (x) (cons x 10)))
(f 8)
1111111
(set 'lambda_
    (lambda (args body) 
        (cons 'lambda (cons args (cons body nil)))
    )
)
(lambda_ a b)


(set 'p (lambda (a b c)
        (set a ((lambda (tx) (ty)) (cons b nil) (cons c nil)))
    )
)

(set 'test (eval (lambda_ (x y) (set x y))))
10
(test z 10)

(test z (x y) (cons x y))
z

(set 'def
    (lambda (f args body) 
        (set f (cons 'lambda (cons args (cons body nil))))
    )
)
(def x (x y) (cons x y))
99999
((eval (def x (x y) (cons x y))) 78 15)
1000
(x 1 2)
1111

(set 'x '(lambda (x) (cons 1 x)))
x
88
(eval x)
((eval x) 10)
99
(set 'x 2)
x
10000
(defun f (x) (cons 1 x))
(f 2)

(set 'g (lambda (x) (lambda (y) (cons x y))))
(set 'c1 (g 1))
(set 'c2 (g 2))
(c1 3)
(c2 4)

(set 'defun 
    (lambda (fname vars body) 
        (set fname (lambda vars body))
    )
)
defun

(defun f (x) (cons 1 x))
(f 2)

(setf z 10)
z
(setf g (lambda (x y) (cons x (cons 'z y))))
(g 1 2)
