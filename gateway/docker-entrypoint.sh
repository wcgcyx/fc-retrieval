#!/bin/sh

mkdir -p .fc-retrieval/gateway/
file=".fc-retrieval/gateway/admin.key"
echo "deaa5112ef638de9b9120cf65736dedef0f6642d4af42fabe65c55c2a2c806ff" > $file

export CONTAINER_IP="$(awk 'END{print $1}' /etc/hosts)"

echo "Container IP: $CONTAINER_IP"
echo "Starting service ..."

./main