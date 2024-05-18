#!/bin/bash

for i in $(seq 1 1001);
do
    echo $1
    aws s3 cp ./DummyFile.txt s3://coveo12345-poc/folder-1/file-$i
done
