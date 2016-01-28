clispi
======

Common Lisp interpreter written in Go language.
It implements the base part of Common Lisp with some differences and limitations: 
 * No macroses
 * Shared namespace for functional and simple variables. 
 * Variables, which store simple value or lambda - have same nature. So (lambda ...) just return unnamed lambda variable - it can be used without needing to call execution staight after initialising (as it is in Common Lisp).
 * So then functions are just lambdas with name. You can do with it what you can do with any variables.
 * This interpreter presents some core (base) Common Lisp functions - these functions have low level implementation. These base function-variables are consts so they cannot be overwritten for the safe sake.
 * There's some differences with sintax for lazy evaluation - all expressions are markable as lazy evaluated (but with different sign - '$')
