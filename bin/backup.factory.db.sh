#!/bin/bash

backupfile="/opt/backup/factory.sqlite3.db."`date +"%Y%m%d"`

cp /home/iot/gopath/src/espressif.com/chip/factory/factory.sqlite3.db $backupfile

stampfile=`tempfile`
touch -d '-10 day' $stampfile

for x in `ls /opt/backup/factory.sqlite3.db.*`
do
    if [ $stampfile -nt $x ]; then
        rm -f $x
    fi
done

touch -d '-3 day' $stampfile
for x in `ls /opt/backup/dump.*.csv.zip`
do
    if [ $stampfile -nt $x ]; then
        rm -f $x
    fi
done

rm $stampfile
