#!/bin/bash

step="${1:-4}"  

dirs=()

for ((i=1; i<=step; i++)); do
  dirs+="tests/step$i "
done

valid=$(find $dirs -type f -name valid*.json)
invalid=$(find $dirs -type f -name invalid*.json)

for case in $valid; do
  if OUTPUT=$(cat $case | go run . $case 2>&1); then
    echo "PASS $case"
  else
    echo "FAILED: $case"
    IFS=$'\n'
    for line in $OUTPUT; do
      echo "   > $line"
    done
    echo "--------------------------"
  fi
done

for case in $invalid; do
  if OUTPUT=$(cat $case | go run . $case 2>&1); then
    echo "FAILED: $case"
    IFS=$'\n'
    for line in $OUTPUT; do
      echo "   > $line"
    done
  else
    echo "PASS $case"
  fi
done

# for case in $invalid; do
#   if ! OUTPUT=$(cat $case | go run . $case 2>&1); then
#     echo "PASS $case"
#   else
#     echo "FAILED: $case"
#     IFS=$'\n'
#     for line in $OUTPUT; do
#       echo "   > $line"
#     done
#   fi
# done
