#!/bin/bash

# compile each cmd and install it

BIN_DIR="$HOME/go/bin"

go mod tidy
for dir in */; do
    dir="${dir%/}"
    if ! compgen -G "$dir/*.go" > /dev/null; then
        continue
    fi
    if go build -C $dir -o "$BIN_DIR/$dir" .; then
        echo "installed as $BIN_DIR/$dir"
    else
        echo "build failed for $dir" >&2
    fi
done