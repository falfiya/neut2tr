#lang racket

; I am a line comment!
#| I am a block comment |#

; I am a line comment with a #| in it
;not all line comments have spaces

#| oh no, is this a block comment
; or is it a line comment? |#

#| Who gets the |? #|#  |#|#

#|
   ##|#
   ||#||#
#|
   what a mess#|#|## |#
|#
||#

#\##|yep, this is another block comment|#

(print "This is code!")
#\;(print "Hey look another line comment... oh wait")
(print "I'd never do something as sneaky as putting a ;comment in a String")
(print "\"; surely this must be a comment since the string is over")
(print "Did you hand\\e string escapes \'proper\\y\"?\"\\")

#\"(print "Oh, you thought this was all a string?")

#| is that #\| a character? (print "this is not code!")|#

#;(print "wait, this is a comment too? or it's in between?") #|
oh lord, this is a block comment
   #|
      this is a block comment INSIDE a block comment
   |#|
      this is only 1 comment deep
|#
(print "This is also code")

#\#;(print "I thought that #; was supposed to comment out an entire s-expression?")

(print #\" "Hah, tricked you! We're still inside a string here #|")
#|
   #\|#
(print "; |# Are you having fun yet?")
(print "That's not a character! You can't have #\\| inside a comment!")
;|# (print "I am not code!")
