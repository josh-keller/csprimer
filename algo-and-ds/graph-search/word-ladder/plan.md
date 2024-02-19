# Word Ladder

## Problem

Find the path from one word to another, changing only one letter at a time
- find shortest path
- one letter at a time
- don't add or subtract letters

Considerations:
- load the list of words from file
- put into data structure
- keep track of path
- use bfs to find shortest path
- 


## Example:

HEAD
heal
teal
tell
tall
TAIL

## Data Structures
- Track visited words: set/map
- Relationships between words:
  map of maps
  word: {
      cord: {},
      ford: {},
      ...
      }
  }
  Each word in a map, it points to a set of all its neighbors.

## Algorithm

### Find neighbors
1. Given a word
2. For each letter in that word
3. construct a version of the word with all 25 other possibilities for that letter, add to a list
4. Filter these possibilities based on whether they are in a set of words

### Read in words
1. Read all words into a set

### Search
Takes: 
- target
- to_search (e.g. {[word]: [cord, ford, lord, ward, wood, work, worm, worn...]} or {[word, cord]: [corn, ...]}
- visited (set of words)

next_to_search = (empty map)
next_visited = visited



for each history, poss_steps in to_search:
    for each word in poss_steps:
        if word is target:
            return reversed history ++ word
        get all neigbors
        filter to words that have not been visited
        if there are any words left:
            add to next_visited
            add [word | history]: neighbors to next_to_search

return search(target, next_to_search, next_visited)
