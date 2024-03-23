#lang typed/racket

; A [Listof X] is one of:
; - '()
; - (cons X [Listof X])
(define-type (Listof X) (U '() (cons X (Listof X))))
(: y (Listof Number))
(define y (list 1 1))

; map : (X Y) [X -> Y] [Listof X] -> [Listof Y]
(: map (All (X Y) (-> (-> X Y) [Listof X] [Listof Y])))
