#!/bin/bash

docker exec -it app /bin/bash -c "./wait-for-db.sh ./trading-platform /assets/test.csv"