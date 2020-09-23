#!/bin/bash

binary_path="./build/builds/app"

if [[ ! -f $binary_path ]]; then
    echo "Binary not available, make sure you did run \"local-build.sh\" before"
    exit 1
fi

"${binary_path}"
