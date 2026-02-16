#!/usr/bin/env bash

state=$1
county=$2

python3 entrypoint.py $state $county
mkdir -p done/$state
mv not_done/$state/$county done/$state
