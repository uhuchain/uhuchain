#!/bin/bash
for i in `find . -name '*.go' ! -path '*vendor*'`
do
  if ! grep -q SPDX-License-Identifier $i
  then
    echo adding license header to $i 
    cat test/scripts/licenseheader.txt $i >$i.new && mv $i.new $i
  fi
done