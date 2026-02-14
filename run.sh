#!/usr/bin/env bash

count=0
concurrency=2

wait(){
  wait $1
  count=count-1
}

for state in $(ls not_done); do
  for county in $(ls "not_done/${state}"); do
    while [[ $count < $concurrency ]]; do
      if [ -f "not_done/$state/$county" ]; then
        # sanity check
        #echo "$state/$county"
        date +%s
        ./county.sh $state $county &
        count=count+1
        pid=$!
        wait $pid
      fi
    done
  done
done
