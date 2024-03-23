# Racket Comment Parser

This module's responsibility is to find comments in racket source code and
bundle line comments that seem like they belong to the same comment block.

`awful.rkt` is a difficult-to-parse file used in development for testing if this
module correctly identifies comments.
