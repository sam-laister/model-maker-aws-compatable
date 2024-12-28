#!/bin/bash

target="/home/personal/scripts/07_22_13/ford/$1"
let count=0
for f in "$target"/*; do
	echo $(basename $f)
	let count=count+1
done
echo ""
echo "Count: $count"
