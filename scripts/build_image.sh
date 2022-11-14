#! /bin/bash

if [ -f crypto.go ]
then
  printf "Please run script in the working directory of api/"
fi

docker build -t crypto-app .