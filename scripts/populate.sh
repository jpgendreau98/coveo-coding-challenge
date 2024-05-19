#!/bin/bash

for i in $(seq 1 100);
do
    echo $1
    aws s3 cp ./scripts/DummyFile.txt s3://coveo12345-poc-2/folder-1/file-$i
done

exit 0