#! /bin/bash

docker stop crypto-app
docker rm crypto-app
docker run -it --name crypto-app crypto-app