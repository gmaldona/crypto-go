#! /bin/bash

if [ ! -f crypto.go ]
then
  printf "Please run script in the project root directory."
  exit 1
fi

IMAGE_NAME="gmaldona-poller" # image name

while getopts 'lph' OPTION; do
  case "$OPTION" in
    l)
      printf "Building image for local machine.\n"
      docker image rm "$IMAGE_NAME"-local
      docker build -t "$IMAGE_NAME"-local .
      ;;
    p)
      printf "Building image for AWS linux machine.\n"
      docker image rm "$IMAGE_NAME"
      docker buildx build --platform linux/amd64 -t "$IMAGE_NAME" .
      ;;
    h)
      printf "\055l\t Build a docker image for local machine.\n"
      printf "\055p\t Build a docker image for AWS linux machine.\n"
      ;;
    ?)
      printf "Building image for AWS linux machine.\n"
      docker image rm "$IMAGE_NAME"
      docker buildx build --platform linux/amd64 -t "$IMAGE_NAME" .
  ;;
esac
done