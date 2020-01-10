#!/bin/bash

mkdir -p darwin
mkdir -p linux
mkdir -p windows

echo "Building Linux"
GOOS=linux go build ./cmd/configula
mv configula linux/configula

echo "Building Darwin"
GOOS=darwin go build ./cmd/configula
mv configula darwin/configula

echo "Building Windows"ls

GOOS=windows go build ./cmd/configula
mv configula.exe windows/configula.exe

echo "Packaging Linux"
tar -czf configula.linux.tar.gz ./linux

echo "Packaging Darwin"
tar -czf configula.darwin.tar.gz ./darwin

echo "Packaging Windows"
zip configula.windows.zip windows/*

echo "Cleaning up"
rm -r linux/
rm -r darwin/
rm -r windows/