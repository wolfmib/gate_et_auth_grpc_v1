#!/bin/bash


echo "[Jean]: Testing the Register api....!!!"
echo 
echo "--------------------------------------"
curl -d '{"first_name":"testing_","family_name":"whf","email":"ggggg@gmail.com"}' -H "Content-Type: application/json" -v -X POST http://localhost:8080/register
echo "-------------------------------------"
echo



