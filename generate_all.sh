#!/bin/bash
# Run `go generate ./...` for each module in the workspace
for dir in $(find . -name go.mod -exec dirname {} \;); do
    echo "Running go generate in $dir"
    (cd $dir && go generate ./...)
done
