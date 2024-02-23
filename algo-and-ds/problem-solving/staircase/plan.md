# Staircase Ascent

## Problem

- Given a staircase of N steps
- How many different ways can you ascend the staircase using groups of 1,2, and 3 steps?

## Examples

n = 1 -> [(1)]
n = 2 -> [(1,1), (2)]
n = 3 -> [(1,1,1), (2,1), (1,2), (3)]

n = 4 -> 
(1,3),(3,1)
(2,1,1), (2,2),  (1,1,2), (2,2)
(1,1,1,1), (1,2,1), (1,1,2), (1,3)
(1,1,1,1), (2,1,1), (1,2,1), (3,1)
[(1,3),(3,1),(2,1,1),(2,2),(1,1,2),(1,1,1,1),(1,2,1)]

n = 5 -> 
(1,1,3),(1,3,1),(1,2,1,1),(1,2,2),(1,1,1,2),(1,1,1,1,1),(1,1,2,1)
(3,1,1),(2,1,1,1),(2,2,1),(2,1,2),(2,3),(3,2)

General top down ->
Make a set (exclude duplicates) with:
 - N
 - N-1 with 1 prepended to each
 - N-1 with 1 appended to each 
 - N-2 with 2 prepended
 - N-2 with 2 appended
 - N-3 with 3 prepended
 - N-3 with 3 appended

General bottom up ->
- Start with N, [a(), b(), c(), d()]
- If N <= 0, return the set made from combiniung all elements in accumulator
- Otherwise recursively call: 
  N-1, [(N), pre/ap-pend(1,a()), pre/ap-pend(2,b()), pre/ap-pend(3,c())]

4, 1, (), (), ()

  

4, 

## Data
Set of lists for each N (or maybe just the last 3?)

## Algorithm
- Start with N and [Set()]
- Start with N and [(1)]


