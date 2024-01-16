# Correct Binary Search

## Problem
- Given a sorted array, a target and a comparison function
- Write a binary search function
- It should return the index of the target if found
- If not found it should return -1

## Examples

```
# Middle of array
search([1, 3, 4, 5, 8, 9], 4) == 2

# Beginning of array
search([1, 3, 4, 5, 8, 9], 1) == 0

# End of array
search([1, 3, 4, 5, 8, 9], 9) == 5

# Not in, small
search([1, 3, 4, 5, 8, 9], 0) == -1

# Not in, large
search([1, 3, 4, 5, 8, 9], 10) == -1

# Not in, middle
search([1, 3, 4, 5, 8, 9], 6) == -1



```

## Data Structures
- Array/slice to be searched over

## Algorithm

High level approach:
- search at the middle index
- if it matches, return the current idx
- if target is larger, repeat for right have
- else repeat for left half
- if low > high return -1

Detailed algorithm:
- SET low = 0, high = len-1
- while low <= high:
    - mid = (high - low) / 2 + low 
    - if arr[mid] == target:
      - return mid
    - else if target > arr[mid]:
      - low = mid + 1
    - else (target < arr[mid]):
      - high = mid - 1
- return -1 




search([1, 3, 4, 5, 8, 9], 4) == 2
low = 0, high = 5
mid = 2
arr[2] = 4, return 2


search([1, 3, 4, 5, 8, 9], 1) == 0
low = 0, high = 5
mid = 2
arr[2] = 4, target lower
low = 0, high = 1
mid = 0
arr[0] = 1, return 0

search([1, 3, 4, 5, 8, 9], 9) == 5
low = 0, high = 5
mid = 2
arr[2] = 4, target higher
low = 3, high = 5
mid = 4
arr[4] = 8, target higher
low = 5, high = 5
mid = 5
arr[5] = 9, return 5

search([1, 3, 4, 5, 8, 9], 0) == -1
low = 0, high = 5
mid = 2
arr[2] = 4, target lower
low = 0, high = 1
mid = 0
arr[0] = 1, target lower
low = 0, high = -1
low ! <= high, return -1

search([]int{1, 3, 4, 5, 8}, 9) == 4
low = 0, high = 4
mid = 2
arr[2] = 4, target higher
low = 3, high = 4
mid = 3
arr[3] = 5, target higher
low = 4, high = 4
mid = 
