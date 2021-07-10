#!/bin/bash

if ! [[ -v WWW_ROOT ]]
then
    1>&2 echo "WWW_ROOT is not set, not doing anything."
    sleep 3600
    exit 1
fi

echo "Cloning pages to $WWW_ROOT"
dbfile="$WWW_ROOT/db.json"
dburl="localhost:8001/apis/ahojukka5.com/v1/namespaces/dummysites/dummysites"

while :
do
    echo "Downloading sync list to $dbfile"
    curl -so $dbfile $dburl
    for SRC in $(cat $dbfile | jq -r '.items[].spec.website_url')
    do
        echo "Cloning $SRC to $WWW_ROOT/$SRC"
        wget -mpck --html-extension --user-agent="" -e robots=off --wait 1 -P $WWW_ROOT $SRC
    done
    
    echo "Creating index.html"
    echo "<html><body><ul>" > $WWW_ROOT/index.html
    for SRC in $(cat $dbfile | jq -r '.items[].spec.website_url')
    do
        echo "<li><a href=/$SRC>$SRC</a></li>" >> $WWW_ROOT/index.html
    done
    echo "</ul></body></html>" >> $WWW_ROOT/index.html
    
    echo "Cycle done. Sleeping 600 seconds ..."
    sleep 600
done
