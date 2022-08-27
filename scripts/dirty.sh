#!/bin/bash

if [[ $(git diff --stat) != '' ]]; then
    echo 'Dirty working directory, commit your changes before releasing'
    exit 1
fi
