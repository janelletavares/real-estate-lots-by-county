#!/usr/bin/env bash

count=0
concurrency=2

wait(){
  wait $1
  count=count-1
}

for state in $(ls not_done); do
  for county in $(ls "not_done/${state}"); do
    if [ -f "not_done/$state/$county" ]; then
      # sanity check
      echo "$state/$county"
      date +%s
      python3 entrypoint.py $state $county
      mkdir -p done/$state
      mv not_done/$state/$county done/$state
    fi
  done
done
