#!/bin/sh

set -xe

go run main.go data.go > trie.dot
dot -Tsvg trie.dot > /mnt/c/Users/anujZ/Downloads/trie.svg
