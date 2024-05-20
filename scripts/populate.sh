#!/bin/bash

for i in $(seq 1 50);
do
    echo $1
    aws s3 cp ./scripts/DummyFile.txt s3://coveo12345-poc-useast/folder-1/file-$i
done

exit 0