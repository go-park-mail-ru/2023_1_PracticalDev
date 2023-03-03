#!/bin/bash

network=$(docker network ls | grep 'db-network' | awk '{print $2}')

if [ -z $network ];
then
    echo No database network found
    exit 1
else
docker run --rm \
    -e PGPASSWORD=pickpinpswd \
    --network=$network \
    -v $(pwd)/scripts/migrations/:/scripts/migrations/ \
    postgres \
      psql -h db -U pickpin -d pickpindb -f ./scripts/migrations/init.sql
fi


