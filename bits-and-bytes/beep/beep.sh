#!/bin/bash

old_term_settings=$(stty -g)

reset_term() {
  stty "$old_term_settings"
  exit 0
}

trap reset_term INT

stty cbreak

do_beeps() {
  for ((i = 0; i < $1; i++)); do
    printf "\a"
    sleep 0.3
  done
}

repeats=0
echo "How many beeps? "
while true; do
  read repeats
  case "$repeats" in
    1|2|3|4|5|6|7|8|9)
      echo "$repeats"
      do_beeps "$repeats"
    ;;
  esac
done
