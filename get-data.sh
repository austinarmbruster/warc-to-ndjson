#! /bin/bash

desiredFiles=${1:-3}
fileCnt=0

if [[ $desiredFiles -lt 0 ]]
then
    echo "desired files must be a positive integer"
    echo "$0 <desiredFiles>"
    exit 1
fi

s3Path=s3://commoncrawl/crawl-data/CC-MAIN-2021-39/segments/1631780061350.42/warc/
for f in $(aws s3 ls ${s3Path} --no-sign-request | sed 's/.* //')
do
    if [[ ! -f $f ]]
    then
        aws s3 cp ${s3Path}${f} . --no-sign-request
        let fileCnt=fileCnt+1
    fi

    if [[ $fileCnt -eq $desiredFiles ]]
    then
        exit 0
    fi
done
