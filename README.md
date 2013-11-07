clispi
======

Common Lisp interpreter written on Go language.
It implements the base part of Common Lisp but with some differences: 
 * I didn't have time to add maintainance of macros
 * It have one namespace for function and simple variables. 
 * Variable, stored simple value or lambda - have same nature. So (lambda ...) just return unnamed lambda variable - you can use it without needing to implement staight after initialising (as it is in Common Lisp).
 * So then functions - just lambdas with name. You can do with them what you can do with any variables.
 * This interpreter realises some core (base) Common Lisp functions - these functions have low level implementation. These base function-variables are consts so you can't overwrite them for the safe sake.
 * There's some differences with sintax for lazy evaluation - you can mark all expressions as lazy evaluated (but it has differ sign - '$')
