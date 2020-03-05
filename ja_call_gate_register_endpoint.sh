#!/bin/bash

echo "[Mary]: Default sending to localhost:8080/endpoints....."
echo "[Mary]: Hey there, please input your first_name"
read first_name

echo "family_name ?"
read family_name

echo "eamil ?"
read email 


# Bug to use that 
# my_data="'{\"first_name\":\"${first_name}\",\"family_name\":\"${family_name}\",\"email\":\"${email}\"}'"
# echo $my_data
# read nothing

my_option="\"Content-Type: application/json\""
echo $my_option
read nothing

curl -X POST -H "Content-Type: application/json" --data \
"{\"first_name\":\"${first_name}\",\"family_name\":\"${family_name}\",\"email\":\"${email}\"}"  \
-v  http://localhost:8080/register