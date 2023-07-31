#!/bin/bash

# Check if a file name is provided as an argument
if [ $# -eq 0 ]; then
    echo "Usage: $0 <file>"
    exit 1
fi

file="$1"
output="$2"
# Declare an empty array
lines=()

# Read the file line by line and add each line to the array
while IFS= read -r line; do
    lines+=("\"$line\"")
done < "$file"

# Join array elements with commas and print the result enclosed in square brackets
result="["
for ((i = 0; i < ${#lines[@]}; i++)); do
    result+="${lines[i]}"
    if ((i < ${#lines[@]} - 1)); then
        result+=","
    fi
done
result+="]"
if [ $# -eq 1 ]; then
	echo "$result"
fi
if [ $# -eq 2 ]; then
	echo "$result" | tee $output
fi

