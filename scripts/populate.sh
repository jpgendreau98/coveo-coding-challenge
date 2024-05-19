#!/bin/bash

for i in $(seq 1 100);
do
    echo $1
    aws s3 cp ./DummyFile.txt s3://coveo12345-poc/folder-1/file-$i
done
