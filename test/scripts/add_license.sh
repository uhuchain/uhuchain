#!/bin/bash
find . -name '*.[go]' ! -path '*vendor*'
do
  echo $i
  if ! grep -q Copyright $i
  then
    cat copyright.txt $i >$i.new && mv $i.new $i
  fi
done