#! /bin/bash

./skip-data.sh 2

for i in $(seq 1)
do
    ./get-data.sh 1
    ./warc-to-json

    echo "Do something with the NDJSON"

    # track the completion
    for f in $(ls *.warc.gz)
    do
        rm $f
        touch $f
    done

    # remove the temporary file
    for f in $(ls *.warc.gz.ndjson)
    do
        rm $f
    done
done