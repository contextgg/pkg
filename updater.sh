#!/bin/bash

for d in */ ; do
  cd $d
  go get -u
  cd ..
done