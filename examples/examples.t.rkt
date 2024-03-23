#lang typed/racket

; A TrafficLightColor is one of:
; - "Red"
; - "Yellow"
; - "Green"
(define-type TrafficLightColor (U "Red" "Yellow" "Green"))

; A CustomTruthy is one of:
; - #t
; - 1
; - "true"
(define-type CustomTruthy (U #t 1 "true"))

; A StringOrZero is one of
; - String
; - 0
(define-type StringOrZero (U String 0))

; A [Listof X] is one of:
; - '()
; - (cons X [Listof X])
(define-type (Listof X) (U '() (cons X [Listof X])))

; A [Maybe X] is one of:
; - X
; - #f
(define-type (Maybe X) (U X #f))

; display-clock : Minute -> Image
(: display-clock (-> Minute Image))

; generate-next : [Listof Real] -> [String -> Real]
(: generate-next (-> [Listof Real] (-> String Real)))

; map : {X Y} [X -> Y] [Listof X] -> [Listof Y]
(: map (All (X Y) (-> (-> X Y) [Listof X] [Listof Y])))

; foldr : {X Y} [X Y -> Y] Y [Listof X] -> Y
(: foldr (All (X Y) (-> (-> X Y Y) Y [Listof X] Y)))
