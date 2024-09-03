#!/bin/bash

OUTPUT_DIR="./bin"
BINARY_NAME="tgw-server"

mkdir -p $OUTPUT_DIR
rm -f $OUTPUT_DIR/*

go build -o $OUTPUT_DIR/$BINARY_NAME ../

if [ $? -eq 0 ]; then
    echo "Build successful: $OUTPUT_DIR/$BINARY_NAME"
else
    echo "Build failed"
    exit 1
fi