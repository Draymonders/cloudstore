#!/bin/bash

filepath=$1
filestore=$2

echo "filepath: " $filepath
echo "filestorepath: " $filestore

if [ ! -f $filestore ]; then
        echo "$filestore is not exist"
else
        rm -f $filestore
        echo "rm ${filestore} ok"
fi

for item in `ls $filepath | sort -n`
do
        `cat $filepath/${item} >> ${filestore}`
        echo "merge ${filepath/${item}}  to $filestore ok"
done

echo "file store ok"