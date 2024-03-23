#lang racket

; A TrafficLightColor is one of:
; - "Red"
; - "Yellow"
; - "Green"

; A CustomTruthy is one of:
; - #t
; - 1
; - "true"

; A StringOrZero is one of
; - String
; - 0

; A [Listof X] is one of:
; - '()
; - (cons X [Listof X])

; A [Maybe X] is one of:
; - X
; - #f

; display-clock : Minute -> Image

; generate-next : [Listof Real] -> [String -> Real]

; map : {X Y} [X -> Y] [Listof X] -> [Listof Y]

; foldr : {X Y} [X Y -> Y] Y [Listof X] -> Y
