#!/bin/bash

docker build -t cave-logger ./docker/.

docker run \
  --interactive \
  --tty \
  --volume $HOME/.config/cave-logger:/root/.config/cave-logger \
  cave-logger \
  cave-logger

