#!/bin/bash

check_db() {
    pg_isready -h localhost
}

count=10

until check_db; do
    echo 'DB not responding. Waiting...'
    echo $count
    sleep $(expr 9 - $count)
    count=$(expr $count - 1)
    if [ $count -eq 0 ]; then
        echo 'Time expired, terminating test'
        exit 1
    fi
done

echo 'DB is ready for work'
