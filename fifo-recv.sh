#!/bin/sh

if [ -z $1 ]; then
    echo "Usage: fifo-recv.sh <fifo>"
    exit 1
fi

cat >> "$2"
