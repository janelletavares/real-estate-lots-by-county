#!/usr/bin/env bash


for state in $(ls not_done); do
  for county in $(ls "not_done/${state}"); do
    if [ -f "not_done/$state/$county" ]; then
        # sanity check
        #echo "$state/$county"
        python entrypoint.py $state $county
    fi
  done
done

