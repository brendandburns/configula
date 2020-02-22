#!/bin/bash
set -e

echo "Running unit tests."
go test ./pkg/configula

go build ./cmd/configula

echo "Running integration tests."
for x in examples/*.py examples/*.yml; do
    file=$(basename $x)
    echo "Testing $x"
    ./configula examples/${file} > test-output/${file}.test
    diff --strip-trailing-cr test-output/${file}.test test-output/${file}.out
done

rm test-output/*.test

echo "Tests passed!"
