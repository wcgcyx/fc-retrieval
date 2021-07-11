#!/bin/sh

mkdir -p .fc-retrieval/provider/files
file=".fc-retrieval/provider/admin.key"
echo "deaa5112ef638de9b9120cf65736dedef0f6642d4af42fabe65c55c2a2c806ff" > $file
testFile1=".fc-retrieval/provider/files/test1.txt"
echo "This is the test file 1" > $testFile1
testFile2=".fc-retrieval/provider/files/test2.txt"
echo "This is the test file 2" > $testFile2
testFile3=".fc-retrieval/provider/files/test3.txt"
echo "This is the test file 3" > $testFile3

export CONTAINER_IP="$(awk 'END{print $1}' /etc/hosts)"

echo "Container IP: $CONTAINER_IP"
echo "Starting service ..."

./main