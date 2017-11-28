#!/bin/bash
for i in `find . -name '*.[go]' ! -path '*vendor*'`
do
  echo $i
  if ! grep -q SPDX-License-Identifier $i
  then
    cat test/scripts/licenseheader.txt $i >$i.new && mv $i.new $i
  fi
done