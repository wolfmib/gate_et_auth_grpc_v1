#!/bin/bash

echo "[Jean]: Je vais creer 'user_proto' folder pour toi "
echo "         T'a sure ?"
echo "----------------------------"
ls
echo "----------------------------"
echo "Enter ...."
read nothing_var

mkdir user_proto
protoc -I proto proto/user_proto.proto --go_out=plugins=grpc:user_proto

echo
echo "[Jean]: Done... "
echo "----------------------"
ls user_proto
echo "-----------------------"
echo 