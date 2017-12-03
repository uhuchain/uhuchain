#!/bin/bash


SEARCH="$1"

echo "Cleaning docker images containing $SEARCH"
{ # your 'try' block
    docker rm $(docker ps -a | grep $SEARCH | awk "{print \$1}")
    docker rmi $(docker images | grep $SEARCH | awk "{print \$3}")
} || { # your 'catch' block
    echo "nothing to do"
}
