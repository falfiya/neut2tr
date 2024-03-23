# neut2tr

*Northeastern University Type Comments to typed/tacket*

## Example NEU Type Comments

### Enum

```rkt
; A TrafficLightColor is one of:
; - "Red"
; - "Yellow"
; - "Green"
(define-type TrafficLightColor (U "Red" "Yellow" "Green"))
```

```rkt
; A CustomTruthy is one of:
; - #t
; - 1
; - "true"
(define-type CustomTruthy (U #t 1 "true"))
```

### Union

```rkt
; A StringOrZero is one of
; - String
; - 0
(define-type StringOrZero (U String 0))
```

```rkt
; A [Listof X] is one of:
; - '()
; - (cons X [Listof X])
(define-type (Listof X) (U '() (cons X [Listof X])))
```

```rkt
; A [Maybe X] is one of:
; - X
; - #f
(define-type (Maybe X) (U X #f))
```

### Function

```rkt
; display-clock : Minute -> Image
(: display-clock (-> Minute Image))
```

```rkt
; generate-next : [Listof Real] -> [String -> Real]
(: generate-next (-> [Listof Real] (-> String Real)))
```

```rkt
; map : {X Y} [X -> Y] [Listof X] -> [Listof Y]
(: map (All (X Y) (-> (-> X Y) [Listof X] [Listof Y])))
```

```rkt
; foldr : {X Y} [X Y -> Y] Y [Listof X] -> Y
(: foldr (All (X Y) (-> (-> X Y Y) Y [Listof X] Y)))
```
